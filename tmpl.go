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

func (t *TemplateManager) walkTemplate(ctx *Context) (map[string]*template.Template, error) {
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
			tpl = template.New(withoutPrefix).Funcs(t.funcMap(ctx))
			tpl, err = tpl.Parse(defaultVars)
			if err != nil {
				return err
			}
		}
		if tpl.Name() == withoutPrefix {
			tmpl = tpl
		} else {
			tmpl = tpl.New(withoutPrefix).Funcs(t.funcMap(ctx))
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

func (t *TemplateManager) funcMap(ctx *Context) template.FuncMap {
	return template.FuncMap{
		"flash":      flash(ctx),
		"hasSession": hasSession(ctx),
	}
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
		tmpls, err := r.walkTemplate(ctx)
		if err != nil {
			panic(err)
		}
		tpls = tmpls
	} else {
		r.o.Do(func() {
			tpls, err := r.walkTemplate(ctx)
			if err != nil {
				panic(err)
			}
			r.templates = tpls
		})
		tpls = r.templates
	}
	theme := ctx.theme
	if theme == "" {
		theme = filepath.Dir(name)
	} else {
		name = theme + "/" + name
	}
	tpl, ok := tpls[theme]
	if !ok {
		panic("template not found")
	}
	if data == nil {
		data = H{}
	}
	if ctx.data != nil {
		for k, v := range ctx.data {
			data[k] = v
		}
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
	buf := bufpool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufpool.Put(buf)
	}()

	err := r.template.ExecuteTemplate(buf, r.name, r.data)
	if err != nil {
		return err
	}

	if r.ctx.beforeRender != nil {
		r.ctx.beforeRender()
	}

	_, err = buf.WriteTo(w)
	return err
}

// WriteContentType (HTML) writes HTML ContentType.
func (r HtmlRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, htmlContentType)
}
