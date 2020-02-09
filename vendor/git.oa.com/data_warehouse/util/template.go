package util

import (
	"html/template"
	"path/filepath"
)

func NewTemplate(tplFileName string) (t *template.Template, err error) {
	_, baseName := filepath.Split(tplFileName)
	t, err = template.New(baseName).Funcs(template.FuncMap{
		"number_format": NumberFormat,
		"render_color":  RenderColor,
		"html": func(x string) interface{} {
			return template.HTML(x)
		},
		"add": func(nums ...interface{}) float64 {
			n := float64(0)
			for _, val := range nums {
				n += AnyToFloat(val)
			}
			return n
		},
		"minus": func(nums ...interface{}) float64 {
			n := AnyToFloat(nums[0])
			for _, val := range nums[1:] {
				n -= AnyToFloat(val)
			}
			return n
		},
	}).ParseFiles(tplFileName)
	return
}
