package moz

import "text/template"

// ToTemplate returns a template instance with the giving templ string and functions.
func ToTemplate(name string, templ string, mx template.FuncMap) (*template.Template, error) {
	tml, err := template.New(name).Funcs(mx).Parse(templ)
	if err != nil {
		return nil, err
	}

	return tml, nil
}
