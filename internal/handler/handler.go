package handler

import (
	"html/template"
)

var Templates *template.Template
var spotifyClient *SpotifyClient

type SpotifyClient struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectUri  string `yaml:"redirect_uri"`
}

func SetSpotifyClient(client *SpotifyClient) {
	spotifyClient = client
}
