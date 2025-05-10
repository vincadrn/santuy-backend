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
	"vincadrn.com/santuy/internal/repository"
	"vincadrn.com/santuy/internal/service"
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

	// Initialize repo and service
	accountRepo := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepo)
	itineraryRepo := repository.NewItineraryRepository(db)
	itineraryService := service.NewItineraryService(itineraryRepo)

	ctx := context.Background()

	mux := http.NewServeMux()

	mux.Handle("/v1/auth/login", auth.RequestAuth(accountService, ctx))
	mux.Handle("/v1/auth/session", auth.RequestSession(accountService, ctx))
	mux.Handle("/v1/auth/logout", auth.RequestLogout())

	// User & group
	mux.Handle("/v1/user", handler.UserHandler(accountService, ctx))
	mux.Handle("/v1/group", handler.GroupHandler(accountService, ctx))

	// Itinerary
	mux.Handle("/v1/itineraries", handler.ItineraryHandler(itineraryService, ctx))
	mux.Handle("/v1/itineraries/{itineraryId}/details", handler.ItineraryDetailHandler(itineraryService, ctx))

	mux.Handle("/v1/picture/1", handler.GetPicture(db, ctx))

	mux.Handle("/v1/logout", auth.RequestLogout())

	var handler http.Handler = mux
	handler = middleware.AuthMiddleware(handler)
	handler = middleware.CORSMiddleware(handler)

	server := new(http.Server)
	server.Addr = ":9000"
	server.Handler = handler

	log.Println("Server started at " + server.Addr)
	server.ListenAndServe()
}
