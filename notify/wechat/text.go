package wechat

type TextMessage struct {
	Receivers
	MsgType string `json:"msgtype"`
	Content string `json:"content"`
}

func NewTextMessage(Builder Receivers, Content string) *TextMessage {
	return &TextMessage{
		Receivers: Builder,
		Content:   Content,
	}
}

func (m *TextMessage) ToJSON() (map[string]interface{}, error) {
	j := map[string]interface{}{
		"touser":  m.ToUser,
		"toparty": m.ToParty,
		"totag":   m.ToTag,
		"agentid": m.AgentId,
		"msgtype": "text",
		"text": map[string]string{
			"content": m.Content,
		},
	}
	return j, nil
}
