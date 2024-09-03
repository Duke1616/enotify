package wechat

type CardMessage struct {
	Receivers
	TemplateCard
}

type TemplateCard interface {
	ToMap() (map[string]any, error)
}

func NewCardMessage(builder Receivers, card TemplateCard) *CardMessage {
	return &CardMessage{
		Receivers:    builder,
		TemplateCard: card,
	}
}

func (m *CardMessage) ToJSON() (map[string]interface{}, error) {
	templateCardMap, err := m.TemplateCard.ToMap()
	if err != nil {
		return nil, err
	}

	j := map[string]interface{}{
		"touser":        m.ToUser,
		"toparty":       m.ToParty,
		"totag":         m.ToTag,
		"agentid":       m.AgentId,
		"msgtype":       "template_card",
		"template_card": templateCardMap,
	}

	return j, nil
}
