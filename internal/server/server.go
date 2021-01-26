package server

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/alexmati/spotify-automation/internal/handler"

	"github.com/gorilla/mux"
)

func Run() {
	handler.Templates = template.Must(template.ParseGlob("internal/templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", handler.ConnectedHandler).Methods("GET")
	r.HandleFunc("/welcome", handler.WelcomeHandler).Methods("GET")

	cfg, err := getConfig()
	if err != nil {
		log.Fatal("error")
	}
	handler.SetSpotifyClient(cfg.SpotifyClient)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

type Config struct {
	SpotifyClient *handler.SpotifyClient `yaml:"spotify_client"`
}

func getConfig() (*Config, error) {
	var cfg Config

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Cannot find file: %v", err)
	}
	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		log.Fatalf("Cannot unmarshall file: %v", err)
	}
	return &cfg, err
}
