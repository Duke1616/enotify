package failover

import (
	"context"
	"enotify/notify/sms"
	"errors"

	"log"
)

type SMSFailOverService struct {
	svcs []sms.Service

	idx uint64
}

func NewFailOverSMSService(svcs []sms.Service) sms.Service {
	return &SMSFailOverService{
		svcs: svcs,
	}
}

func (f *SMSFailOverService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		// 发送成功
		if err == nil {
			return nil
		}
		// 正常这边，输出日志
		// 要做好监控
		log.Println(err)
	}
	return errors.New("全部服务商都失败了")
}
