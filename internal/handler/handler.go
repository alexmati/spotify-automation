package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

var Templates *template.Template

func ConnectedHandler(w http.ResponseWriter, r *http.Request) {
	Templates.ExecuteTemplate(w, "index.html", nil)
}

func DisconnectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Disconnected!")
}
