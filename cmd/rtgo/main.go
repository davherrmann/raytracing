package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/davherrmann/rtgo"
	"github.com/davherrmann/rtgo/raytracing"
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
	randomColorPalette := []raytracing.Color{
		{R: 0.4, G: 0.8, B: 0.97},
		{R: 0.97, G: 0.55, B: 0.28},
		{R: 0.30, G: 0.89, B: 1},
	}
	camera := rtgo.GenerateCamera(0, 1, 400, 300)
	world := rtgo.GenerateWorld(randomColorPalette)

	// server
	server := rtgo.NewServer(camera, world)

	log.Println("listening on port " + config.Port)
	http.ListenAndServe(":"+config.Port, server)
}
