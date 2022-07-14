package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
)

func main() {
	// g := ginplus.New()
	WalkTemplate("examples/templates", ".gohtml")
}

func WalkTemplate(root string, suffix string) {
	var templates = make(map[string]*template.Template)
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, suffix) {
			return nil
		}
		tplPath := strings.TrimSuffix(strings.TrimPrefix(path, filepath.Join(root)+string(filepath.Separator)), suffix)
		tplPath = strings.Replace(tplPath, "\\", "/", -1)
		tpl, err := template.New(tplPath).ParseFiles(path)
		if err != nil {
			return err
		}
		templates[tplPath] = tpl
		return nil
	})
	fmt.Println(templates)
}
