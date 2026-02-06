package feishu

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Duke1616/enotify/notify/feishu/card"
	"github.com/Duke1616/enotify/template"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestHandler(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

type HandlerTestSuite struct {
	suite.Suite
	tmpl    *template.Template
	handler *Handler
}

func (s *HandlerTestSuite) SetupSuite() {
	var err error
	appId, ok := os.LookupEnv("FEISHU_APP_ID")
	if !ok {
		// 这里改为 Skip 或者 Fatal 取决于是否强制要求本地有环境变量。
		// 参考 feishu_test.go 是 Fatal
		s.T().Fatal("FEISHU_APP_ID environment variable not set")
	}
	appSecret, ok := os.LookupEnv("FEISHU_APP_SECRET")
	if !ok {
		s.T().Fatal("FEISHU_APP_SECRET environment variable not set")
	}
	lark, err := NewLarkClient(appId, appSecret)
	require.NoError(s.T(), err)

	s.handler, err = NewHandler(lark)
	require.NoError(s.T(), err)

	s.tmpl, err = template.FromGlobs([]string{})
	require.NoError(s.T(), err)
}

func (s *HandlerTestSuite) TestSendCreate() {
	t := s.T()

	// 构造 Create 消息
	msg := NewCreateBuilder("bcegag66").
		SetReceiveIDType(ReceiveIDTypeUserID).
		SetContent(NewFeishuCustomCard(s.tmpl, "feishu-card-callback", card.NewApprovalCardBuilder().
			SetToTitle("Handler Test Create").SetToFields([]card.Field{
			{
				IsShort: false,
				Tag:     "plain_text",
				Content: "这是通过 Handler 发送的测试消息",
			},
		}).SetToCallbackValue([]card.Value{
			{
				Key:   "task_id",
				Value: "100",
			},
		}).SetToHideForm(false).Build())).
		Build()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err := s.handler.Send(ctx, msg)
	require.NoError(t, err)
}

func (s *HandlerTestSuite) TestSendPatch() {
	t := s.T()

	// 构造 Patch 消息
	// 注意：Patch 需要真实的 MessageID
	msg := NewPatchBuilder("om_2bd4af328d5a0c33c02290e59be98a72").
		SetContent(NewFeishuCustomCard(s.tmpl, "feishu-card-want", card.NewApprovalCardBuilder().
			SetToTitle("Handler Test Patch").SetToFields([]card.Field{
			{
				IsShort: false,
				Tag:     "plain_text",
				Content: "这是通过 Handler 更新的测试消息",
			},
		}).SetToCallbackValue([]card.Value{
			// Note: 原始代码这里没有 Callback Value，我保持原样，如果有的话加上
		}).SetWantResult("通过 Handler 更新成功").Build())).
		Build()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err := s.handler.Send(ctx, msg)
	require.NoError(t, err)
}

func (s *HandlerTestSuite) TestSendWithForm() {
	t := s.T()

	// 使用 Builder 模式构造消息
	msg := NewCreateBuilder("bcegag66").
		SetReceiveIDType(ReceiveIDTypeUserID).
		SetContent(NewFeishuCustomCard(s.tmpl, "feishu-card-callback", card.NewApprovalCardBuilder().
			SetToTitle("Handler Test Form").
			SetInputFields([]card.InputField{
				{
					Name:     "Reason",
					Key:      "reason",
					Type:     card.FieldTextarea,
					Required: true,
					Props: map[string]string{
						"placeholder": "请输入申请理由",
					},
				},
				{
					Name:     "Date",
					Key:      "date",
					Type:     card.FieldDate,
					Required: true,
				},
				{
					Name:     "Level",
					Key:      "level",
					Type:     card.FieldSelect,
					Required: true,
					Options: []card.InputOption{
						{Label: "Low", Value: "low"},
						{Label: "High", Value: "high"},
					},
				},
				{
					Name:     "Tags",
					Key:      "tags",
					Type:     card.FieldMultiSelect,
					Required: false,
					Options: []card.InputOption{
						{Label: "Urgent", Value: "urgent"},
						{Label: "Work", Value: "work"},
					},
				},
			}).
			SetToHideForm(false).Build())).
		Build()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err := s.handler.Send(ctx, msg)
	require.NoError(t, err)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}
