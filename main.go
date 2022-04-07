package main

import (
	"log"
	"net/http"
)

func main() {
	server := NewServer()

	port := "8080"

	log.Println("listening on port " + port)
	http.ListenAndServe(":"+port, server)
}
