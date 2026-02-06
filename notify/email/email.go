package email

import (
	"context"
	"crypto/tls"

	"github.com/Duke1616/enotify/notify"
	"github.com/gotomicro/ego/core/elog"
	"gopkg.in/gomail.v2"
)

type Notifier struct {
	logger *elog.Component

	smtpUrl      string
	smtpPort     int
	smtpUsername string
	smtpPassword string
}

// NewEmailNotify 构造函数
func NewEmailNotify(url, username, password string, port int) notify.Notifier[Email] {
	return &Notifier{
		logger:       elog.DefaultLogger,
		smtpUrl:      url,
		smtpPort:     port,
		smtpUsername: username,
		smtpPassword: password,
	}
}

// Send 发送邮件
func (n *Notifier) Send(ctx context.Context, notify notify.BasicNotificationMessage[Email]) (bool, error) {
	d := gomail.NewDialer(n.smtpUrl, n.smtpPort, n.smtpUsername, n.smtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	msg, err := notify.Message()
	if err != nil {
		return false, err
	}

	// 创建邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", msg.From)
	m.SetHeader("To", msg.From)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody(msg.ContentType, msg.Body)

	// 发送邮件
	if err = d.DialAndSend(m); err != nil {
		return false, err
	}

	return true, nil
}
