package template

import (
	"reflect"
	"regexp"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type FuncMap map[string]any

var DefaultFuncs = FuncMap{
	"json": func(v interface{}) string {
		// 如果是字符串，先处理可能存在的字面量 \n
		if s, ok := v.(string); ok {
			v = strings.ReplaceAll(s, "\\n", "\n")
		}
		out, _ := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(v)
		return out
	},
	"toUpper": strings.ToUpper,
	"toLower": strings.ToLower,
	"title": func(text string) string {
		// Casers should not be shared between goroutines, instead
		// create a new caser each time this function is called.
		return cases.Title(language.AmericanEnglish).String(text)
	},
	"trimSpace": strings.TrimSpace,
	// join is equal to strings.Join but inverts the argument order
	// for easier pipelining in templates.
	"join": func(sep string, s []string) string {
		return strings.Join(s, sep)
	},
	"match": regexp.MatchString,
	"reReplaceAll": func(pattern, repl, text string) string {
		re := regexp.MustCompile(pattern)
		return re.ReplaceAllString(text, repl)
	},
	"stringSlice": func(s ...string) []string {
		return s
	},
	// date returns the text representation of the time in the specified format.
	"date": func(fmt string, t time.Time) string {
		return t.Format(fmt)
	},
	// tz returns the time in the timezone.
	"tz": func(name string, t time.Time) (time.Time, error) {
		loc, err := time.LoadLocation(name)
		if err != nil {
			return time.Time{}, err
		}
		return t.In(loc), nil
	},
	"since": time.Since,
	"last": func(index int, slice interface{}) bool {
		v := reflect.ValueOf(slice)
		return index == v.Len()-1
	},
	"safeHTML": func(html string) string {
		return html
	},
}
