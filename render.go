package main

import (
	"image/color"
)

type Vec = Vector

type Ray struct {
	Origin    Vec
	Direction Vec
}

func (r *Ray) At(t float64) Vec {
	return r.Origin.Add(r.Direction.Multiply(t))
}

var world = World(
	Sphere(Vec{0, 0, -1}, 0.5),
	Sphere(Vec{0.3, 0, -1}, 0.4),
)

func rayColor(ray Ray) color.RGBA {
	hit := world(ray, 0, 10000)
	if hit != nil {
		t := hit.T
		if t > 0 {
			normal := ray.At(t).Subtract(Vec{0, 0, -1}).Normalized()
			return color.RGBA{
				R: uint8((normal.X + 1) * 128),
				G: uint8((normal.Y + 1) * 128),
				B: uint8((normal.Z + 1) * 128),
			}
		}
	}

	t := 0.5*ray.Direction.Normalized().Y + 1

	return color.RGBA{
		R: uint8(((1.0 - t) + t*0.5) * 0xff),
		G: uint8(((1.0 - t) + t*0.7) * 0xff),
		B: uint8(((1.0 - t) + t*1.0) * 0xff),
	}
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
			u := float64(x) / float64(width-1)
			v := float64(y) / float64(height-1)

			rayDirection := lowerLeftCorner. //
								Add(horizontal.Multiply(u)). //
								Add(vertical.Multiply(v)).   //
								Subtract(origin)

			ray := Ray{origin, rayDirection}
			color := rayColor(ray)
			drawFn(x, height-y, color)
		}
	}
}
