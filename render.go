package main

import (
	"image/color"
	"math"
	"math/rand"
)

type Vec = Vector

type Color struct {
	R float64
	G float64
	B float64
}

func (cA Color) Add(cB Color) Color {
	return Color{
		R: cA.R + cB.R,
		G: cA.G + cB.G,
		B: cA.B + cB.B,
	}
}

func (cA Color) Mix(cB Color) Color {
	return Color{
		R: cA.R * cB.R,
		G: cA.G * cB.G,
		B: cA.B * cB.B,
	}
}

func (c Color) Multiply(factor float64) Color {
	return Color{
		R: c.R * factor,
		G: c.G * factor,
		B: c.B * factor,
	}
}

type Ray struct {
	Origin    Vec
	Direction Vec
}

func (r *Ray) At(t float64) Vec {
	return r.Origin.Add(r.Direction.Multiply(t))
}

var world = World(
	Sphere(Vec{0, 0, -1}, 0.5, DiffuseColor{0.5, 0.7, 1.0}),
	Sphere(Vec{0.3, 0, -1}, 0.4, DiffuseColor{0.2, 0.8, 0.3}),
	Sphere(Vec{0, -100.5, -1}, 100, DiffuseColor{0.95, 0.1, 0.1}),
)

func randomInUnitSphere() Vec {
	return Vec{rand.Float64(), rand.Float64(), rand.Float64()}.Normalized()
}

var samplesPerPixel = 100
var maxBounces = 50

func rayColor(ray Ray, bounces int) Color {
	if bounces >= maxBounces {
		return Color{} // black
	}

	hit := world(ray, 0.001, 10000)
	if hit != nil {
		target := hit.Point.Add(hit.Normal).Add(randomInUnitSphere())
		nextRay := Ray{hit.Point, target.Subtract(hit.Point)}
		rayColor := rayColor(nextRay, bounces+1)

		switch material := hit.Material.(type) {
		case DiffuseColor:
			rayColor = rayColor.Mix(Color(material))
		default:
		}

		return Color{
			R: rayColor.R / 2,
			G: rayColor.G / 2,
			B: rayColor.B / 2,
		}
	}

	t := 0.5 * (ray.Direction.Normalized().Y + 1)

	return Color{1, 1, 1}.Multiply(1 - t).Add(Color{0.2, 0.5, 0.7}.Multiply(t))
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
			gamma := 2.2
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
