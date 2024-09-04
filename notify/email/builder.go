package email

import (
	"enotify/template"
)

type Builder interface {
	Build() Email
	SetToForm(form []string) Builder
	SetToUser(to string) Builder
	SetToSubject(subject string) Builder
}

type Email struct {
	From        string
	To          []string
	Subject     string
	ContentType string
	Body        string
}

func (m *Email) Message() (Email, error) {
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
	//var buf bytes.Buffer
	//if err := b.Template.Execute(&buf, b.Email); err != nil {
	//	panic(err)
	//}
	//b.Email.Body = buf.String()
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
