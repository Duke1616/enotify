package feishu

import (
	"context"
	"github.com/Duke1616/enotify/notify"
	"github.com/Duke1616/enotify/notify/feishu/card"
	"github.com/Duke1616/enotify/notify/feishu/message"
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
	tmpl         *template.Template
	createNotify notify.Notifier[*larkim.CreateMessageReq]
	updateNotify notify.Notifier[*larkim.PatchMessageReq]
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

	s.createNotify, err = NewCreateFeishuNotify(appId, appSecret)
	require.NoError(s.T(), err)

	s.updateNotify, err = NewPatchFeishuNotify(appId, appSecret)
	require.NoError(s.T(), err)

	s.tmpl, err = template.FromGlobs([]string{})
	require.NoError(s.T(), err)
}

func (s *FeishuNotifyTestSuite) TestCreateFeishuMessage() {
	t := s.T()

	testCases := []struct {
		name       string
		wrap       notify.NotifierWrap
		wantResult bool
	}{
		{
			name: "发送自定义-卡片消息",
			wrap: notify.WrapNotifierStatic(s.createNotify, message.NewCreateFeishuMessage("user_id", "bcegag66",
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
				}).SetToHideForm(false).Build()))),
			wantResult: true,
		},
		{
			name: "发送自定义-卡片消息-流程反馈",
			wrap: notify.WrapNotifierStatic(s.createNotify, message.NewCreateFeishuMessage("user_id", "bcegag66",
				NewFeishuCustomCard(s.tmpl, "feishu-card-progress-want", card.NewApprovalCardBuilder().SetToTitle("德玛西亚").SetToFields([]card.Field{
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
				}).SetToHideForm(false).SetWantResult("你已同意该申请, 批注：无").Build()))),
			wantResult: true,
		},
		{
			name: "发送自定义-卡片消息-图片返回",
			wrap: notify.WrapNotifierStatic(s.createNotify, message.NewCreateFeishuMessage("user_id", "bcegag66",
				NewFeishuCustomCard(s.tmpl, "feishu-card-progress-image-result", card.NewApprovalCardBuilder().SetToTitle("德玛西亚").SetToFields([]card.Field{
					{
						Tag:     "markdown",
						Content: `**审批人：** <at id=a579e467></at>`,
					},
					{
						Tag:     "markdown",
						Content: `**状态：<font color='green'> 审批中 </font>**`,
					},
				}).SetImageKey("img_v2_74ca810e-e59e-4b47-82a1-fc227fe66a8g").Build()))),
			wantResult: true,
		},
		{
			name: "发送生成模版-静态卡片消息",
			wrap: notify.WrapNotifierStatic(s.createNotify, message.NewCreateFeishuMessage("user_id", "bcegag66", NewFeishuTemplateCard(
				"AAqCtHtCQMglP", "1.0.1", map[string]string{
					"title": "德玛西亚",
				}))),
			wantResult: true,
		},
		{
			name: "发送生成模版-动态卡片消息",
			wrap: notify.WrapNotifierDynamic(s.createNotify, func() (notify.BasicNotificationMessage[*larkim.CreateMessageReq], error) {
				return message.NewCreateFeishuMessage("user_id", "bcegag66", NewFeishuTemplateCard(
					"AAqCtHtCQMglP", "1.0.1", map[string]string{
						"title": "艾欧尼亚",
					})), nil
			}),
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

func (s *FeishuNotifyTestSuite) TestUpdateFeishuMessage() {
	t := s.T()

	testCases := []struct {
		name       string
		wrap       notify.NotifierWrap
		wantResult bool
	}{
		{
			name: "工单审批修改-卡片消息",
			wrap: notify.WrapNotifierStatic(s.updateNotify, message.NewPatchFeishuMessage("om_2bd4af328d5a0c33c02290e59be98a72",
				NewFeishuCustomCard(s.tmpl, "feishu-card-want", card.NewApprovalCardBuilder().SetToTitle("德玛西亚").SetToFields([]card.Field{
					{
						IsShort: false,
						Tag:     "plain_text",
						Content: "字段1内容",
					},
				}).SetWantResult("你已同意该休假申请，并批注：好的，玩得开心").Build()))),
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
