package sms

import (
	"context"
	"github.com/Duke1616/enotify/notify"
	"github.com/Duke1616/enotify/notify/sms/aliyun"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sms "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/joho/godotenv"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
	"time"
)

func TestSmsNotify(t *testing.T) {
	suite.Run(t, new(SMSNotifyTestSuite))
}

type SMSNotifyTestSuite struct {
	suite.Suite
	notify *Notifier
}

func (s *SMSNotifyTestSuite) SetupSuite() {
	var err error
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		s.T().Fatal()
	}

	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		s.T().Fatal()
	}

	config := &openapi.Config{
		AccessKeyId:     tea.String(secretId),
		AccessKeySecret: tea.String(secretKey),
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	c, err := sms.NewClient(config)
	if err != nil {
		s.T().Fatal()
	}

	smsSvc := aliyun.NewService(c, "阿里云短信测试")
	s.notify = NewSmsNotifier(smsSvc)
	require.NoError(s.T(), err)
}

func (s *SMSNotifyTestSuite) TestEmailMessage() {
	t := s.T()
	number, ok := os.LookupEnv("SMS_NUMBER")
	if !ok {
		s.T().Fatal()
	}

	testCases := []struct {
		name       string
		req        notify.BasicNotificationMessage[Sms]
		wantResult bool
	}{
		{
			name: "阿里云短信",
			req: NewSms("SMS_154950909", []Args{{
				Val:  "code",
				Name: "hello",
			}}, []string{number}...),
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
