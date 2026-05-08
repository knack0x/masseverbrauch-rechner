package main

import (
	"embed"
	"fmt"
	"html/template"
	"strings"
)

//go:embed templates
var templateFiles embed.FS

func tmplSlots(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i + 1
	}
	return s
}

func tmplSub(a, b int) int {
	return a - b
}

func tmplNumberFormat(v float64, decimals int) string {
	format := fmt.Sprintf("%%.%df", decimals)
	s := fmt.Sprintf(format, v)
	parts := strings.Split(s, ".")
	intPart := parts[0]
	var result []byte
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, byte(c))
	}
	if len(parts) > 1 {
		result = append(result, ',')
		result = append(result, parts[1]...)
	}
	return string(result)
}

func tmplTrSlotName(name string) string {
	tr := map[string]string{
		"Tower Slot 1": "Turmposition 1",
		"Tower Slot 2": "Turmposition 2",
		"Tower Slot 3": "Turmposition 3",
		"Tower Slot 4": "Turmposition 4",
		"Tower Slot 5": "Turmposition 5",
	}
	if v, ok := tr[name]; ok {
		return v
	}
	return name
}

var tmpls *template.Template

func initTemplates() error {
	funcMap := template.FuncMap{
		"slots":        tmplSlots,
		"sub":          tmplSub,
		"numberFormat": tmplNumberFormat,
		"trSlotName":   tmplTrSlotName,
	}
	var err error
	tmpls, err = template.New("").Funcs(funcMap).ParseFS(templateFiles, "templates/*.html")
	return err
}
