package feishu

import (
	"context"
	"github.com/Duke1616/enotify/notify"
	"github.com/Duke1616/enotify/notify/feishu/card"
	"github.com/Duke1616/enotify/template"
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
	tmpl   *template.Template
	notify notify.Notifier[*larkim.CreateMessageReq]
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

	s.tmpl, err = template.FromGlobs([]string{})
	require.NoError(s.T(), err)
}

func (s *FeishuNotifyTestSuite) TestFeishuMessage() {
	t := s.T()

	testCases := []struct {
		name       string
		wrap       notify.NotifierWrap
		wantResult bool
	}{
		{
			name: "发送自定义-卡片消息",
			wrap: notify.WrapNotifier(s.notify, NewFeishuMessage("user_id", "bcegag66",
				NewFeishuCustomCard(s.tmpl, "feishu-card-callback", card.NewApprovalCardBuilder().SetToTitle("德玛西亚").SetToFields([]card.Field{
					{
						IsShort: false,
						Tag:     "plain_text",
						Content: "字段1内容",
					},
				}).SetToCallbackValue([]card.Value{
					{
						Key:   "task_id",
						Value: "10",
					},
					{
						Key:   "user_id",
						Value: "123",
					},
				}).Build()))),
			wantResult: true,
		},
		{
			name: "发送生成模版-卡片消息",
			wrap: notify.WrapNotifier(s.notify, NewFeishuMessage("user_id", "bcegag66", NewFeishuTemplateCard(
				"AAqCtHtCQMglP", "1.0.1", map[string]string{
					"title": "德玛西亚",
				}))),
			wantResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()

			ok, err := tc.wrap.Send(ctx)
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
