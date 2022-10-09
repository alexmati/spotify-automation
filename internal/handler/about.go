package handler

import (
	"fmt"
	"net/http"
)

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	accessToken, err := getAccessToken(r)
	if err != nil {
		http.Redirect(w, r, "http://localhost:8080/", 301)
		fmt.Printf("could not get user details: %v", err)
		return
	}

	details, err := GetUserDetails(accessToken)
	if err != nil {
		fmt.Printf("could not get user details: %v", err)
		return
	}

	m := map[string]interface{}{
		"DisplayName": details.DisplayName,
		"UserID":      details.UserID,
	}

	if err := Templates.ExecuteTemplate(w, "about.html", m); err != nil {
		fmt.Printf("could not execute about template: %v", err)
		return
	}
}
