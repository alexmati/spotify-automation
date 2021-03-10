package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]interface{}{
		"ClientId":    spotifyClient.ClientId,
		"RedirectUri": spotifyClient.RedirectUri,
	}

	if err := Templates.ExecuteTemplate(w, "index.html", m); err != nil {
		fmt.Printf("Could not execute template: %v", err)
		return
	}
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")

	urlForm := url.Values{}
	urlForm.Set("grant_type", "authorization_code")
	urlForm.Set("code", code)
	urlForm.Set("redirect_uri", spotifyClient.RedirectUri)
	body := strings.NewReader(urlForm.Encode())

	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", body)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req.SetBasicAuth(spotifyClient.ClientId, spotifyClient.ClientSecret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Could not set basic auth: %v", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to get spotify access token %v: \n%v", resp.StatusCode, string(resBody))
		return
	}

	type SpotifyTokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		Expiry       int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	tokenRes := &SpotifyTokenResponse{}
	if err := json.Unmarshal(resBody, tokenRes); err != nil {
		fmt.Printf("Failed to unmarshal: %v", err)
		return
	}

	const (
		cookieName = "spotify_access"
	)
	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: base64.StdEncoding.EncodeToString(resBody),
	})

	req, err = http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me/playlists", nil)
	if err != nil {
		fmt.Printf("Failed to create playlists request: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+tokenRes.AccessToken)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Failed to perform get playlists request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Enexpected http status code when getting playlists: %v", resp.StatusCode)
		return
	}

	resBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read playlists response: %v", err)
		return
	}

	fmt.Println(string(resBody))
	Templates.ExecuteTemplate(w, "welcome.html", r)
}
