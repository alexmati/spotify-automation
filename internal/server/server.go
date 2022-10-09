package server

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alexmati/spotify-automation/internal/handler"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

func Run() {
	handler.Templates = template.Must(template.ParseGlob("internal/templates/*.html"))

	r := mux.NewRouter()
	r.HandleFunc("/", handler.LoginHandler).Methods(http.MethodGet)
	r.HandleFunc("/callback", handler.CallbackHandler).Methods(http.MethodGet)

	authRouter := r.NewRoute().Subrouter()
	authRouter.Use(handler.Authentication)
	authRouter.HandleFunc("/about", handler.AboutHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/create", handler.CreatePlaylistHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/playlist", handler.SelectPlaylistHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/songs", handler.TopSongsHandler).Methods(http.MethodGet)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	cfg, err := getConfig()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	handler.SetSpotifyClient(cfg.SpotifyClient)

	http.Handle("/", r)
	fmt.Println("Serving on PORT 8080")
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
