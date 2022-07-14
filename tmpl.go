package ginplus

import (
	"bytes"
	"errors"
	"github.com/clearcodecn/ginplus/render"
	"html/template"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

func (t *TemplateManager) walkTemplate() (map[string]*template.Template, error) {
	var templates = make(map[string]*template.Template)
	err := filepath.Walk(t.root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		withoutPrefix := strings.TrimPrefix(path, t.root+string(filepath.Separator))
		withoutPrefix = strings.Replace(withoutPrefix, "\\", "/", -1)
		paths := strings.Split(withoutPrefix, "/")
		if len(paths) == 1 {
			panic("directory at least has 2 level")
		}
		var tmpl *template.Template
		tpl := templates[paths[0]]
		if tpl == nil {
			tpl = template.New(withoutPrefix)
			tpl, err = tpl.Parse(defaultVars)
			if err != nil {
				return err
			}
		}
		if tpl.Name() == withoutPrefix {
			tmpl = tpl
		} else {
			tmpl = tpl.New(withoutPrefix)
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		tmpl, err = tmpl.Parse(string(data))
		if err != nil {
			return err
		}
		templates[paths[0]] = tmpl
		return nil
	})
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (t *TemplateManager) funcMap() template.FuncMap {
	return template.FuncMap{}
}

type TemplateManager struct {
	root      string
	suffix    string
	debug     bool
	templates map[string]*template.Template
	o         sync.Once
}

func NewTemplateManager(root string, suffix string, debug bool) *TemplateManager {
	abs, err := filepath.Abs(root)
	if err != nil {
		panic(err)
	}
	return &TemplateManager{
		root:      abs,
		suffix:    suffix,
		debug:     debug,
		templates: map[string]*template.Template{},
	}
}

func (tm *TemplateManager) Render(ctx *Context, name string, data H) (string, error) {
	dir := filepath.Dir(name)
	tpl, ok := tm.templates[dir]
	if !ok {
		return "", errors.New("template not found")
	}
	var buf = bytes.Buffer{}
	if data == nil {
		data = H{}
	}
	data["Ctx"] = ctx
	err := tpl.ExecuteTemplate(&buf, name+tm.suffix, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

const (
	defaultVars = `{{- $global := . -}}
{{- $ctx := .Context -}}
`
)

// Instance (HTMLDebug) returns an HTML instance which it realizes Render interface.
func (r *TemplateManager) Instance(ctx *Context, name string, data H) render.Render {
	var tpls map[string]*template.Template
	if r.debug {
		tmpls, err := r.walkTemplate()
		if err != nil {
			panic(err)
		}
		tpls = tmpls
	} else {
		r.o.Do(func() {
			tpls, err := r.walkTemplate()
			if err != nil {
				panic(err)
			}
			r.templates = tpls
		})
		tpls = r.templates
	}

	dir := filepath.Dir(name)
	tpl, ok := tpls[dir]
	if !ok {
		panic("template not found")
	}
	if data == nil {
		data = H{}
	}
	data["Context"] = ctx
	return &HtmlRender{
		template: tpl,
		name:     name + r.suffix,
		data:     data,
		ctx:      ctx,
	}
}

type HtmlRender struct {
	template *template.Template
	name     string
	data     any
	ctx      *Context
}

func (r *HtmlRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	if r.name == "" {
		return r.template.Execute(w, r.data)
	}
	return r.template.ExecuteTemplate(w, r.name, r.data)
}

// WriteContentType (HTML) writes HTML ContentType.
func (r HtmlRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, htmlContentType)
}
