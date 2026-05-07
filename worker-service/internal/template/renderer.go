package template

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"math/rand"
	"time"
)

//go:embed email/*.html
var emailFS embed.FS

type EmailData struct {
	Title   string
	Content string
}

var houses = []string{"stark", "lannister", "targaryen"}

var templates = make(map[string]*template.Template)

func init() {
	rand.Seed(time.Now().UnixNano())

	for _, house := range houses {
		fileName := fmt.Sprintf("email/%s.html", house)
		t, err := template.ParseFS(emailFS, fileName)
		if err != nil {
			panic(fmt.Sprintf("failed to parse template %s: %v", fileName, err))
		}
		templates[house] = t
	}
}

func RenderEmail(data EmailData) (string, error) {
	house := houses[rand.Intn(len(houses))]
	tmpl, ok := templates[house]
	if !ok {
		return "", fmt.Errorf("template for house %s not found", house)
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
