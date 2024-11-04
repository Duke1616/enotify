package feishu

import (
	"context"
	"errors"
	"fmt"
	"github.com/Duke1616/enotify/notify"
	"github.com/gotomicro/ego/core/elog"
	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"net/url"
)

const feishuApiURL = "https://qyapi.weixin.qq.com/cgi-bin/"

type CreateNotify struct {
	logger *elog.Component
	url    *url.URL

	larkC *lark.Client
}

type UpdateNotify struct {
	logger *elog.Component
	url    *url.URL

	larkC *lark.Client
}

// NewCreateFeishuNotify SDK 使用文档：https://github.com/larksuite/oapi-sdk-go/tree/v3_main
func NewCreateFeishuNotify(appId, appSecret string, opts ...lark.ClientOptionFunc) (
	notify.Notifier[*larkim.CreateMessageReq], error) {
	if appId == "" || appSecret == "" {
		return nil, errors.New("appId and appSecret must not be empty")
	}

	n := &CreateNotify{
		larkC:  lark.NewClient(appId, appSecret, opts...),
		logger: elog.DefaultLogger,
	}

	// 返回接口类型
	return n, nil
}

func NewUpdateFeishuNotify(appId, appSecret string, opts ...lark.ClientOptionFunc) (
	notify.Notifier[*larkim.UpdateMessageReq], error) {
	if appId == "" || appSecret == "" {
		return nil, errors.New("appId and appSecret must not be empty")
	}

	n := &UpdateNotify{
		larkC:  lark.NewClient(appId, appSecret, opts...),
		logger: elog.DefaultLogger,
	}

	// 返回接口类型
	return n, nil
}

func (n *UpdateNotify) Send(ctx context.Context, notify notify.BasicNotificationMessage[*larkim.UpdateMessageReq]) (
	bool, error) {
	msg, err := notify.Message()
	if err != nil {
		return false, err
	}

	resp, err := n.larkC.Im.Message.Update(ctx, msg)
	if err != nil {
		return false, err

	}

	if !resp.Success() {
		return false, fmt.Errorf("发送消息失败： Code: %d, Msg: %s, requestId: %s",
			resp.Code, resp.Msg, resp.RequestId())
	}

	fmt.Println(larkcore.Prettify(resp))
	return true, nil
}

func (n *CreateNotify) Send(ctx context.Context, notify notify.BasicNotificationMessage[*larkim.CreateMessageReq]) (
	bool, error) {
	msg, err := notify.Message()
	if err != nil {
		return false, err
	}

	resp, err := n.larkC.Im.Message.Create(ctx, msg)
	if err != nil {
		return false, err

	}

	if !resp.Success() {
		return false, fmt.Errorf("发送消息失败： Code: %d, Msg: %s, requestId: %s",
			resp.Code, resp.Msg, resp.RequestId())
	}

	fmt.Println(larkcore.Prettify(resp))
	return true, nil
}
