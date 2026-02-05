package dispatch

import (
	"context"
	"fmt"

	"github.com/Duke1616/enotify/notify"
	"github.com/Duke1616/enotify/notify/feishu/message"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"go.uber.org/multierr"
)

// simpleMessage 适配 notify.BasicNotificationMessage 接口 (Create)
type simpleMessage struct {
	req *larkim.CreateMessageReq
}

func (s *simpleMessage) Message() (*larkim.CreateMessageReq, error) {
	return s.req, nil
}

// BatchCreate 批量发送创建消息
func BatchCreate(ctx context.Context, targets []string, idType string, content message.Content, notifier notify.Notifier[*larkim.CreateMessageReq]) error {
	var errs error
	for _, to := range targets {
		if err := sendCreate(ctx, to, idType, content, notifier); err != nil {
			errs = multierr.Append(errs, fmt.Errorf("send to %s failed: %w", to, err))
		}
	}
	return errs
}

// sendCreate 发送单条创建消息
func sendCreate(ctx context.Context, receiveID, idType string, content message.Content, notifier notify.Notifier[*larkim.CreateMessageReq]) error {
	fm := message.NewCreateFeishuMessage(idType, receiveID, content)
	req, err := fm.Message()
	if err != nil {
		return err
	}
	_, err = notifier.Send(ctx, &simpleMessage{req: req})
	return err
}
