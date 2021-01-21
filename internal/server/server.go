package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()
	r.HandleFunc("/connected", connectedHandler).Methods("GET")
	r.HandleFunc("/disconnected", disconnectedHandler).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func connectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Connected!")
}

func disconnectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Disconnected!")
}
