package handler

import (
	"html/template"
	"net/http"
)

var Templates *template.Template
var spotifyClient *SpotifyClient

func ConnectedHandler(w http.ResponseWriter, r *http.Request) {
	Templates.ExecuteTemplate(w, "index.html", nil)
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	Templates.ExecuteTemplate(w, "welcome.html", nil)
}

type SpotifyClient struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectUri  string `yaml:"redirect_uri"`
}

func SetSpotifyClient(client *SpotifyClient) {
	spotifyClient = client
}
