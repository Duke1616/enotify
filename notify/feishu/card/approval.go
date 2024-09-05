package card

type Builder interface {
	Build() map[string]interface{}
	SetToTitle(title string) Builder
	SetToFields(Fields []Field) Builder
}

type Field struct {
	IsShort bool
	Tag     string
	Content string
}

type Approval struct {
	Fields []Field
	Title  string
}

func (a *Approval) Build() map[string]interface{} {
	fieldsMap := make([]map[string]interface{}, len(a.Fields))
	for i, field := range a.Fields {
		fieldsMap[i] = map[string]interface{}{
			"IsShort": field.IsShort,
			"Tag":     field.Tag,
			"Content": field.Content,
		}
	}

	return map[string]interface{}{
		"Title":  a.Title,
		"Fields": fieldsMap,
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

func NewApprovalCardBuilder() Builder {
	return &Approval{}
}
