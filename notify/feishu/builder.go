package feishu

import (
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type Content interface {
	Builder() (string, error)
	MsgType() string
}

type Feishu struct {
	ReceiveIdType string
	ReceiveId     string
	MsgType       string
	Content       string
	Error         error
}

func (f *Feishu) Message() (*larkim.CreateMessageReq, error) {
	if f.Error != nil {
		return nil, f.Error
	}

	return larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(f.ReceiveIdType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(f.ReceiveId).
			MsgType(f.MsgType).
			Content(f.Content).Build()).Build(), nil
}

func NewFeishuMessage(ReceiveIdType, ReceiveId string, c Content) *Feishu {
	content, err := c.Builder()
	if err != nil {
		return &Feishu{
			Error: err,
		}
	}

	return &Feishu{
		ReceiveIdType: ReceiveIdType,
		ReceiveId:     ReceiveId,
		MsgType:       c.MsgType(),
		Content:       content,
	}
}
