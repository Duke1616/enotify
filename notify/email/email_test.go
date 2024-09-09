package email

import (
	"context"
	"fmt"
	"github.com/Duke1616/enotify/notify"
	"github.com/Duke1616/enotify/template"
	"github.com/joho/godotenv"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
	"time"
)

func TestEmailNotify(t *testing.T) {
	suite.Run(t, new(EmailNotifyTestSuite))
}

type EmailNotifyTestSuite struct {
	suite.Suite
	notify notify.Notifier[Email]
}

func (s *EmailNotifyTestSuite) SetupSuite() {
	var err error
	username, ok := os.LookupEnv("EMAIL_USERNAME")
	if !ok {
		s.T().Fatal()
	}
	password, ok := os.LookupEnv("EMAIL_PASSWORD")
	if !ok {
		s.T().Fatal()
	}

	s.notify = NewEmailNotify("smtp.qq.com", username, password, 465)

	require.NoError(s.T(), err)
}

func (s *EmailNotifyTestSuite) TestEmailMessage() {
	t := s.T()

	form, ok := os.LookupEnv("EMAIL_FORM")
	if !ok {
		s.T().Fatal()
	}
	to, ok := os.LookupEnv("EMAIL_TO")
	if !ok {
		s.T().Fatal()
	}

	globs, err := template.FromGlobs([]string{"email.tmpl"})
	if err != nil {
		fmt.Print(err)
		return
	}

	data := map[string]interface{}{
		"Title":   "My Page",
		"Heading": "Welcome to My Page",
		"Content": "This is a sample content.",
	}

	testCases := []struct {
		name       string
		warp       notify.NotifierWrap
		wantResult bool
	}{
		{
			name: "成功发送消息-文本",
			warp: func() notify.NotifierWrap {
				b, er := NewEmailBuilder(globs).SetToUser(form).
					SetToForm([]string{to}).
					SetToSubject("title").
					SetToContentType(HTML).
					SetToBody("notify-email", data).
					Build()

				require.NoError(t, er)
				return notify.WrapNotifierStatic(s.notify, b)
			}(),
			wantResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()

			ok, err = tc.warp.Send(ctx)
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
