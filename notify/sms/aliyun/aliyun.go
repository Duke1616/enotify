package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	mySms "github.com/Duke1616/enotify/notify/sms"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"

	"github.com/ecodeclub/ekit"
	"strings"
)

type Service struct {
	client   *dysmsapi.Client
	signName string
}

func NewService(client *dysmsapi.Client, signName string) *Service {
	return &Service{
		client:   client,
		signName: signName,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []mySms.Args, numbers ...string) error {
	req := &dysmsapi.SendSmsRequest{}
	req.PhoneNumbers = ekit.ToPtr[string](strings.Join(numbers, ","))
	req.SignName = ekit.ToPtr[string](s.signName)

	argsMap := make(map[string]string, len(args))
	for _, arg := range args {
		argsMap[arg.Val] = arg.Name
	}

	bCode, err := json.Marshal(argsMap)
	if err != nil {
		return err
	}
	req.TemplateParam = ekit.ToPtr[string](string(bCode))
	req.TemplateCode = ekit.ToPtr[string](tplId)

	var resp *dysmsapi.SendSmsResponse
	runtime := &util.RuntimeOptions{}
	resp, err = s.client.SendSmsWithOptions(req, runtime)
	if err != nil {
		return err
	}

	if *resp.Body.Code != "OK" {
		return fmt.Errorf("发送失败，code: %s, 原因：%s", *resp.Body.Code, *resp.Body.Message)
	}

	return nil
}
