package rtgo

import (
	"math"
	"math/rand"

	"github.com/davherrmann/rtgo/raytracing"
)

type Vec = raytracing.Vector

func GenerateCamera(angle float64, zoom float64, width, height int) raytracing.Camera {
	from := Vec{
		X: math.Cos(angle),
		Z: math.Sin(angle),
		Y: 0.5,
	}
	from = from.Normalized().Multiply(10)

	camera := raytracing.Camera{
		Up:     Vec{X: 0, Y: 1, Z: 0},
		From:   from,
		LookAt: Vec{X: 0, Y: 0, Z: 0},
		Zoom:   zoom,
	}
	return camera
}

func GenerateWorld(colors []raytracing.Color) raytracing.Hittable {
	randomColor := func() raytracing.Color {
		return colors[rand.Intn(len(colors)-1)]
	}
	materialGround := raytracing.Lambertian(randomColor())
	materialCenter := raytracing.Dielectric(1.5)
	materialLeft := raytracing.Metal(randomColor(), 0.3)
	materialRight := raytracing.Metal(randomColor(), 1.0)

	world := raytracing.World(
		raytracing.Sphere(Vec{X: 0, Y: -100.5, Z: -1}, 100, materialGround),
		raytracing.Sphere(Vec{X: 0, Y: 0.3, Z: -1}, 0.5, materialCenter),
		raytracing.Sphere(Vec{X: 0, Y: 0.3, Z: -1}, -0.48, materialCenter),
		raytracing.Sphere(Vec{X: -1, Y: 0, Z: -1}, 0.5, materialLeft),
		raytracing.Sphere(Vec{X: 1, Y: 0, Z: -1}, 0.5, materialRight),
	)

	return world
}
