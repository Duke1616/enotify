package message

import (
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type CreateFeishuMessage struct {
	ReceiveIdType string
	ReceiveId     string
	MsgType       string
	Content       string
	Error         error
}

func (f *CreateFeishuMessage) Message() (*larkim.CreateMessageReq, error) {
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

func NewCreateFeishuMessage(ReceiveIdType, ReceiveId string, c Content) *CreateFeishuMessage {
	content, err := c.Builder()
	if err != nil {
		return &CreateFeishuMessage{
			Error: err,
		}
	}

	return &CreateFeishuMessage{
		ReceiveIdType: ReceiveIdType,
		ReceiveId:     ReceiveId,
		MsgType:       c.MsgType(),
		Content:       content,
	}
}
