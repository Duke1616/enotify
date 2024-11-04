package message

import (
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type UpdateFeishuMessage struct {
	MessageId string
	MsgType   string
	Content   string
	Error     error
}

func (f *UpdateFeishuMessage) Message() (*larkim.UpdateMessageReq, error) {
	if f.Error != nil {
		return nil, f.Error
	}

	return larkim.NewUpdateMessageReqBuilder().
		MessageId(f.MessageId).
		Body(larkim.NewUpdateMessageReqBodyBuilder().
			MsgType(f.MsgType).
			Content(f.Content).
			Build()).
		Build(), nil
}

func NewUpdateFeishuMessage(MessageId string, c Content) *UpdateFeishuMessage {
	content, err := c.Builder()
	if err != nil {
		return &UpdateFeishuMessage{
			Error: err,
		}
	}

	return &UpdateFeishuMessage{
		MessageId: MessageId,
		MsgType:   c.MsgType(),
		Content:   content,
	}
}
