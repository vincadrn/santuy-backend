package main

import (
	"log"
	"net/http"

	"vincadrn.com/santuy/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/api/v1/user/{userId}", handler.GetUserHandler())
	mux.Handle("/api/v1/group/{groupId}", handler.GetGroupHandler())

	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":9000", mux))
}
