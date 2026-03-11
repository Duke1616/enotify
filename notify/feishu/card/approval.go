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
	// FieldTips 提示选项
	FieldTips FieldType = "tips"
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
	Value    string            `json:"value"`    // 数据值
	ReadOnly bool              `json:"readonly"` // 只读字段，比如提示用户时候使用
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
	var tipsList []string
	var filteredInputFields []InputField

	for _, field := range a.InputFields {
		if field.Type == FieldTips {
			tipsList = append(tipsList, field.Value)
		} else {
			filteredInputFields = append(filteredInputFields, field)
		}
	}

	inputFieldsList := slice.Map(filteredInputFields, func(idx int, src InputField) map[string]interface{} {
		return map[string]interface{}{
			"Name":     src.Name,
			"Key":      src.Key,
			"Type":     src.Type,
			"Value":    src.Value,
			"Required": src.Required,
			"Readonly": src.ReadOnly,
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
	for i := 0; i < len(inputFieldsList); {
		current := inputFieldsList[i]
		isReadOnly := current["Readonly"] == true

		// 如果当前是只读字段，强制独占一行
		if isReadOnly {
			inputRows = append(inputRows, []map[string]interface{}{current})
			i++
			continue
		}

		// 尝试和下一个字段合并
		if i+1 < len(inputFieldsList) {
			next := inputFieldsList[i+1]
			nextIsReadOnly := next["Readonly"] == true

			if !nextIsReadOnly {
				inputRows = append(inputRows, []map[string]interface{}{current, next})
				i += 2
				continue
			}
		}

		// 无法合并，自己一行
		inputRows = append(inputRows, []map[string]interface{}{current})
		i++
	}

	fields := slice.Map(a.Fields, func(idx int, src Field) map[string]interface{} {
		return map[string]interface{}{
			"IsShort": src.IsShort,
			"Tag":     src.Tag,
			"Content": src.Content,
		}
	})

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
		"Fields":      fields,
		// Sections 将原始 Fields 按 IsShort:false 切割成分组，
		// 供 carbon-copy-sections 模板渲染带 hr 分隔线的结构化卡片。
		"Sections":  buildSections(a.Fields),
		"InputRows": inputRows,
		"Tips":      tipsList,
	}
}

// Section 表示卡片中的一个内容分区，包含小标题和其下属的字段列表
type Section struct {
	Title  string                   // 分区小标题（lark_md 格式，IsShort:false 的字段内容）
	Fields []map[string]interface{} // 该分区内的数据字段（IsShort:true），已序列化为模板所需的 map 格式
}

// buildSections 将原始 []Field 按 IsShort:false 切割为分组。
// IsShort:false 的字段作为区块小标题，其后连续的 IsShort:true 字段归入该区块。
// 若开头无标题直接是数据字段，则将其归入一个无标题区块，保持向后兼容。
func buildSections(raw []Field) []Section {
	var sections []Section
	var current *Section

	for _, f := range raw {
		if !f.IsShort {
			// IsShort:false → 区块小标题，结束上一个 Section，开启新的
			if current != nil {
				sections = append(sections, *current)
			}
			current = &Section{Title: f.Content}
		} else {
			// IsShort:true → 数据字段，序列化后追加到当前 Section
			if current == nil {
				current = &Section{}
			}
			current.Fields = append(current.Fields, map[string]interface{}{
				"IsShort": f.IsShort,
				"Tag":     f.Tag,
				"Content": f.Content,
			})
		}
	}

	if current != nil {
		sections = append(sections, *current)
	}

	return sections
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
