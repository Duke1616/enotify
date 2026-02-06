package feishu

import (
	"context"
	"fmt"

	"github.com/Duke1616/enotify/notify"
	"github.com/Duke1616/enotify/notify/feishu/dispatch"
	"github.com/Duke1616/enotify/notify/feishu/message"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// Handler 飞书消息分发器，充当 Context 角色
type Handler struct {
	createNotifier notify.Notifier[*larkim.CreateMessageReq]
	patchNotifier  notify.Notifier[*larkim.PatchMessageReq]
}

func NewHandler(lark *lark.Client) (*Handler, error) {
	c, err := NewCreateLarkNotifyByClient(lark)
	if err != nil {
		return nil, err
	}

	p, err := NewPatchLarkNotifyByClient(lark)
	if err != nil {
		return nil, err
	}
	return &Handler{
		createNotifier: c,
		patchNotifier:  p,
	}, nil
}

// prepareContent 将通用 notify.Message 转换为 message.Content 接口
func (d *Handler) prepareContent(msg *notify.Message) (message.Content, error) {
	// A. 优先处理预置对象 (复用现有卡片逻辑)
	if obj, ok := msg.Extra[KeyContentObj]; ok {
		// 断言它实现了 message.Content 接口
		if cb, ok := obj.(message.Content); ok {
			return cb, nil
		}
		return nil, fmt.Errorf("extra['%s'] does not implement message.Content interface", KeyContentObj)
	}

	return nil, fmt.Errorf("feishu content object (key: %s) required", KeyContentObj)
}

// Send 处理统一消息发送
func (d *Handler) Send(ctx context.Context, msg *notify.Message) error {
	// 1. 准备 Content 对象
	contentObj, err := d.prepareContent(msg)
	if err != nil {
		return err
	}

	// 2. 获取操作类型 (默认为 create)
	action := ActionCreate
	if val, ok := msg.Extra[KeyAction].(string); ok {
		action = val
	}

	// 3. 分发逻辑
	switch action {
	case ActionPatch:
		// Patch 模式：msg.To 列表被视为 MessageID 列表
		return dispatch.BatchPatch(ctx, msg.To, contentObj, d.patchNotifier)
	default:
		// Create 模式：msg.To 列表被视为 ReceiveID (UserID/OpenID etc)
		// 获取 ID 类型，默认为 user_id
		idType := ReceiveIDTypeUserID
		if val, ok := msg.Extra[KeyReceiveIDType].(string); ok {
			idType = val
		}
		return dispatch.BatchCreate(ctx, msg.To, idType, contentObj, d.createNotifier)
	}
}
