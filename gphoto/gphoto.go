package gphoto

import (
  "os"
  "errors"
  "time"
  "path/filepath"
  "crypto/sha256"
  "fmt"
  "io"

  "github.com/boltdb/bolt"
)

const dbFileName = ".gphoto.db"
var bucketName = []byte("ShaToPath")

type GPhoto struct {
  db  *bolt.DB
  Dir string
}

func (gphoto *GPhoto) Path(checksum []byte) (string) {
  var path string
  gphoto.db.View(func(tx *bolt.Tx) error {
  	b := tx.Bucket(bucketName)
  	path = string(b.Get(checksum))
  	return nil
  })
  return path
}

func (gphoto *GPhoto) AddPath(checksum []byte, path string) (error) {
  relativePath, err := filepath.Rel(gphoto.Dir, path)
  if err != nil {
    return err
  }

  return gphoto.db.Update(func(tx *bolt.Tx) error {
  	b := tx.Bucket(bucketName)
  	err := b.Put(checksum, []byte(relativePath))
  	return err
  })
}

func (gphoto *GPhoto) Close() (error) {
  return gphoto.db.Close()
}

func (gphoto *GPhoto) Build() (error) {
  return filepath.Walk(gphoto.Dir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }

    if info.IsDir() {
      return nil
    }

    if filepath.Base(path) == dbFileName {
      return nil
    }

    f, err := os.Open(path)
    if err != nil {
      return err
    }
    defer f.Close()

    h := sha256.New()
    if _, err := io.Copy(h, f); err != nil {
      return err
    }

    gphoto.AddPath(h.Sum(nil), path)

    return nil
  })
}

func (gphoto *GPhoto) ListFiles() (error) {
  return gphoto.db.View(func(tx *bolt.Tx) error {
  	b := tx.Bucket(bucketName)

  	c := b.Cursor()

    count := 0
    for k, v := c.First(); k != nil; k, v = c.Next() {
      fmt.Printf("%v - %s: %x\n", count, v, k)
      count += 1
    }

  	return nil
  })
}

func Open(dir string) (GPhoto, error) {
  if !dirExists(dir) {
    return GPhoto{}, errors.New("dir (" + dir + ") must exist")
  }

  dbFile := filepath.Join(dir, dbFileName)
  if !fileExists(dbFile) {
    return GPhoto{}, errors.New("db (" + dbFile + ") must exist")
  }

  db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
  if err != nil {
    return GPhoto{}, err
  }

  return GPhoto{db, dir}, nil
}

func Init(dir string) (error) {
  if !dirExists(dir) {
    return errors.New("dir (" + dir + ") must exist")
  }

  dbFile := filepath.Join(dir, dbFileName)
  if fileExists(dbFile) {
    return errors.New("db (" + dbFile + ") already exist")
  }

  db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
  if err != nil {
    return err
  }
  defer db.Close()

  db.Update(func(tx *bolt.Tx) error {
  	_, err = tx.CreateBucket([]byte(bucketName))
  	return nil
  })
  return err
}

func fileExists(path string) bool {
  stats, err := os.Stat(path)
  if errors.Is(err, os.ErrNotExist) {
    return false
  }
  return !stats.IsDir()
}

func dirExists(path string) bool {
  stats, err := os.Stat(path)
  if errors.Is(err, os.ErrNotExist) {
    return false
  }
  return stats.IsDir()
}
