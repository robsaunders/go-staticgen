package main

import (
  "fmt"
  "os"
  "path/filepath"
  "strings"
  "io/ioutil"
)

func processSites() {
  files, err := ioutil.ReadDir("sites")
  if err != nil { checkFatal(err) }

  for _, file := range files {
    if file.IsDir() {
      dot := filepath.Base(file.Name())[0]
      if dot != '.' && dot != '_' {
        println(dot)
        consoleInfo("\nProcessing Site: "+ file.Name())
        processSite(file.Name())
      }
    }
  }
}

func processSite(sitename string) {

  filepath.Walk("sites/"+sitename, func(name string, info os.FileInfo, err error) error {
    if info == nil || name == "sites/"+sitename {
      return err
    }

    from := filepath.ToSlash(name)
    dot := filepath.Base(name)[0]
    ext := filepath.Ext(name)

    if info.IsDir() {
      if dot == '.' || dot == '_' {
        return filepath.SkipDir
      } else {
        makeDirIfMissing(convertSrcToDestPath(from))
      }
    } else {
      // for _, exclude := range [] {
      //   if strings.HasSuffix(from, exclude) {
      //     return err
      //   }
      // }
      if dot != '.' && dot != '_' {
        switch ext {
        case ".ace":
          compileAce(from)
        case ".sass":
          compileGcss(from, sitename)
        default:
          copyFile(from)
        }
      }
    }
    return err
  })
}

func convertSrcToDestPath(filename string) string {
  return strings.Replace(filename, "sites", "public", 1)
}

func findLayoutFile(filename string) string {
  return filepath.Dir(filename) + "/_layout"
}

func trimExt(filename string) string {
  return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func makeDirIfMissing(dir string) {
  if _, err := os.Stat(dir); os.IsNotExist(err) {
    consoleInfo("[MKDIR] " + dir)
    os.MkdirAll(dir, os.ModePerm)
  }
}

func writeStringToFile(filepath string, content string) {
  fo, err := os.Create(filepath)
  if err != nil { panic(err) }

  defer fo.Close()
  fmt.Fprintf(fo, content)
}