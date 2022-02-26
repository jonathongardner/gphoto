package backup

import (
  "os"
  "archive/zip"
  "path/filepath"
  "io"
  "log"

  "github.com/jonathongardner/gphoto/photo"
  "github.com/jonathongardner/gphoto/gphoto"
  "github.com/jonathongardner/gphoto/iom"
)

func AddZip(archiveFile string, gp gphoto.GPhoto) (error) {
  archive, err := zip.OpenReader(archiveFile)
  if err != nil {
    return err
  }
  defer archive.Close()

  for i, f := range archive.File {
    // skip if directory
    if f.FileInfo().IsDir() {
      continue
    }

    fileInArchive, err := f.Open()
    if err != nil {
      return err
    }
    // copy to memory so can be read multiple times
    source := &iom.MemoryReadWriteSeeker{}
    _, err = io.Copy(source, fileInArchive)
    fileInArchive.Close()
    if err != nil {
      return err
    }
    source.Seek(0, io.SeekStart)

    sha256 := source.SHA256()
    currentPath := gp.Path(sha256)
    if len(currentPath) != 0 {
      log.Printf("File (%v) already in db %v\n", f.Name, currentPath)
      continue
    }

    year := photo.YearPhotoTaken(source)
    source.Seek(0, io.SeekStart)

    newFile := filepath.Join(gp.Dir, year, filepath.Base(f.Name))
    err = photo.Write(source, newFile)

    if err != nil {
      return err
    }

    err = gp.AddPath(sha256, newFile)
    if err != nil {
      return err
    }
    if i % 100 == 0 {
      log.Printf("Processed: %v\n", i)
    }
  }
  return nil
}

func AddDir(backupFolder string, gp gphoto.GPhoto) (error) {
  i := 0
  return filepath.Walk(backupFolder, func(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }

    if info.IsDir() {
      return nil
    }

    file, err := os.Open(path)
    if err != nil {
      return err
    }
    // copy to memory so can be read multiple times
    source := &iom.MemoryReadWriteSeeker{}
    _, err = io.Copy(source, file)
    file.Close()

    if err != nil {
      return err
    }
    source.Seek(0, io.SeekStart)

    sha256 := source.SHA256()
    currentPath := gp.Path(sha256)
    if len(currentPath) != 0 {
      log.Printf("File (%v) already in db %v\n", path, currentPath)
      return nil
    }

    relativePath, err := filepath.Rel(backupFolder, path)
    if err != nil {
      return err
    }
    newFile := filepath.Join(gp.Dir, relativePath)

    err = photo.Write(source, newFile)
    if err != nil {
      return err
    }

    err = gp.AddPath(sha256, newFile)
    if err != nil {
      return err
    }

    if i % 100 == 0 {
      log.Printf("Processed: %v\n", i)
    }
    i += 1

    return nil
  })
}
