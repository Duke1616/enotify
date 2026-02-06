package card

import "github.com/ecodeclub/ekit/slice"

type FieldType string

const (
	// FieldInput 单行文本
	FieldInput FieldType = "input"
	// FieldTextarea 多行文本
	FieldTextarea FieldType = "textarea"
	// FieldNumber 数字
	FieldNumber FieldType = "number"
	// FieldDate 日期
	FieldDate FieldType = "date"
	// FieldSelect 下拉选择
	FieldSelect FieldType = "select"
	// FieldMultiSelect 多项选择
	FieldMultiSelect FieldType = "multi_select"
)

type InputOption struct {
	Label string `json:"label"` // 选项显示名
	Value string `json:"value"` // 选项值
}

type InputField struct {
	Name     string            `json:"name"`     // 表单字段显示名
	Key      string            `json:"key"`      // 表单字段键名（对应 Order Data Key）
	Type     FieldType         `json:"type"`     // 字段类型：input, textarea, date, number...
	Required bool              `json:"required"` // 是否必填
	Options  []InputOption     `json:"options"`  // 选项列表（用于 select 等）
	Props    map[string]string `json:"props"`    // 额外组件属性（如 placeholder）
}

// Builder 审批卡片构建器接口
type Builder interface {
	// Build 构建卡片数据
	Build() map[string]interface{}
	// SetToTitle 设置卡片标题
	SetToTitle(title string) Builder
	// SetToFields 设置卡片字段列表
	SetToFields(Fields []Field) Builder
	// SetToHideForm 设置是否隐藏表单
	SetToHideForm(hideForm bool) Builder
	// SetToCallbackValue 设置回调参数
	SetToCallbackValue(callbackValues []Value) Builder
	// SetInputFields 设置用户输入表单字段
	SetInputFields(inputFields []InputField) Builder
	// SetImageKey 设置图片 Key
	SetImageKey(imageKey string) Builder
	// SetWantResult 设置期望结果内容
	SetWantResult(content string) Builder
}

type Field struct {
	IsShort bool   `json:"is_short"`
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type Approval struct {
	Fields        []Field      `json:"fields"`
	CallbackValue []Value      `json:"callback_value"`
	HideForm      bool         `json:"hide_form"`
	Title         string       `json:"title"`
	ImageKey      string       `json:"image_key"`
	WantContent   string       `json:"want_content"`
	InputFields   []InputField `json:"input_fields"`
}

type Value struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (a *Approval) Build() map[string]interface{} {
	inputFieldsList := slice.Map(a.InputFields, func(idx int, src InputField) map[string]interface{} {
		return map[string]interface{}{
			"Name":     src.Name,
			"Key":      src.Key,
			"Type":     src.Type,
			"Required": src.Required,
			"Options": slice.Map(src.Options, func(idx int, opt InputOption) map[string]interface{} {
				return map[string]interface{}{
					"Label": opt.Label,
					"Value": opt.Value,
				}
			}),
			"Props": src.Props,
		}
	})

	var inputRows [][]map[string]interface{}
	for i := 0; i < len(inputFieldsList); i += 2 {
		end := i + 2
		if end > len(inputFieldsList) {
			end = len(inputFieldsList)
		}
		inputRows = append(inputRows, inputFieldsList[i:end])
	}

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
		"InputRows": inputRows,
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

func (a *Approval) SetInputFields(inputFields []InputField) Builder {
	a.InputFields = inputFields
	return a
}

func (a *Approval) SetWantResult(content string) Builder {
	a.WantContent = content
	return a
}

func NewApprovalCardBuilder() Builder {
	return &Approval{}
}
