package main

import (
	"math"
	"math/rand"
)

type MaterialHit struct {
	Attenuation Color
	Scattered   Ray
}

type Material func(ray Ray, hit Hit) *MaterialHit

func Lambertian(albedo Color) Material {
	return func(ray Ray, hit Hit) *MaterialHit {
		scatterDirection := hit.Normal.Add(randomUnitVector())

		// catch zero scatter direction
		if scatterDirection.Length() < 1e-8 {
			scatterDirection = hit.Normal
		}

		return &MaterialHit{
			Scattered: Ray{
				Origin:    hit.Point,
				Direction: scatterDirection,
			},
			Attenuation: albedo,
		}
	}
}

func Metal(albedo Color, fuzz float64) Material {
	return func(ray Ray, hit Hit) *MaterialHit {
		reflected := ray.Direction.Normalized().Reflect(hit.Normal)
		scattered := Ray{hit.Point, reflected.Add(randomUnitVector().Multiply(fuzz))}

		if reflected.Dot(hit.Normal) > 0 {
			return &MaterialHit{
				Scattered:   scattered,
				Attenuation: albedo,
			}
		}

		return nil
	}
}

func schlickReflectance(cosTheta, refractionRatio float64) float64 {
	// Schlick approximation for reflectance
	r0 := (1 - refractionRatio) / (1 * refractionRatio)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosTheta, 5)
}

func Dielectric(indexOfRefraction float64) Material {
	return func(ray Ray, hit Hit) *MaterialHit {
		refractionRatio := indexOfRefraction
		if hit.FrontFace {
			refractionRatio = 1 / indexOfRefraction
		}

		normalizedDirection := ray.Direction.Normalized()
		cosTheta := math.Min(normalizedDirection.Multiply(-1).Dot(hit.Normal), 1)
		sinTheta := math.Sqrt(1 - cosTheta*cosTheta)

		cannotRefract := refractionRatio*sinTheta > 1
		schlickReflect := schlickReflectance(cosTheta, refractionRatio) > rand.Float64()

		var scatterDirection Vec
		if cannotRefract || schlickReflect {
			scatterDirection = normalizedDirection.Reflect(hit.Normal)
		} else {
			scatterDirection = normalizedDirection.Refract(hit.Normal, refractionRatio)
		}

		scattered := Ray{hit.Point, scatterDirection}

		return &MaterialHit{
			Scattered:   scattered,
			Attenuation: Color{1, 1, 1},
		}
	}
}
