package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type ctxKey int

const (
	ctxAccessToken ctxKey = iota
)

func setAccessToken(r *http.Request, accessToken string) context.Context {
	return context.WithValue(r.Context(), ctxAccessToken, accessToken) // contexts are immutable, so we return a new one with the access token
}

func getAccessToken(r *http.Request) (accessToken string, err error) {
	accessToken, ok := r.Context().Value(ctxAccessToken).(string) // ctx.Value(..) returns interface{} but we want string, so we do "type assertion" -> .(type)
	if !ok {
		return "", errors.New("couldn't find access token in context")
	}
	return accessToken, nil
}

type Claims struct {
	jwt.StandardClaims
	AccessToken string `json:"access_token"`
}

func verifyJWT(r *http.Request, key []byte) (accessToken string, err error) {
	cookie, err := r.Cookie(authCookie)
	if err != nil {
		return "", fmt.Errorf("failed to get cookie: %v", err)
	}
	token := cookie.Value

	signerFunc := func(_jwt *jwt.Token) (interface{}, error) {
		if _, ok := _jwt.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("received an unexpected signing method for token")
		}
		return key, nil
	}

	c := &Claims{}
	t, err := jwt.ParseWithClaims(token, c, signerFunc)
	if err != nil {
		return "", fmt.Errorf("failed to parse jwt claims: %v", err)
	}

	if !t.Valid {
		return "", errors.New("invalid cookie")
	}
	return c.AccessToken, nil
}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			accessToken, err := verifyJWT(r, jwtKey)
			if err != nil {
				fmt.Printf("failed to verify jwt: %v", err)
				http.Redirect(w, r, "http://localhost:8080/", 302)
				return
			}

			ctx := setAccessToken(r, accessToken)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
}
