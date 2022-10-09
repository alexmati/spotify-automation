package handler

import (
	"fmt"
	"net/http"
)

func TopSongsHandler(w http.ResponseWriter, r *http.Request) {
	accessToken, err := getAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "http://localhost:8080/", 301)
		fmt.Printf("could not get user details: %v", err)
		return
	}

	songs, err := GetTopSongs(accessToken)
	if err != nil {
		fmt.Printf("could not get users top songs: %v", err)
		return
	}
	details, err := GetUserDetails(accessToken)
	if err != nil {
		fmt.Printf("could not get user details: %v", err)
		return
	}

	m := map[string]interface{}{
		"SongName":    songs,
		"DisplayName": details.DisplayName,
	}
	for _, m := range songs {
		fmt.Println(m.Request)
		fmt.Println(m.SongName)
	}

	if err := Templates.ExecuteTemplate(w, "songs.html", m); err != nil {
		fmt.Printf("could not execute top songs template: %v", err)
		return
	}
}
