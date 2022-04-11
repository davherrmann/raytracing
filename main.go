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

	// world
	colors := []Color{{0.4, 0.8, 0.97}, {0.97, 0.55, 0.28}, {0.30, 0.89, 1}}
	world := generateWorld(colors)

	// server
	server := NewServer(world)

	log.Println("listening on port " + config.Port)
	http.ListenAndServe(":"+config.Port, server)
}
