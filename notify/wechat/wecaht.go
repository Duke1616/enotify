package wechat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Duke1616/enotify/notify"
	"github.com/Duke1616/enotify/pkg/httpx"
	"github.com/gotomicro/ego/core/elog"
	"io"
	"net/http"
	"net/url"
	"time"
)

const wechatApiURL = "https://qyapi.weixin.qq.com/cgi-bin/"

type Notifier struct {
	logger *elog.Component
	client *http.Client
	url    *url.URL

	corpId        string
	agentId       string
	corpSecret    string
	accessToken   string
	accessTokenAt time.Time
}

type Option func(*Notifier)

func NewWechatNotify(corpId, corpSecret string, opts ...Option) (notify.Notifier[map[string]any], error) {
	if corpId == "" || corpSecret == "" {
		return nil, errors.New("corpId and corpSecret must not be empty")
	}

	n := &Notifier{
		client:     &http.Client{},
		corpId:     corpId,
		corpSecret: corpSecret,
		logger:     elog.DefaultLogger,
	}

	for _, o := range opts {
		o(n) // 直接应用选项
	}

	// 如果没有 URL，设置默认值
	if n.url == nil {
		parsedURL, err := url.Parse(wechatApiURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse base URL: %w", err)
		}
		n.url = parsedURL
	}

	return n, nil
}

// WithUrl 选项函数，用于自定义 URL
func WithUrl(wxUrl *url.URL) Option {
	return func(n *Notifier) {
		n.url = wxUrl
	}
}

// WithAgentId 选项函数，用于自定义 AgentId
func WithAgentId[T any](agentId string) Option {
	return func(n *Notifier) {
		n.agentId = agentId
	}
}

func (n *Notifier) Send(ctx context.Context, notify notify.BasicNotificationMessage[map[string]any]) (bool, error) {
	// 每隔两个小时刷新 Token
	err := n.renNewOrRefreshToken(ctx)
	if err != nil {
		return false, err
	}

	// 发送消息通知
	err = n.send(ctx, notify)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (n *Notifier) CopyUrl() *url.URL {
	v := *n.url
	return &v
}

func (n *Notifier) send(ctx context.Context, message notify.BasicNotificationMessage[map[string]any]) error {
	// 传递参数
	parameters := url.Values{}
	parameters.Add("access_token", n.accessToken)

	// 拼接URL
	tUrl := n.CopyUrl()
	tUrl.Path += "message/send"
	tUrl.RawQuery = parameters.Encode()

	// 发送数据
	headers := make(map[string]string, 0)
	j, err := message.Message()
	if err != nil {
		return err
	}

	body, err := json.Marshal(j)
	if err != nil {
		return err
	}

	resp, err := httpx.PostRequest(ctx, n.client, tUrl.String(), bytes.NewBuffer(body), headers)
	if resp.Header.Get("Error-Code") != "0" {
		return fmt.Errorf("err: %s code: %s", resp.Header.Get("Error-Msg"), resp.Header.Get("Error-Code"))
	}

	if err != nil {
		return err
	}

	return nil
}

func (n *Notifier) renNewOrRefreshToken(ctx context.Context) error {
	if n.accessToken == "" || time.Since(n.accessTokenAt) > 2*time.Hour {
		// 传递参数
		parameters := url.Values{}
		parameters.Add("corpsecret", n.corpSecret)
		parameters.Add("corpid", n.corpId)

		// 拼接URL
		tUrl := n.CopyUrl()
		tUrl.Path += "gettoken"
		tUrl.RawQuery = parameters.Encode()

		resp, err := httpx.GetRequest(ctx, n.client, tUrl.String())
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		var wechatToken struct {
			AccessToken string `json:"access_token"`
		}

		if err = json.Unmarshal(body, &wechatToken); err != nil {
			return err
		}

		if wechatToken.AccessToken == "" {
			return fmt.Errorf("invalid APISecret for CorpID: %s", n.corpId)
		}

		// Cache accessToken
		n.accessToken = wechatToken.AccessToken
		n.accessTokenAt = time.Now()
	}

	return nil
}
