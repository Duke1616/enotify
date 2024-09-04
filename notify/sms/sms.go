package sms

type Sms struct {
	tmpl    string
	args    []Args
	numbers []string
}

func (s *Sms) Message() (Sms, error) {
	return Sms{
		tmpl:    s.tmpl,
		args:    s.args,
		numbers: s.numbers,
	}, nil
}

func NewSms(tmpl string, args []Args, numbers ...string) *Sms {
	return &Sms{
		tmpl:    tmpl,
		args:    args,
		numbers: numbers,
	}
}
