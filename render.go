package main

import (
	"context"
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
	return Vec{rand.Float64()*2 - 1, rand.Float64()*2 - 1, rand.Float64()*2 - 1}.Normalized()
}

var samplesPerPixel = 10
var maxBounces = 10

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

func draw(ctx context.Context, width int, height int, angle float64, drawFn drawFn) {
	from := Vec{
		X: math.Cos(angle),
		Z: math.Sin(angle),
		Y: 0.5,
	}
	from = from.Normalized().Multiply(2)

	camera := Camera(Vec{0, 1, 0}, from, Vec{0, 0, -1}, width, height)

	colorSums := make([]Color, width*height)

	// TODO parallelize
	for s := 0; s < samplesPerPixel; s++ {
		for y := height - 1; y >= 0; y-- {
			for x := 0; x < width; x++ {
				if ctx.Err() != nil {
					return
				}

				u := (float64(x) + rand.Float64()) / float64(width-1)
				v := (float64(y) + rand.Float64()) / float64(height-1)

				ray := camera(u, v)
				singleColor := rayColor(ray, 0)

				i := y*width + x
				colorSums[i] = colorSums[i].Add(singleColor)

				if s%3 == 0 || s == samplesPerPixel-1 {
					// average color
					averageColor := colorSums[i].Multiply(1 / float64(s+1))

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
	}
}
