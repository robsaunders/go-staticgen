package main

import (
  "io"
  "os"
  "fmt"
)

func copyFile(from string) error {
  if err := cp(from, convertSrcToDestPath(from)); err != nil {
    return err
  } else {
    consoleSuccess(fmt.Sprintf("\t%s\n", convertSrcToDestPath(from)))
  }
  return nil
}

func cp(src, dst string) error {
  s, err := os.Open(src)
  if err != nil {
    return err
  }
  // no need to check errors on read only file, we already got everything
  // we need from the filesystem, so nothing can go wrong now.
  defer s.Close()
  d, err := os.Create(dst)
  if err != nil {
    return err
  }
  if _, err := io.Copy(d, s); err != nil {
    d.Close()
    return err
  }
  return d.Close()
}