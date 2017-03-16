package main

import (
  "bytes"
  "text/template"
  "path/filepath"
  "strings"
  "fmt"
)

func processTemplate(from string, dir string) (bytes.Buffer, error) {
  var doc bytes.Buffer
  var siteName = filepathToSitename(from)

  funcMap := template.FuncMap {
    "copy": func (rel string) string { return helperCopyFile(rel, findSharedFile(siteName, rel)) },
  }

  baseName := filepath.Base(from)
  globbedFiles, _ := PartialGlob("sites/"+siteName+"/styles", ".sass")

  tpl := template.New(baseName).Funcs(funcMap)

  if parsedTpl, err := tpl.ParseFiles(globbedFiles...); err != nil {
    return doc, err
  } else {
    tpl = parsedTpl
  }

  err := tpl.Execute(&doc, nil)
  return doc, err
}

func helperCopyFile(rel string, src string) string {
  if src == "" { createError(rel, nil) }

  dest := strings.Replace(convertSrcToDestPath(src), "_shared", filepathToSitename(src), 1)

  if err := makeDirIfMissing(filepath.Dir(dest)); err != nil { createError(src, err) }

  if err := cp(src, dest); err != nil {
    createError(src, err)
  } else {
    consoleSuccess(fmt.Sprintf("\t%s\n", dest))
  }

  println(rel + "?checksum=" + checksum(src))

  return rel + "?checksum=" + checksum(src)
}

func findSharedFile(site string, from string) string {
  if fileExists("sites/" + site + "/" + from) {
    return "sites/" + site + "/" + from
  } else {
    if fileExists("sites/_shared/" + from) {
      return "sites/_shared/" + from
    } else { return "" }
  }
}