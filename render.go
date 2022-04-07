package main

import (
	"image/color"
	"math"
	"math/rand"
)

type Vec = Vector

type Ray struct {
	Origin    Vec
	Direction Vec
}

func (r *Ray) At(t float64) Vec {
	return r.Origin.Add(r.Direction.Multiply(t))
}

var (
	materialGround = Lambertian(Color{0.8, 0.8, 0.0})
	materialCenter = Dielectric(1.5)
	materialLeft   = Metal(Color{0.8, 0.8, 0.8}, 0.3)
	materialRight  = Metal(Color{0.8, 0.6, 0.2}, 1.0)
)

var world = World(
	Sphere(Vec{0, -100.5, -1}, 100, materialGround),
	Sphere(Vec{0, 0.3, -1}, 0.5, materialCenter),
	Sphere(Vec{0, 0.3, -1}, -0.48, materialCenter),
	Sphere(Vec{-1, 0, -1}, 0.5, materialLeft),
	Sphere(Vec{1, 0, -1}, 0.5, materialRight),
)

func randomUnitVector() Vec {
	return Vec{rand.Float64(), rand.Float64(), rand.Float64()}.Normalized()
}

var samplesPerPixel = 100
var maxBounces = 100

func rayColor(ray Ray, bounces int) Color {
	if bounces >= maxBounces {
		return Black
	}

	hit := world(ray, 0.001, 10000)
	if hit != nil {
		materialHit := hit.Material(ray, *hit)

		if materialHit == nil {
			return Black
		}

		return rayColor(materialHit.Scattered, bounces+1).Mix(materialHit.Attenuation)
	}

	t := 0.5 * (ray.Direction.Normalized().Y + 1)

	return Color{1, 1, 1}.Multiply(1 - t).Add(Color{0.5, 0.7, 1.0}.Multiply(t))
}

type drawFn func(x, y int, color color.RGBA)

func draw(width int, height int, drawFn drawFn) {
	aspectRatio := float64(width) / float64(height)

	// camera & viewport
	viewportHeight := 2.0
	viewportWidth := aspectRatio * viewportHeight
	focalLength := 1.0

	// ray origin
	origin := Vec{0, 0, 0}
	horizontal := Vec{viewportWidth, 0, 0}
	vertical := Vec{0, viewportHeight, 0}
	lowerLeftCorner := origin.Add(horizontal.Multiply(-1. / 2)).Add(vertical.Multiply(-1. / 2)).Add(Vec{0, 0, -focalLength})

	// TODO parallelize
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			colorSum := Color{}
			for s := 0; s < samplesPerPixel; s++ {
				u := (float64(x) + rand.Float64()) / float64(width-1)
				v := (float64(y) + rand.Float64()) / float64(height-1)

				rayDirection := lowerLeftCorner. //
									Add(horizontal.Multiply(u)). //
									Add(vertical.Multiply(v)).   //
									Subtract(origin)

				ray := Ray{origin, rayDirection}
				color := rayColor(ray, 0)
				colorSum = colorSum.Add(color)
			}

			// average color
			averageColor := colorSum.Multiply(1 / float64(samplesPerPixel))

			// gamma correction
			gamma := 2.0
			gammaCorrected := Color{
				R: math.Pow(averageColor.R, 1/gamma),
				G: math.Pow(averageColor.G, 1/gamma),
				B: math.Pow(averageColor.B, 1/gamma),
			}

			// convert color
			converted := color.RGBA{
				R: uint8(gammaCorrected.R * 0xff),
				G: uint8(gammaCorrected.G * 0xff),
				B: uint8(gammaCorrected.B * 0xff),
			}

			drawFn(x, height-y, converted)
		}
	}
}
