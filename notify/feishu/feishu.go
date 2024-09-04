package feishu

import (
	"context"
	"errors"
	"fmt"
	"github.com/Duke1616/enotify/notify"
	"github.com/gotomicro/ego/core/elog"
	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"net/url"
)

const feishuApiURL = "https://qyapi.weixin.qq.com/cgi-bin/"

type Notifier struct {
	logger *elog.Component
	url    *url.URL

	larkC *lark.Client
}

func NewFeishuNotify(appId, appSecret string, opts ...lark.ClientOptionFunc) (*Notifier, error) {
	if appId == "" || appSecret == "" {
		return nil, errors.New("appId and appSecret must not be empty")
	}

	n := &Notifier{
		larkC:  lark.NewClient(appId, appSecret, opts...),
		logger: elog.DefaultLogger,
	}

	return n, nil
}

func (n *Notifier) Send(ctx context.Context, notify notify.BasicNotificationMessage[Feishu]) (bool, error) {
	msg, err := notify.Message()
	if err != nil {
		return false, err
	}

	resp, err := n.larkC.Im.Message.Create(ctx, msg.CreateMessageReq)
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
