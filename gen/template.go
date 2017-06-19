package gen

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

var (
	defaultFuncs = template.FuncMap{
		"greaterThanEqualF": func(b, a float64) bool {
			return b >= a
		},
		"lessThanEqualF": func(b, a float64) bool {
			return b <= a
		},
		"greaterThanEqual": func(b, a int) bool {
			return b >= a
		},
		"lessThanEqual": func(b, a int) bool {
			return b <= a
		},
		"greaterThanF": func(b, a float64) bool {
			return b > a
		},
		"lessThanF": func(b, a float64) bool {
			return b < a
		},
		"greaterThan": func(b, a int) bool {
			return b > a
		},
		"lessThan": func(b, a int) bool {
			return b < a
		},
		"trimspace": func(b string) string {
			return strings.TrimSpace(b)
		},
		"equal": func(b, a interface{}) bool {
			return b == a
		},
		"not": func(b bool) bool {
			return !!b
		},
		"notequal": func(b, a interface{}) bool {
			return b != a
		},
		"quote": func(b interface{}) string {
			switch bo := b.(type) {
			case string:
				return strconv.Quote(bo)
			case int:
				return strconv.Quote(strconv.Itoa(bo))
			case bool:
				return strconv.Quote(strconv.FormatBool(bo))
			case int64:
				return strconv.Quote(strconv.Itoa(int(bo)))
			case float32:
				mo := strconv.FormatFloat(float64(bo), 'f', 4, 32)
				return strconv.Quote(mo)
			case float64:
				mo := strconv.FormatFloat(bo, 'f', 4, 32)
				return strconv.Quote(mo)
			case byte:
				return strconv.QuoteRune(rune(bo))
			case rune:
				return strconv.QuoteRune(bo)
			default:
				return "Unconvertible Type"
			}
		},
		"prefixInt": func(prefix string, b int) string {
			return fmt.Sprintf("%s%d", prefix, b)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"multiply": func(a, b int) int {
			return a * b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"divide": func(a, b int) int {
			return a / b
		},
		"percentage": func(a, b float64) float64 {
			return (a / b) * 100
		},
	}
)

// ToTemplate returns a template instance with the giving templ string and functions.
func ToTemplate(name string, templ string, mx template.FuncMap) (*template.Template, error) {
	tml, err := template.New(name).Funcs(defaultFuncs).Funcs(mx).Parse(templ)
	if err != nil {
		return nil, err
	}

	return tml, nil
}
