package main

import (
	"log"
	"net/http"

	"vincadrn.com/santuy/internal/auth"
	"vincadrn.com/santuy/internal/handler"
	"vincadrn.com/santuy/internal/middleware"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mux := http.NewServeMux()

	mux.Handle("/v1/auth/login", auth.RequestAuth())
	mux.Handle("/v1/auth/session", auth.RequestSession())

	mux.Handle("/v1/users/{userId}", handler.GetUserHandler())
	mux.Handle("/v1/groups/{groupId}", handler.GetGroupHandler())

	// assume for now that these following data:
	// 1. user info (name, email, etc)
	// 2. current active group and role in that group
	// can be retrieved via session.
	// TODO: Create a middleware to retrieve the key-value info
	// from the session
	// mux.Handle("/v1/itineraries", nil)
	// mux.Handle("/v1/itineraries/{itineraryId}", nil)

	var handler http.Handler = mux
	handler = middleware.AuthMiddleware(handler)
	handler = middleware.CORSMiddleware(handler)

	server := new(http.Server)
	server.Addr = ":9000"
	server.Handler = handler

	log.Println("Server started at " + server.Addr)
	server.ListenAndServe()
}
