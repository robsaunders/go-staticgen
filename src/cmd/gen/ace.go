package main

import (
	"html/template"
	"strings"

	"github.com/yosssi/ace"
)

type AceProcessor struct{}

func (p AceProcessor) compile(tpl *TemplateWriter) error {
	funcMap := template.FuncMap{}

	t, err := ace.Load(tpl.src, "", &ace.Options{FuncMap: funcMap, DynamicReload: true})

	if err != nil {
		return err
	}

	err = t.Execute(tpl.writer, nil)

	return err
}

func (p AceProcessor) dstfile(filename string) string {
	return aceOutputFilePath(filename)
}

func aceOutputFilePath(from string) string {
	toMake := strings.Replace(from, "sites", "public", 1)
	toMake = strings.Replace(toMake, "ace", "html", 1)
	return strings.Replace(toMake, "pages/", "", 1)
}

// func compileAce(filename string) error {
// 	var doc bytes.Buffer
// 	var siteName = filepathToSitename(filename)

// 	funcMap := template.FuncMap{
// 		"current_template": func() string { return filename },
// 		"copy":             func(rel string) string { return helperCopyFile(rel, findSharedFile(siteName, rel), siteName) },
// 	}

// 	if tpl, err := ace.Load(aceInputFilePath(filename), "", &ace.Options{
// 		FuncMap:       funcMap,
// 		DynamicReload: true,
// 		BaseDir:       filepath.Join("sites", filepathToSitename(filename), "pages"),
// 	}); err == nil {

// 		if err := tpl.Execute(&doc, nil); err != nil {
// 			return err
// 		}

// 		writeStringToFile(aceOutputFilePath(filename), doc.String())

// 	} else {
// 		return err
// 	}
// 	return nil
// }

// func aceInputFilePath(from string) string {
// 	return trimExt(strings.Replace(from, filepath.Join("sites", filepathToSitename(from), "pages"), "", 1))
// }