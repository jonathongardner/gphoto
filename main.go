package main

import (
  "log"
  "os"
  "errors"

  "github.com/urfave/cli/v2"
  "github.com/jonathongardner/gphoto/backup"
)


func main() {
  app := &cli.App{
    Name: "gphoto",
    Usage: "backup google photos",
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name: "dir",
        Aliases: []string{"f"},
        Value: "gphoto/",
        EnvVars: []string{"GPHOTO_DIR"},
      },
    },
  }
  app.Commands = []*cli.Command{
    {
      Name:  "add-zip",
      Usage: "Add photos in zip to gphoto folder",
      Action: func(c *cli.Context) error {
        if !c.Args().Present() {
          return errors.New("must pass argument")
        }

        dir := c.String("dir")
        if !dirExists(dir) {
          return errors.New("dir " + dir + " must exist")
        }

        for i := 0; i < c.Args().Len(); i++ {
          path := c.Args().Get(i)
          if fileExists(path) {
            err := backup.AddZip(path, dir)
            if err != nil {
              return err
            }
          } else {
            log.Println("Skipping " + path + " not found")
          }
        }

        return nil
      },
    },
    {
      Name:  "add-dir",
      Usage: "Add photos in dir to gphoto folder",
      Action: func(c *cli.Context) error {
        if !c.Args().Present() {
          return errors.New("must pass argument")
        }

        dir := c.String("dir")
        if !dirExists(dir) {
          return errors.New("dir " + dir + " must exist")
        }

        for i := 0; i < c.Args().Len(); i++ {
          path := c.Args().Get(i)
          if dirExists(path) {
            err := backup.AddDir(path, dir)
            if err != nil {
              return err
            }
          } else {
            log.Println("Skipping " + path + " not found")
          }
        }

        return nil
      },
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
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
