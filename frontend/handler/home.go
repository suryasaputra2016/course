package handler

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
)

type HomeHandler struct {
	Template *template.Template
}

func NewHomeHandler(tmpl *template.Template) *HomeHandler {
	return &HomeHandler{Template: tmpl}
}

func (hh HomeHandler) ShowHome(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	type Home struct {
		Title string
		Body  string
	}

	home := new(Home)

	json.Unmarshal(body, &home)

	hh.Template.ExecuteTemplate(w, "home.html", home)
}
