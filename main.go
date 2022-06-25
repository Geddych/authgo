package main

import (
	sessions "authgo/Sessions"
	"authgo/tokens"
	"authgo/user"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		var result map[string]interface{}
		json.NewDecoder(r.Body).Decode(&result)
		uid := user.FindUserIDbyUsername(fmt.Sprint(result["username"]))
		if uid == mongo.ErrNoDocuments {
			w.Write([]byte("No user's with that username!"))
		} else {
			at := tokens.CreateAccessToken(uid)
			bat := tokens.BeareringAccessToken(at)
			rt := tokens.CreateRefreshToken(at)
			sessions.CreateSession(tokens.HashToken(rt), tokens.GetIdFromToken(rt))
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Authorization", bat)
			w.Header().Add("Refresh", rt)
		}
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Values("Authorization")
		bearerHeader := authHeader[0]
		if authHeader != nil {
			if strings.Contains(bearerHeader, "Bearer") {
				isTrue := tokens.CheckTokensLifetime(tokens.DeBeareringToken(bearerHeader))
				if isTrue == true {
					w.Write([]byte("Token is right"))
				} else {
					uid := tokens.GetIdFromToken(tokens.DeBeareringToken(bearerHeader))
					at := tokens.CreateAccessToken(uid)
					bat := tokens.BeareringAccessToken(at)
					rt := tokens.CreateRefreshToken(at)
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("Authorization", bat)
					w.Header().Add("Refresh", rt)
				}
			} else {
				fmt.Print("False")
			}
		} else {
			fmt.Print(authHeader)
		}
	})

	http.ListenAndServe(":3000", r)
}
