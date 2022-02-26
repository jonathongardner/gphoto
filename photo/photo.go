package photo

import (
  "io"
  "os"
  "path/filepath"
  "strings"
  "encoding/hex"

  "github.com/jonathongardner/gphoto/iom"
)
// source := &iom.MemoryReadWriteSeeker{}
// _, err := io.Copy(source, reader)
// if err != nil {
//   return err
// }
// source.Seek(0, io.SeekStart)
//
// // Get year photo was taken so we know what folder
// // returns "unkown" if cant find it
// year := yearPhotoTaken(source)
// source.Seek(0, io.SeekStart)
// // Name file by sha256 so its uniq and dont have duplicates
// hash256 := source.SHA256()
// file := filepath.Join(yFolder, hash256 + strings.ToLower(ext))

func Write(source *iom.MemoryReadWriteSeeker, newFile string) (error) {
  // If year folder doesnt exist yet create
  folder := filepath.Dir(newFile)
  if _, err := os.Stat(folder); os.IsNotExist(err) {
    err = os.MkdirAll(folder, os.ModePerm)
    if err != nil {
      return err
    }
  }

  // check if file already exists
  _, err := os.Stat(newFile)
  if err == nil {
    // if it does add sha256 to name so its uniq
    ext := filepath.Ext(newFile)
    hash256 := hex.EncodeToString(source.SHA256())
    newFile = strings.TrimSuffix(newFile, ext) + " (" + hash256 + ")" + ext
  }

  // copy file over
  destination, err := os.Create(newFile)
  if err != nil {
    return err
  }
  defer destination.Close()
  _, err = io.Copy(destination, source)

  return err
}
