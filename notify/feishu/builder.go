package feishu

import (
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type Feishu struct {
	CreateMessageReq *larkim.CreateMessageReq
}

func (f *Feishu) Message() (Feishu, error) {
	return Feishu{
		CreateMessageReq: f.CreateMessageReq,
	}, nil
}

func NewFeishuMessage(CreateMessageReq *larkim.CreateMessageReq) *Feishu {
	return &Feishu{
		CreateMessageReq: CreateMessageReq,
	}
}
