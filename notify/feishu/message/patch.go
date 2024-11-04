package message

import (
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type PatchFeishuMessage struct {
	MessageId string
	Content   string
	Error     error
}

func (f *PatchFeishuMessage) Message() (*larkim.PatchMessageReq, error) {
	if f.Error != nil {
		return nil, f.Error
	}

	return larkim.NewPatchMessageReqBuilder().
		MessageId(f.MessageId).
		Body(larkim.NewPatchMessageReqBodyBuilder().
			Content(f.Content).
			Build()).
		Build(), nil
}

func NewPatchFeishuMessage(MessageId string, c Content) *PatchFeishuMessage {
	content, err := c.Builder()
	if err != nil {
		return &PatchFeishuMessage{
			Error: err,
		}
	}

	return &PatchFeishuMessage{
		MessageId: MessageId,
		Content:   content,
	}
}
