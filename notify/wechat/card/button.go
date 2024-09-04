package card

import (
	"crypto/rand"
	"encoding/hex"
)

type Builder interface {
	Build() *buttonCard
	SetToMailTitle(mailTitle MainTitle) Builder
	SetSelection(selection ButtonSelection) Builder
	SetButtonList(buttons []Button) Builder
	SetContentList(list []HorizontalContentList) Builder
	SetTaskId(taskId string) Builder

	SetCustomFields(fields map[string]any) Builder
}

type buttonCard struct {
	CardType              string                  `json:"card_type"`
	TaskId                string                  `json:"task_id"`
	MainTitle             MainTitle               `json:"main_title"`
	ButtonSelection       ButtonSelection         `json:"button_selection"`
	ButtonList            []Button                `json:"button_list"`
	HorizontalContentList []HorizontalContentList `json:"horizontal_content_list"`
	CustomFields          map[string]any          `json:"-"`
}

type MainTitle struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type HorizontalContentList struct {
	KeyName string `json:"keyname"`
	Value   string `json:"value"`
}

type ButtonSelection struct {
	QuestionKey string   `json:"question_key"`
	Title       string   `json:"title"`
	OptionList  []Option `json:"option_list"`
	SelectedID  string   `json:"selected_id"`
}

type Option struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type Button struct {
	Type  int    `json:"type"`
	Text  string `json:"text"`
	Style int    `json:"style"`
	Key   string `json:"key"`
	Url   string `json:"url"`
}

func NewMailTitle(title, desc string) MainTitle {
	return MainTitle{
		Title: title,
		Desc:  desc,
	}
}

func NewButtonCardBuilder() Builder {
	return &buttonCard{
		CardType: "button_interaction",
	}
}

func (b *buttonCard) SetTaskId(taskId string) Builder {
	b.TaskId = taskId
	return b
}

func (b *buttonCard) SetToMailTitle(mailTitle MainTitle) Builder {
	b.MainTitle = mailTitle
	return b
}

func (b *buttonCard) SetContentList(list []HorizontalContentList) Builder {
	b.HorizontalContentList = list
	return b
}

func (b *buttonCard) SetSelection(selection ButtonSelection) Builder {
	b.ButtonSelection = selection
	return b
}

func (b *buttonCard) SetButtonList(buttons []Button) Builder {
	b.ButtonList = buttons
	return b
}

func (b *buttonCard) SetCustomFields(fields map[string]any) Builder {
	for k, v := range fields {
		b.CustomFields[k] = v
	}
	return b
}

func (b *buttonCard) Build() *buttonCard {
	if b.TaskId == "" {
		b.TaskId = generateRandomString(6)
	}
	return b
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func (b *buttonCard) Message() (map[string]any, error) {
	result := map[string]interface{}{
		"card_type":               b.CardType,
		"task_id":                 b.TaskId,
		"main_title":              b.MainTitle,
		"button_selection":        b.ButtonSelection,
		"button_list":             b.ButtonList,
		"horizontal_content_list": b.HorizontalContentList,
	}

	for k, v := range b.CustomFields {
		result[k] = v
	}

	return result, nil
}
