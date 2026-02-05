package sms

import (
	"context"
	"fmt"

	"github.com/Duke1616/enotify/notify"
)

// Handler 实现了 notify.Handler 接口，封装了 SMS 的发送逻辑
type Handler struct {
	notifier notify.Notifier[Sms]
}

// NewHandler 创建一个新的 SMS 处理器
func NewHandler(svc Service) *Handler {
	// 复用现有的 NewSmsNotifier
	n := NewSmsNotifier(svc)
	return &Handler{
		notifier: n,
	}
}

// Send 发送统一消息
func (h *Handler) Send(ctx context.Context, msg *notify.Message) error {
	// 短信通常强依赖模版 ID
	if msg.TemplateID == "" {
		// 如果没有模版ID，尝试检查是否可以降级为纯文本（取决于服务商，通常不可以）
		// 或者从 Content 自动提取（不推荐）
		return fmt.Errorf("sms channel requires TemplateID")
	}

	// 转换模版变量
	var args []Args
	if msg.TemplateVars != nil {
		for k, v := range msg.TemplateVars {
			args = append(args, Args{
				Name: k,
				Val:  fmt.Sprintf("%v", v),
			})
		}
	}

	// 使用现有的工厂方法创建 Sms 对象
	// *Sms 已经实现了 BasicNotificationMessage[Sms] 接口
	smsData := NewSms(msg.TemplateID, args, msg.To...)

	// 调用底层 Notifier 发送
	// 注意：底层 Notifier Expects BasicNotificationMessage[Sms]
	_, err := h.notifier.Send(ctx, smsData)
	if err != nil {
		return err
	}

	return nil
}
