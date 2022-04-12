package main

import (
	"flag"

	"github.com/davherrmann/rtgo"
	"github.com/davherrmann/rtgo/raytracing"
)

type Config struct {
	Port string
}

func parseConfig() Config {
	config := Config{}

	flag.StringVar(&config.Port, "port", "8080", "listen port for server")
	flag.Parse()

	return config
}

func main() {
	config := parseConfig()

	// world & camera
	randomColorPalette := []raytracing.Color{
		{R: 0.4, G: 0.8, B: 0.97},
		{R: 0.97, G: 0.55, B: 0.28},
		{R: 0.30, G: 0.89, B: 1},
	}
	camera := rtgo.GenerateCamera(0, 1, 400, 300)
	world := rtgo.GenerateWorld(randomColorPalette)

	// http server
	handler := rtgo.NewServer(camera, world)
	handler.ListenAndServe(config.Port)
}
