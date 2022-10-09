package handler

import (
	"fmt"
	"net/http"
)

func SelectPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	accessToken, err := getAccessToken(r) // gets the access token from the context that was set in the Authentication middleware
	if err != nil {
		http.Redirect(w, r, "http://localhost:8080/", 301)
		fmt.Printf("failed to get access token: %v", err)
		return
	}

	playlists, err := GetPlaylists(accessToken)
	if err != nil {
		fmt.Printf("failed to get playlists: %v", err)
		return
	}
	details, err := GetUserDetails(accessToken)
	if err != nil {
		fmt.Printf("could not get user details: %v", err)
		return
	}

	m := map[string]interface{}{
		"PlaylistName": playlists,
		"DisplayName":  details.DisplayName,
	}

	if err := Templates.ExecuteTemplate(w, "playlist.html", m); err != nil {
		fmt.Printf("could not execute select playlist template: %v", err)
		return
	}
}
