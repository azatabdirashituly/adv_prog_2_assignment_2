package cmd

import (
	"html/template"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := renderTemplate(w, "home.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	t, err := template.ParseFiles("frontend/templates/" + tmpl)
	if err != nil {
		return err
	}
	err = t.Execute(w, data) 
	if err != nil {
		return err
	}
	return nil
}