package main

import (
	"fmt"
	"os"
	"path/filepath"
)

/// Process all sites.
func processSites() error {
	sites, err := allSites()
	if err != nil {
		return err
	}

	for _, site := range sites {
		consoleInfo("\nProcessing Site: " + cfg.ServerURL() + site)
		if err := makeDirIfMissing(filepath.Join(cfg.DestDir, site)); err != nil {
			return err
		}
		if err := processSite(site); err != nil {
			return err
		}
	}
	return err
}

/// Process a single site.
func processSite(sitename string) error {
	os.RemoveAll(cfg.ErrorFile())
	os.RemoveAll(filepath.Join(cfg.DestDir, sitename, "*.*"))

	var err error
	var files []string

	if files, err = RecursiveGlob(filepath.Join(cfg.SrcDir, sitename)); err == nil {
		for _, name := range files {

			// guard against partials and dotfiles
			prefix := filepath.Base(name)[0:1]
			if prefix == "_" || prefix == "." {
				continue
			}

			// make the directory at the target
			if err := makeDirIfMissing(convertSrcToDestPath(filepath.Dir(name))); err != nil {
				return err
			}

			switch filepath.Ext(name) {
			case ".sass":
				processSASS(name)
			case ".ace":
				// processACE(name)
			}
			println(name)
		}
	}

	// if err := processPages(sitename); err != nil {
	// 	return err
	// }

	// if err := processStyles(sitename); err != nil {
	// 	return err
	// }

	// if fileExists(filepath.Join(cfg.SrcDir, sitename, cfg.ImageDir)) {
	// 	if err := processImages(sitename); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// func processImages(sitename string) error {
// 	var err error
// 	var files []string

// 	if files, err = RecursiveGlob(filepath.Join(cfg.SrcDir, sitename, cfg.ImageDir)); err == nil {
// 		for _, name := range files {
// 			err = copyFile(name)
// 		}
// 	}
// 	return err
// }

func processPages(sitename string) error {

	if err := makeDirIfMissing(filepath.Join(cfg.DestDir, sitename, cfg.PageDir)); err != nil {
		return err
	}

	return processDir(
		filepath.Join(cfg.SrcDir, sitename, cfg.PageDir),
		".ace",
		&AceProcessor{BaseDir: filepath.Join(cfg.SrcDir, sitename, cfg.PageDir)})
}

func processStyles(sitename string) error {
	if err := makeDirIfMissing(filepath.Join(cfg.DestDir, sitename, cfg.StyleDir)); err != nil {
		return err
	}

	// var processor = &GcssProcessor{}
	// var processor = &SassProcessor{}
	return processDir(filepath.Join(cfg.SrcDir, sitename, cfg.StyleDir), ".sass", &SassProcessor{})
}

func processDir(srcdir string, filetype string, processor Processor) error {
	var err error

	if files, err := FileTypeGlob(srcdir, filetype); err == nil {
		for _, name := range files {
			// create directories if needed
			dir, _ := filepath.Split(processor.dstfile(name))
			makeDirIfMissing(dir)

			tpl := NewTemplateWriter(name, processor.dstfile(name))
			if tpl.err != nil {
				return tpl.err
			}

			err = processor.compile(tpl)

			if err == nil {
				consoleSuccess(fmt.Sprintf("\t" + processor.dstfile(name) + "\n"))
			} else {
				// print(name)
				return err
			}
		}
	}

	return err
}
