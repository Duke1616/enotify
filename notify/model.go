package notify

import (
	"context"
)

// MessageType 定义消息类型
type MessageType string

const (
	MsgTypeText     MessageType = "text"
	MsgTypeMarkdown MessageType = "markdown" // 部分渠道支持，不支持的可能降级为 Text
	MsgTypeCard     MessageType = "card"     // 复杂卡片，通常需要 channel 特定的 Extra 数据
)

// Message 是业务层使用的统一消息模型
// 它屏蔽了底层渠道的差异，提供了最常用的字段
type Message struct {
	// 基础字段
	To      []string    // 接收者列表 (UserID, Phone, Email 等，取决于 Channel)
	Subject string      // 标题 (邮件、飞书富文本等使用)
	Content string      // 内容 (文本内容)
	Type    MessageType // 消息类型

	// 扩展字段
	TemplateID   string         // 模版 ID (如短信、飞书模板消息)
	TemplateVars map[string]any // 模版变量
	Files        []string       // 附件链接
	Extra        map[string]any // 渠道特定的额外参数 (如飞书的 card_id, 交互配置等)
}

// Handler 定义了统一的发送接口
// 业务层调用此接口，不关心底层的泛型实现
type Handler interface {
	Send(ctx context.Context, msg *Message) error
}
