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
	colors := []Color{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
	world := generateWorld(colors)

	// server
	server := NewServer(world)

	log.Println("listening on port " + config.Port)
	http.ListenAndServe(":"+config.Port, server)
}
