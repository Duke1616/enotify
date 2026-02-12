package card

import (
	"fmt"
	"testing"

	"github.com/Duke1616/enotify/template"
	"github.com/stretchr/testify/require"
)

func TestApprovalTemplate(t *testing.T) {
	// 加载模版
	// 注意：这里使用相对路径加载模版文件，确保测试时能找到
	tmpl, err := template.FromGlobs([]string{"../../../template/default/approval.tmpl"})
	require.NoError(t, err)

	// 构造测试数据
	approval := &Approval{
		Title: "测试审批单",
		InputFields: []InputField{
			{
				Name:     "申请理由",
				Key:      "reason",
				Type:     FieldTextarea,
				Required: true,
				Props: map[string]string{
					"placeholder": "请输入申请理由",
				},
			},
			{
				Name:  "提示信息",
				Key:   "tips_1",
				Type:  FieldTips,
				Value: "这是一段提示信息，只有展示作用",
			},
			{
				Name:     "只读字段",
				Key:      "readonly_1",
				Type:     FieldInput,
				Value:    "这是只读内容",
				ReadOnly: true,
			},
			{
				Name: "日期",
				Key:  "date",
				Type: FieldDate,
			},
		},
		CallbackValue: []Value{
			{Key: "key1", Value: "val1"},
		},
	}

	// 执行 Build 转换数据
	data := approval.Build()

	// 渲染模版
	result, err := tmpl.Execute("approval", data)
	require.NoError(t, err)

	// 打印结果查看效果
	fmt.Println(result)
}
