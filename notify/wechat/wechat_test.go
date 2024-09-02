package wechat

import (
	"context"
	"enotify/notify"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

func TestWechatNotify(t *testing.T) {
	suite.Run(t, new(WechatNotifyTestSuite))
}

type WechatNotifyTestSuite struct {
	suite.Suite
	notify notify.Notifier
}

func (s *WechatNotifyTestSuite) SetupSuite() {
	var err error
	corpId := os.Getenv("WECHAT_CORP_ID")
	corpSecret := os.Getenv("WECHAT_CORP_SECRET")

	s.notify, err = NewWechatNotify(corpId, corpSecret)

	require.NoError(s.T(), err)
}

func (s *WechatNotifyTestSuite) TestWechatMessage() {
	t := s.T()

	testCases := []struct {
		name       string
		req        notify.NotificationMessage
		wantResult bool
	}{
		{
			name: "成功发送消息",
			req: NewTextMessage(NewReceiversBuilder().SetAgentId(1000004).SetToUser(
				[]string{"sunwk"}).Build(), "hello world!"),
			wantResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()

			ok, err := s.notify.Send(ctx, tc.req)
			require.NoError(t, err)
			assert.Equal(t, ok, tc.wantResult)
		})
	}
}
