package template

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
)

//go:embed default/*.tmpl
var templates embed.FS

type Template struct {
	template *template.Template
}

type Option func(tmpl *template.Template)

func newTemplate(options ...Option) (*Template, error) {
	t := &Template{
		template: template.New("enotify").Option("missingkey=zero"),
	}

	for _, o := range options {
		o(t.template)
	}

	t.template.Funcs(template.FuncMap(DefaultFuncs))

	return t, nil
}

func (t *Template) Parse(r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if t.template, err = t.template.Parse(string(b)); err != nil {
		return err
	}

	return nil
}

func FromGlobs(paths []string, options ...Option) (*Template, error) {
	t, err := newTemplate(options...)
	if err != nil {
		return nil, err
	}

	err = fs.WalkDir(templates, "default", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 排除目录本身
		if d.IsDir() {
			return nil
		}

		// 加载所有默认模版文件
		f, err := templates.Open(path)
		if err = t.Parse(f); err != nil {
			f.Close()
			return err
		}

		f.Close()
		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, tp := range paths {
		if err = t.FromGlob(tp); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (t *Template) Execute(name string, data interface{}) (string, error) {
	tmpl, err := t.template.Clone()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, name, data)

	return buf.String(), err
}

func (t *Template) FromGlob(path string) error {
	p, err := filepath.Glob(path)
	if err != nil {
		return err
	}
	if len(p) > 0 {
		if t.template, err = t.template.ParseGlob(path); err != nil {
			return err
		}
	}
	return nil
}
