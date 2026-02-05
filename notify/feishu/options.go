package feishu

import (
	"github.com/Duke1616/enotify/notify"
	"github.com/Duke1616/enotify/notify/feishu/message"
)

// 常量定义，避免魔法字符串
const (
	ActionCreate = "create"
	ActionPatch  = "patch"

	ReceiveIDTypeUserID  = "user_id"
	ReceiveIDTypeOpenID  = "open_id"
	ReceiveIDTypeUnionID = "union_id"
	ReceiveIDTypeEmail   = "email"
	ReceiveIDTypeChatID  = "chat_id"

	KeyAction        = "action"
	KeyReceiveIDType = "receive_id_type"
	KeyContentObj    = "content_obj"
)

// Options 定义飞书渠道特有的配置选项
type Options struct {
	// Action 操作类型: create (默认), patch, update
	Action string

	// ReceiveIDType 接收用户 ID 类型: user_id (默认), open_id, union_id, email, chat_id
	// 仅在 Action 为 create 时有效
	ReceiveIDType string

	// Content 飞书特定的消息内容对象 (Card, Text, Image 等)
	// 如果设置了此字段，将优先使用此对象而不是通用的 msg.Content/Subject
	Content message.Content
}

// WithOptions 将飞书特定选项注入到通用消息中
func WithOptions(msg *notify.Message, opts Options) *notify.Message {
	if msg.Extra == nil {
		msg.Extra = make(map[string]any)
	}

	if opts.Action != "" {
		msg.Extra[KeyAction] = opts.Action
	}

	if opts.ReceiveIDType != "" {
		msg.Extra[KeyReceiveIDType] = opts.ReceiveIDType
	}

	if opts.Content != nil {
		msg.Extra[KeyContentObj] = opts.Content
	}

	return msg
}
