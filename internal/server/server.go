package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/alexmati/spotify-automation/internal/handler"

	"github.com/gorilla/mux"
)

func Run() {
	handler.Templates = template.Must(template.ParseGlob("internal/templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", handler.ConnectedHandler).Methods("GET")
	r.HandleFunc("/disconnected", handler.DisconnectedHandler).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
