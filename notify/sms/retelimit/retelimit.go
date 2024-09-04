package ratelimit

import (
	"context"
	"github.com/Duke1616/enotify/notify/sms"
	ratelimit "github.com/Duke1616/enotify/pkg/retelimit"

	"fmt"
)

var errLimited = fmt.Errorf("触发了限流")

type SMSRatelimitService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &SMSRatelimitService{
		svc:     svc,
		limiter: limiter,
	}
}

func (s *SMSRatelimitService) Send(ctx context.Context, tpl string, args []sms.Args, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		return fmt.Errorf("短信服务判断是否限流出现问题，%w", err)
	}
	if limited {
		return errLimited
	}
	err = s.svc.Send(ctx, tpl, args, numbers...)
	return err
}
