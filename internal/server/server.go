package server

import (
	"io"
	"log"
	"net/http"
)

func Run() {
	con := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Connected\n")
	}

	http.HandleFunc("/", con)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
