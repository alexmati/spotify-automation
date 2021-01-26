package handler

import (
	"html/template"
	"net/http"
)

var Templates *template.Template
var spotifyClient *SpotifyClient

func ConnectedHandler(w http.ResponseWriter, r *http.Request) {
	if err := Templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	if err := Templates.ExecuteTemplate(w, "welcome.html", nil); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

type SpotifyClient struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectUri  string `yaml:"redirect_uri"`
}

func SetSpotifyClient(client *SpotifyClient) {
	spotifyClient = client
}
