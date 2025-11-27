package utils

import (
	"bytes"
	"html/template"
	"log"
	"path/filepath"
)

func LoadTemplate(templatePath string, args map[string]interface{}) (templateString string) {

	templatePath = "users/views/" + templatePath
	_, templateName := filepath.Split(templatePath)

	temp, err := template.New(templateName).Delims("${", "}").ParseFiles(templatePath)
	if err != nil {
		log.Println("Error in parse files:", err.Error())
		return
	}

	b := bytes.Buffer{}

	err = temp.Delims("${", "}").Execute(&b, args)
	if err != nil {
		log.Println("Error in execute template:", err.Error())
		return
	}

	templateString = b.String()
	return
}
