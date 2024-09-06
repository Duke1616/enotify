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

func (s *FeishuNotifyTestSuite) TestFeishuMessage() {
	t := s.T()

	tmpl, err := template.FromGlobs([]string{"approval.tmpl"})
	require.NoError(s.T(), err)
	data := card.NewApprovalCardBuilder().SetToTitle("德玛西亚").SetToFields([]card.Field{
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
	}).Build()

	testCases := []struct {
		name       string
		req        notify.BasicNotificationMessage[*larkim.CreateMessageReq]
		wantResult bool
	}{
		{
			name: "发送自定义-卡片消息",
			req: NewFeishuMessage("user_id", "bcegag66",
				NewFeishuCustomCard(tmpl, "app", data)),
			wantResult: true,
		},
		{
			name: "发送生成模版-卡片消息",
			req: NewFeishuMessage("user_id", "bcegag66", NewFeishuTemplateCard(
				"AAqCtHtCQMglP", "1.0.1", map[string]string{
					"title": "德玛西亚",
				})),
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
