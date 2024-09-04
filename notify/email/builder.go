package email

import (
	"github.com/Duke1616/enotify/template"
)

type ContextType string

const (
	TEXT ContextType = "text/plain"
	HTML ContextType = "text/html"
)

type Builder interface {
	Build() Email
	SetToForm(form []string) Builder
	SetToUser(to string) Builder
	SetToSubject(subject string) Builder
	SetToBody(dynamicTmpl string, body map[string]any) Builder
	SetToContentType(ContextType ContextType) Builder
}

type Email struct {
	From        string
	To          []string
	Subject     string
	ContentType string
	Body        string
}

func (m Email) Message() (Email, error) {
	j := Email{
		From:        m.From,
		To:          m.To,
		Subject:     m.Subject,
		ContentType: m.ContentType,
		Body:        m.Body,
	}

	return j, nil
}

type emailBuilder struct {
	Template *template.Template
	Email    Email
}

func NewEmailBuilder(tmpl *template.Template) Builder {
	return &emailBuilder{
		Template: tmpl,
		Email:    Email{},
	}
}

func (b *emailBuilder) Build() Email {
	return b.Email
}

func (b *emailBuilder) SetToForm(form []string) Builder {
	b.Email.To = form
	return b
}

func (b *emailBuilder) SetToUser(to string) Builder {
	b.Email.From = to
	return b
}

func (b *emailBuilder) SetToSubject(subject string) Builder {
	b.Email.Subject = subject
	return b
}

func (b *emailBuilder) SetToBody(dynamicTmpl string, body map[string]any) Builder {
	// 模板内容
	dynamic := `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
</head>
<body>
    <h1>{{.Heading}}</h1>
    <p>{{.Content}}</p>
</body>
</html>
`

	// 模板数据
	w := map[string]interface{}{
		"Title":   "My Page",
		"Heading": "Welcome to My Page",
		"Content": "This is a sample content.",
	}

	data, err := b.Template.Execute(dynamic, w)
	if err != nil {
	}

	b.Email.Body = data
	return b
}

func (b *emailBuilder) SetToContentType(ContextType ContextType) Builder {
	b.Email.ContentType = string(ContextType)
	return b
}
