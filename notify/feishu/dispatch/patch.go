package dispatch

import (
	"context"
	"fmt"

	"github.com/Duke1616/enotify/notify"
	"github.com/Duke1616/enotify/notify/feishu/message"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"go.uber.org/multierr"
)

// simplePatchMessage 适配 notify.BasicNotificationMessage 接口 (Patch)
type simplePatchMessage struct {
	req *larkim.PatchMessageReq
}

func (s *simplePatchMessage) Message() (*larkim.PatchMessageReq, error) {
	return s.req, nil
}

// BatchPatch 批量发送更新消息
func BatchPatch(ctx context.Context, messageIDs []string, content message.Content, notifier notify.Notifier[*larkim.PatchMessageReq]) error {
	var errs error
	for _, msgID := range messageIDs {
		if err := sendPatch(ctx, msgID, content, notifier); err != nil {
			errs = multierr.Append(errs, fmt.Errorf("patch message %s failed: %w", msgID, err))
		}
	}
	return errs
}

// sendPatch 发送单条更新消息
func sendPatch(ctx context.Context, messageID string, content message.Content, notifier notify.Notifier[*larkim.PatchMessageReq]) error {
	pm := message.NewPatchFeishuMessage(messageID, content)
	req, err := pm.Message()
	if err != nil {
		return err
	}
	_, err = notifier.Send(ctx, &simplePatchMessage{req: req})
	return err
}
