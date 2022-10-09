package handler

import (
	"html/template"
)

var Templates *template.Template
var spotifyClient *SpotifyClient

func SetSpotifyClient(client *SpotifyClient) {
	spotifyClient = client
}
