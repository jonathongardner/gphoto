package backup

import (
  "os"
  "archive/zip"
  "path/filepath"

  "github.com/jonathongardner/gphoto/photo"
)

func AddZip(archiveFile string, folder string) (error) {
  archive, err := zip.OpenReader(archiveFile)
  if err != nil {
    return err
  }
  defer archive.Close()

  for _, f := range archive.File {
    // skip if directory
    if f.FileInfo().IsDir() {
      continue
    }

    // copy to buffere so can read and set back to 0
    fileInArchive, err := f.Open()
    if err != nil {
      return err
    }

    ext := filepath.Ext(f.Name)
    err = photo.Write(fileInArchive, folder, ext)

    fileInArchive.Close()
    if err != nil {
      return err
    }
  }
  return nil
}

func AddDir(backupFolder string, folder string) (error) {
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
    defer file.Close()

    ext := filepath.Ext(path)

    return photo.Write(file, folder, ext)
  })
}
