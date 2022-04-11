package main

import "math/rand"

func generateWorld(colors []Color) Hittable {
	randomColor := func() Color {
		return colors[rand.Intn(len(colors)-1)]
	}
	materialGround := Lambertian(randomColor())
	materialCenter := Dielectric(1.5)
	materialLeft := Metal(randomColor(), 0.3)
	materialRight := Metal(randomColor(), 1.0)

	world := World(
		Sphere(Vec{0, -100.5, -1}, 100, materialGround),
		Sphere(Vec{0, 0.3, -1}, 0.5, materialCenter),
		Sphere(Vec{0, 0.3, -1}, -0.48, materialCenter),
		Sphere(Vec{-1, 0, -1}, 0.5, materialLeft),
		Sphere(Vec{1, 0, -1}, 0.5, materialRight),
	)

	return world
}
