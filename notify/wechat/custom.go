package wechat

type CustomMessage struct {
	Custom map[string]interface{}
}

func NewCustomMessage(custom map[string]any) *CustomMessage {
	return &CustomMessage{
		Custom: custom,
	}
}

func (m *CustomMessage) Message() (map[string]any, error) {
	return m.Custom, nil
}
