package photo

import (
  "log"
  "io"
  "os"
  "path/filepath"
  "strings"

  "github.com/jonathongardner/gphoto/iom"
  "github.com/rwcarlsen/goexif/exif"
)

func yearPhotoTaken(file io.Reader) (string) {
  photo, err := exif.Decode(file)
  if err != nil {
    return "unknown"
  }
  tm, _ := photo.DateTime()
  return string(tm.Year())
}

func Write(reader io.Reader, folder string, ext string) (error) {
  source := &iom.MemoryReadWriteSeeker{}
  _, err := io.Copy(source, reader)
  if err != nil {
    return err
  }
  // Get year photo was taken so we know what folder
  // returns "unkown" if cant find it
  year := yearPhotoTaken(source)
  source.Seek(0, io.SeekStart)

  // If year folder doesnt exist yet create
  yFolder := filepath.Join(folder, year)
  if _, err = os.Stat(yFolder); os.IsNotExist(err) {
    err = os.Mkdir(yFolder, os.ModePerm)
    if err != nil {
      return err
    }
  }

  // Name file by sha256 so its uniq and dont have duplicates
  hash256 := source.SHA256()
  file := filepath.Join(yFolder, hash256 + strings.ToLower(ext))

  // check if file already exists
  _, err = os.Stat(file)
  if err == nil {
    log.Printf("File already exists: %v %v\n", year, hash256)
    return nil
  }

  // copy file over
  destination, err := os.Create(file)
  if err != nil {
    return err
  }
  defer destination.Close()
  _, err = io.Copy(destination, source)

  return err
}
