package main

import (
	"flag"
	"log"
	"net/http"
)

type Config struct {
	Port string
}

func main() {
	config := Config{}

	// flags
	flag.StringVar(&config.Port, "port", "8080", "listen port for server")
	flag.Parse()

	log.Println("listening on port " + config.Port)

	server := NewServer()
	http.ListenAndServe(":"+config.Port, server)
}
