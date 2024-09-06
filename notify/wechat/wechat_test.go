package wechat

import (
	"context"
	"github.com/Duke1616/enotify/notify"
	"github.com/Duke1616/enotify/notify/wechat/card"
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
	notify notify.Notifier[map[string]interface{}]
}

func (s *WechatNotifyTestSuite) SetupSuite() {
	var err error
	corpId, ok := os.LookupEnv("WECHAT_CORP_ID")
	if !ok {
		s.T().Fatal()
	}
	corpSecret, ok := os.LookupEnv("WECHAT_CORP_SECRET")
	if !ok {
		s.T().Fatal()
	}

	s.notify, err = NewWechatNotify(corpId, corpSecret)
	require.NoError(s.T(), err)
}

func (s *WechatNotifyTestSuite) TestWechatMessage() {
	t := s.T()

	testCases := []struct {
		name       string
		req        notify.BasicNotificationMessage[map[string]any]
		wantResult bool
	}{
		{
			name: "成功发送消息-文本",
			req: NewTextMessage(NewReceiversBuilder().SetAgentId(1000004).SetToUser(
				[]string{"LuanKaiZhao"}).Build(), "hello world!"),
			wantResult: true,
		},
		{
			name: "成功发送消息-markdown",
			req: NewMarkdownMessage(NewReceiversBuilder().SetAgentId(1000004).SetToUser(
				[]string{"LuanKaiZhao"}).Build(), "## 这是一个Markdown消息\n\n**加粗** 和 *斜体*"),
			wantResult: true,
		},
		{
			name: "成功发送消息-card",
			req: NewCardMessage(NewReceiversBuilder().SetAgentId(1000004).SetToUser(
				[]string{"LuanKaiZhao"}).Build(), card.NewButtonCardBuilder().SetToMailTitle(
				card.NewMailTitle("Example Title", "Example Description")).
				SetSelection(card.ButtonSelection{
					QuestionKey: "btn_question_key1",
					Title:       "企业微信评分",
					OptionList: []card.Option{
						{ID: "btn_selection_id1", Text: "100分"},
						{ID: "btn_selection_id2", Text: "101分"},
					},
					SelectedID: "btn_selection_id1",
				}).
				SetButtonList([]card.Button{
					{Text: "按钮1", Style: 1, Key: "button_key_1"},
					{Text: "按钮2", Style: 2, Key: "button_key_2"},
				}).
				SetContentList([]card.HorizontalContentList{
					{KeyName: "Key1", Value: "Value1"},
					{KeyName: "Key2", Value: "Value2"},
				}).
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
