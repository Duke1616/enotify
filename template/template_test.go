package template

import (
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestTemplate(t *testing.T) {
	suite.Run(t, new(TemplateTestSuite))
}

type TemplateTestSuite struct {
	suite.Suite
}

func (s *TemplateTestSuite) SetupSuite() {
}

func (s *TemplateTestSuite) TestTemplateExecute() {
	t := s.T()

	testCases := []struct {
		name       string
		tmpl       *Template
		req        interface{}
		wantResult string
	}{
		{
			name: "成功发送消息-文本",
			tmpl: func() *Template {
				tmpl, err := FromGlobs([]string{"test/test.tmpl"})
				require.NoError(t, err)
				return tmpl
			}(),
			req: map[string]interface{}{
				"Title": "hello，world",
			},
			wantResult: `
notify: "hello，world"
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.tmpl.Execute("template-test", tc.req)
			require.NoError(t, err)
			assert.Equal(t, result, tc.wantResult)
		})
	}
}
