package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type SpotifyClient struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectUri  string `yaml:"redirect_uri"`
}

type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	Expiry       int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (c SpotifyClient) getAccessToken(code string) (*SpotifyTokenResponse, error) {
	urlForm := url.Values{}
	urlForm.Set("grant_type", "authorization_code")
	urlForm.Set("code", code)
	urlForm.Set("redirect_uri", spotifyClient.RedirectUri)
	body := strings.NewReader(urlForm.Encode())

	endpoint := "https://accounts.spotify.com/api/token"
	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		fmt.Printf("could not make new request: %v", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req.SetBasicAuth(spotifyClient.ClientId, spotifyClient.ClientSecret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("could not set basic auth: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response body: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("failed to get spotify access token %v: \n%v", resp.StatusCode, string(resBody))
		return nil, err
	}

	tokenRes := &SpotifyTokenResponse{}
	if err := json.Unmarshal(resBody, tokenRes); err != nil {
		fmt.Printf("failed to unmarshal token response: %v", err)
		return nil, err
	}
	return tokenRes, nil
}

//---GET USERS PLAYLIST---//

type Image struct {
	PlaylistImage string `json:"url"`
}

type Owner struct {
	UserID string `json:"id"`
}

type Track struct {
	Total int `json:"total"`
}

type Playlist struct {
	SpotifyID    string  `json:"id"`
	Images       []Image `json:"images"`
	PlaylistName string  `json:"name"`
	Owner        Owner   `json:"owner"`
	Tracks       Track   `json:"tracks"`
}

type PlaylistResponse struct {
	Playlists []*Playlist `json:"items"`
}

func GetPlaylists(accessToken string) ([]*Playlist, error) {
	endpoint := "https://api.spotify.com/v1/me/playlists?limit=50"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		fmt.Printf("failed to create playlists request: %v", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("failed to perform get playlists request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("expected http status code when getting playlists: %v", resp.StatusCode)
		return nil, err
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read playlists response: %v", err)
		return nil, err
	}

	playlistRes := &PlaylistResponse{}
	if err := json.Unmarshal(resBody, playlistRes); err != nil {
		fmt.Printf("failed to unmarshal playlist response: %v", err)
		return nil, err
	}
	return playlistRes.Playlists, nil
}

//---GET USERS DETAILS---//

type UserResponse struct {
	DisplayName string `json:"display_name"`
	UserID      string `json:"id"`
}

func GetUserDetails(accessToken string) (*UserResponse, error) {
	endpoint := "https://api.spotify.com/v1/me"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		fmt.Printf("failed to create current user request: %v", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("failed to perform get current user request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("expected http status code when getting current user: %v", resp.StatusCode)
		return nil, err
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read playlists response: %v", err)
		return nil, err
	}

	userRes := &UserResponse{}
	if err := json.Unmarshal(resBody, userRes); err != nil {
		fmt.Printf("failed to unmarshal user response: %v", err)
		return nil, err
	}
	return userRes, nil
}

//---CREATE A NEW PLAYLIST---//

type CreatePlaylistRequest struct {
	PlaylistName string `json:"name"`
	Description  string `json:"description"`
	Public       bool   `json:"public"`
}

func (c CreatePlaylistRequest) Read(p []byte) (n int, err error) {
	panic("todo")
}

func CreateUserPlaylist(accessToken string, playlistRequest *CreatePlaylistRequest, UserID string) {

	endpoint := fmt.Sprintf("https://api.spotify.com/v1/users/%v/playlists", UserID)
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(playlistRequest.PlaylistName))
	if err != nil {
		fmt.Printf("failed to make create playlist request: %v", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("failed to perform create new playlist: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("expected http status code when trying to create a new playlist: %v", resp.StatusCode)
		return
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read playlists response: %v", err)
		return
	}
	fmt.Println(string(resBody))
}

func GetPlaylistDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Printf("incorrect method request")
	}

	playlistRequest := CreatePlaylistRequest{
		PlaylistName: r.FormValue("playlist-name"),
		Description:  r.FormValue("description"),
		Public:       r.FormValue("public") == "",
	}
	fmt.Println(playlistRequest)
	_ = w
}

//---GET USERS TOP SONGS---//

type TopSong struct {
	Request  string `json:"href"`
	SongName string `json:"name"`
}
type TopSongsResponse struct {
	TopSongs []*TopSong `json:"items"`
}

func GetTopSongs(accessToken string) ([]*TopSong, error) {
	endpoint := "https://api.spotify.com/v1/me/top/tracks?limit=30"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		fmt.Printf("failed to obtain users top tracks: %v", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("failed to perform get top songs request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("expected http status code when getting top songs: %v ", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %v ", resp.StatusCode)
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read top songs response: %v ", err)
		return nil, err
	}

	topSongsRes := &TopSongsResponse{}
	if err := json.Unmarshal(resBody, topSongsRes); err != nil {
		fmt.Printf("failed to unmarshal top songs response: %v ", err)
		return nil, err
	}

	return topSongsRes.TopSongs, nil
}
