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
	notify *Notifier
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

	globs, err := template.FromGlobs([]string{})
	if err != nil {
		fmt.Print(err)
		return
	}

	data := make(map[string]any, 0)
	data["name"] = "张三"

	testCases := []struct {
		name       string
		req        notify.BasicNotificationMessage[Email]
		wantResult bool
	}{
		{
			name: "成功发送消息-文本",
			req: NewEmailBuilder(globs).SetToUser(form).
				SetToForm([]string{to}).
				SetToSubject("title").
				SetToContentType(HTML).
				SetToBody("", data).
				Build(),
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
