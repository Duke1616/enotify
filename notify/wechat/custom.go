package wechat

type CustomMessage struct {
	Custom map[string]interface{}
}

func NewCustomMessage(custom map[string]interface{}) *CustomMessage {
	return &CustomMessage{
		Custom: custom,
	}
}

func (m *CustomMessage) ToJSON() (map[string]interface{}, error) {
	return m.Custom, nil
}
