package sms

import (
	"context"
	"github.com/Duke1616/enotify/notify"
)

type Service interface {
	Send(ctx context.Context, tmpl string, args []Args, numbers ...string) error
}

type Notifier struct {
	Service Service
}

func NewSmsNotifier(smsSvc Service) *Notifier {
	return &Notifier{
		Service: smsSvc,
	}
}

func (s *Notifier) Send(ctx context.Context, message notify.BasicNotificationMessage[Sms]) (bool, error) {
	sms, err := message.Message()
	if err != nil {
		return false, err
	}

	err = s.Service.Send(ctx, sms.tmpl, sms.args, sms.numbers...)
	if err != nil {
		return false, err
	}

	return true, nil
}

type Args struct {
	Val  string
	Name string
}
