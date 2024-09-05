package feishu

import "encoding/json"

type cardTemplate struct {
	version    string
	templateId string
	Variables  map[string]string
}

func (c *cardTemplate) Builder() string {
	data := map[string]interface{}{
		"type": "template",
		"data": map[string]interface{}{
			"template_id":           c.templateId,
			"template_version_name": c.version,
			"template_variable":     c.Variables,
		},
	}

	vars, _ := json.Marshal(data)
	return string(vars)
}

func (c *cardTemplate) MsgType() string {
	return "interactive"
}

func NewFeishuTemplateCard(templateId string, version string, variable map[string]string) Content {
	return &cardTemplate{
		templateId: templateId,
		version:    version,
		Variables:  variable,
	}
}
