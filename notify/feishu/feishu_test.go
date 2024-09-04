package feishu

import (
	"context"
	"github.com/Duke1616/enotify/notify"
	"github.com/joho/godotenv"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
	"time"
)

func TestFeishuNotify(t *testing.T) {
	suite.Run(t, new(FeishuNotifyTestSuite))
}

type FeishuNotifyTestSuite struct {
	suite.Suite
	notify *Notifier
}

func (s *FeishuNotifyTestSuite) SetupSuite() {
	var err error
	appId, ok := os.LookupEnv("FEISHU_APP_ID")
	if !ok {
		s.T().Fatal()
	}
	appSecret, ok := os.LookupEnv("FEISHU_APP_SECRET")
	if !ok {
		s.T().Fatal()
	}

	s.notify, err = NewFeishuNotify(appId, appSecret)
	require.NoError(s.T(), err)
}

func (s *FeishuNotifyTestSuite) TestWechatMessage() {
	t := s.T()

	testCases := []struct {
		name       string
		req        notify.BasicNotificationMessage[Feishu]
		wantResult bool
	}{
		{
			name: "成功发送消息-文本",
			req: NewFeishuMessage(larkim.NewCreateMessageReqBuilder().
				ReceiveIdType(`user_id`).
				Body(larkim.NewCreateMessageReqBodyBuilder().
					ReceiveId(`bcegag66`).
					MsgType(`text`).
					Content(`{"text":"test content"}`).
					Build()).
				Build()),
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

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
