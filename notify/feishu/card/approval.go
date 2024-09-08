package card

import "github.com/ecodeclub/ekit/slice"

type Builder interface {
	Build() map[string]interface{}
	SetToTitle(title string) Builder
	SetToFields(Fields []Field) Builder
	SetToCallbackValue(callbackValues []Value) Builder
}

type Field struct {
	IsShort bool   `json:"is_short"`
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type Approval struct {
	Fields        []Field `json:"fields"`
	CallbackValue []Value `json:"callback_value"`
	Title         string  `json:"title"`
}

type Value struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (a *Approval) Build() map[string]interface{} {
	return map[string]interface{}{
		"Title": a.Title,
		"CallbackValue": slice.Map(a.CallbackValue, func(idx int, src Value) map[string]interface{} {
			return map[string]interface{}{
				"Key":   src.Key,
				"Value": src.Value,
			}
		}),
		"Fields": slice.Map(a.Fields, func(idx int, src Field) map[string]interface{} {
			return map[string]interface{}{
				"IsShort": src.IsShort,
				"Tag":     src.Tag,
				"Content": src.Content,
			}
		}),
	}
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

func NewApprovalCardBuilder() Builder {
	return &Approval{}
}
