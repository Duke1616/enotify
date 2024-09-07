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
	Build() (Email, error)
	SetToForm(form []string) Builder
	SetToUser(to string) Builder
	SetToSubject(subject string) Builder
	SetToBody(name string, data interface{}) Builder
	SetToContentType(ContextType ContextType) Builder
}

type Email struct {
	From        string
	To          []string
	Subject     string
	ContentType string
	Body        string
	Error       error
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

	Email Email
}

func NewEmailBuilder(tmpl *template.Template) Builder {
	return &emailBuilder{
		Template: tmpl,
		Email:    Email{},
	}
}

func (b *emailBuilder) Build() (Email, error) {
	if b.Email.Error != nil {
		return Email{}, b.Email.Error
	}

	return b.Email, nil
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

func (b *emailBuilder) SetToBody(name string, data interface{}) Builder {
	body, err := b.Template.Execute(name, data)
	if err != nil {
		b.Email.Error = err
	}

	b.Email.Body = body
	return b
}

func (b *emailBuilder) SetToContentType(ContextType ContextType) Builder {
	b.Email.ContentType = string(ContextType)
	return b
}
