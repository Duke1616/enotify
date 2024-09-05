package template

import (
	"bytes"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
)

type Template struct {
	template *template.Template
}

type Option func(tmpl *template.Template)

func newTemplate(options ...Option) (*Template, error) {
	t := &Template{
		template: template.New("").Option("missingkey=zero"),
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

	defaultTemplates := []string{"default.tmpl", "email.tmpl"}
	for _, file := range defaultTemplates {
		f, er := os.Open(path.Join("/Users/draken/Desktop/enotify/template/", file))
		if er != nil {
			return nil, er
		}

		if er = t.Parse(f); er != nil {
			f.Close()
			return nil, er
		}
		f.Close()
	}

	for _, tp := range paths {
		if er := t.FromGlob(tp); er != nil {
			return nil, er
		}
	}
	return t, nil
}

func FromGlobsV1(paths []string, options ...Option) (*Template, error) {
	t, err := newTemplate(options...)
	if err != nil {
		return nil, err
	}

	defaultTemplates := []string{"default.tmpl", "email.tmpl"}
	for _, file := range defaultTemplates {
		f, er := os.Open(path.Join("/Users/draken/Desktop/enotify/template/", file))
		if er != nil {
			return nil, er
		}

		if er = t.Parse(f); er != nil {
			f.Close()
			return nil, er
		}
		f.Close()
	}

	for _, tp := range paths {
		if er := t.FromGlob(tp); er != nil {
			return nil, er
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

func (t *Template) ExecuteV1(dynamic string, data interface{}) (string, error) {
	if dynamic == "" {
		return "", nil
	}
	tmpl, err := t.template.Clone()
	if err != nil {
		return "", err
	}

	tmpl, err = tmpl.New("").Option("missingkey=zero").Parse(dynamic)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
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
