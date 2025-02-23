package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"vincadrn.com/santuy/internal/auth"
	"vincadrn.com/santuy/internal/handler"
	"vincadrn.com/santuy/internal/middleware"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=verify-full", user, pass, host, port, dbName)
	db, err := sql.Open("postgres", connStr)
	defer db.Close()

	if err != nil {
		log.Fatal("Cannot connect to database.")
	}

	ctx := context.Background()

	mux := http.NewServeMux()

	mux.Handle("/v1/auth/login", auth.RequestAuth())
	mux.Handle("/v1/auth/session", auth.RequestSession())

	mux.Handle("/v1/users/{userId}", handler.GetUserHandler())
	mux.Handle("/v1/groups", handler.ListGroups())
	mux.Handle("/v1/groups/{groupId}", handler.GroupHandler())

	mux.Handle("/v1/itineraries", handler.ListItineraries())
	mux.Handle("/v1/itineraries/{itineraryId}", handler.GetItinerary(db, ctx))

	var handler http.Handler = mux
	handler = middleware.AuthMiddleware(handler)
	handler = middleware.CORSMiddleware(handler)

	server := new(http.Server)
	server.Addr = ":9000"
	server.Handler = handler

	log.Println("Server started at " + server.Addr)
	server.ListenAndServe()
}
