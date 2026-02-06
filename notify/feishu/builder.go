package feishu

import (
	"github.com/Duke1616/enotify/notify"
	"github.com/Duke1616/enotify/notify/feishu/message"
)

// MessageBuilder 用于简化飞书消息的构建
type MessageBuilder struct {
	msg  *notify.Message
	opts Options
}

// NewCreateBuilder 创建发送消息构建器
// receivers: 接收者 ID 列表 (如 UserID, OpenID 等)
func NewCreateBuilder(receivers ...string) *MessageBuilder {
	return &MessageBuilder{
		msg: &notify.Message{
			To: receivers,
		},
		opts: Options{
			// Create 模式默认配置
			Action:        ActionCreate,
			ReceiveIDType: ReceiveIDTypeUserID,
		},
	}
}

// NewPatchBuilder 创建更新消息构建器
// messageIDs: 需要更新的飞书消息 ID 列表
func NewPatchBuilder(messageIDs ...string) *MessageBuilder {
	return &MessageBuilder{
		msg: &notify.Message{
			To: messageIDs,
		},
		opts: Options{
			// Patch 模式配置
			Action: ActionPatch,
		},
	}
}

// SetReceiveIDType 设置接收者 ID 类型 (仅 Create 模式有效)
func (b *MessageBuilder) SetReceiveIDType(idType string) *MessageBuilder {
	b.opts.ReceiveIDType = idType
	return b
}

// SetContent 设置具体的消息内容 (Card, Text 等)
func (b *MessageBuilder) SetContent(content message.Content) *MessageBuilder {
	b.opts.Content = content
	return b
}

// Build 构建并返回最终的 notify.Message
func (b *MessageBuilder) Build() *notify.Message {
	return WithOptions(b.msg, b.opts)
}
