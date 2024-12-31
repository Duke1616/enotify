package card

import "github.com/ecodeclub/ekit/slice"

type Builder interface {
	Build() map[string]interface{}
	SetToTitle(title string) Builder
	SetToFields(Fields []Field) Builder
	SetToHideForm(hideForm bool) Builder
	SetToCallbackValue(callbackValues []Value) Builder
	SetImageKey(imageKey string) Builder
	SetWantResult(content string) Builder
}

type Field struct {
	IsShort bool   `json:"is_short"`
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type Approval struct {
	Fields        []Field `json:"fields"`
	CallbackValue []Value `json:"callback_value"`
	HideForm      bool    `json:"hide_form"`
	Title         string  `json:"title"`
	ImageKey      string  `json:"image_key"`
	WantContent   string  `json:"want_content"`
}

type Value struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (a *Approval) Build() map[string]interface{} {
	return map[string]interface{}{
		"Title":    a.Title,
		"HideForm": a.HideForm,
		"CallbackValue": slice.Map(a.CallbackValue, func(idx int, src Value) map[string]interface{} {
			return map[string]interface{}{
				"Key":   src.Key,
				"Value": src.Value,
			}
		}),
		"WantContent": a.WantContent,
		"ImageKey":    a.ImageKey,
		"Fields": slice.Map(a.Fields, func(idx int, src Field) map[string]interface{} {
			return map[string]interface{}{
				"IsShort": src.IsShort,
				"Tag":     src.Tag,
				"Content": src.Content,
			}
		}),
	}
}

func (a *Approval) SetImageKey(imageKey string) Builder {
	a.ImageKey = imageKey
	return a
}

func (a *Approval) SetToHideForm(hideForm bool) Builder {
	a.HideForm = hideForm
	return a
}

func (a *Approval) SetToTitle(title string) Builder {
	a.Title = title
	return a
}

func (a *Approval) SetToFields(Fields []Field) Builder {
	a.Fields = Fields
	return a
}

func (a *Approval) SetToCallbackValue(callbackValues []Value) Builder {
	a.CallbackValue = callbackValues
	return a
}

func (a *Approval) SetWantResult(content string) Builder {
	a.WantContent = content
	return a
}

func NewApprovalCardBuilder() Builder {
	return &Approval{}
}
