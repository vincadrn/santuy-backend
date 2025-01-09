package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hai gais"))
	})

	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":9000", nil))

	log.Println("Rian was here")
	log.Println("Angel was also here")
}
