package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	authCookie = "token"
)

type JwtSecret struct {
	Secret string `yaml:"secret"`
}

var (
	jwtKey    = []byte("eyfDdDfgSDh2gj4f") //TODO: autogenerate key and store in config
	jwtExpiry = time.Hour * 5
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]interface{}{
		"ClientId":    spotifyClient.ClientId,
		"RedirectUri": spotifyClient.RedirectUri,
	}

	if err := Templates.ExecuteTemplate(w, "index.html", m); err != nil {
		fmt.Printf("could not execute index template: %v", err)
		return
	}
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")
	if code == "" {
		fmt.Printf("invalid authorization code")
		return
	}

	token, err := spotifyClient.getAccessToken(code)
	if err != nil {
		fmt.Printf("could not get access token: %v", err)
		return
	}

	if err := setTokenCookie(w, token.AccessToken); err != nil {
		fmt.Printf("failed to set token cookie: %v", err)
		return
	}

	http.Redirect(w, r, "/about", http.StatusPermanentRedirect)
	return
}

// 1. create your jwt payload/claims
// 2. sign the payload with jwt key
func signJwtToken(accessToken string) (string, error) {
	c := &Claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(jwtExpiry).Unix(),
		},
		AccessToken: accessToken,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString(jwtKey)
}

// header.payload.signature

func setTokenCookie(w http.ResponseWriter, accessToken string) error {
	token, err := signJwtToken(accessToken)
	if err != nil {
		return fmt.Errorf("failed to sign jwt token: %v", err)
	}

	//add CSRF protection (auth cookies for state-changing actions is a security flaw)
	http.SetCookie(w, &http.Cookie{
		Name:     authCookie,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(jwtExpiry),
		HttpOnly: true,
	})
	return nil
}
