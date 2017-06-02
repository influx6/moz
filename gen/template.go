package gen

import (
	"strconv"
	"text/template"
)

var (
	defaultFuncs = template.FuncMap{
		"greaterThanF": func(b, a float64) bool {
			return b > a
		},
		"lessThanF": func(b, a float64) bool {
			return b < a
		},
		"quote": func(b string) string {
			return strconv.Quote(b)
		},
		"greaterThan": func(b, a int) bool {
			return b > a
		},
		"lessThan": func(b, a int) bool {
			return b < a
		},
		"equal": func(b, a interface{}) bool {
			return b == a
		},
		"notequal": func(b, a interface{}) bool {
			return b != a
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
