package raytracing

import (
	"context"
	"image/color"
	"math"
	"math/rand"
)

type RenderOptions struct {
	ResolutionX int
	ResolutionY int

	SamplesPerPixel int
	MaxBounces      int
}

type RayCaster func(u, v float64) Ray

type Vec = Vector

type Ray struct {
	Origin    Vec
	Direction Vec
}

func (r *Ray) At(t float64) Vec {
	return r.Origin.Add(r.Direction.Multiply(t))
}

func randomUnitVector() Vec {
	return Vec{rand.Float64()*2 - 1, rand.Float64()*2 - 1, rand.Float64()*2 - 1}.Normalized()
}

var samplesPerPixel = 10
var maxBounces = 10

func rayColor(world Hittable, ray Ray, bounces int) Color {
	if bounces >= maxBounces {
		return Black
	}

	hit := world(ray, 0.001, 10000)
	if hit != nil {
		materialHit := hit.Material(ray, *hit)

		if materialHit == nil {
			return Black
		}

		return rayColor(world, materialHit.Scattered, bounces+1).Mix(materialHit.Attenuation)
	}

	t := 0.5 * (ray.Direction.Normalized().Y + 1)

	return Color{1, 1, 1}.Multiply(1 - t).Add(Color{0.5, 0.7, 1.0}.Multiply(t))
}

type drawFn func(x, y int, color color.RGBA)

func Render(ctx context.Context, world Hittable, camera Camera, options RenderOptions, drawFn drawFn) {
	aspectRatio := float64(options.ResolutionX) / float64(options.ResolutionY)
	rayCaster := camera.RayCaster(aspectRatio)

	width := options.ResolutionX
	height := options.ResolutionY
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

				ray := rayCaster(u, v)
				singleColor := rayColor(world, ray, 0)

				i := y*width + x
				colorSums[i] = colorSums[i].Add(singleColor)

				// send first, every nth and last sample
				if s == 0 || s%3 == 0 || s == samplesPerPixel-1 {
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
