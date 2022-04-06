package main

import "math"

type Hit struct {
	Point    Vec
	Normal   Vec
	T        float64
	Material Material
}

type Hittable func(ray Ray, tMin, tMax float64) *Hit

type Material interface{}

type DiffuseColor Color

func World(hittables ...Hittable) Hittable {
	return func(ray Ray, tMin, tMax float64) *Hit {
		var closestHit *Hit
		closest := tMax

		for _, hittable := range hittables {
			hit := hittable(ray, tMin, tMax)
			if hit == nil {
				continue
			}

			distance := hit.Point.Subtract(ray.Origin).Length()
			if distance < closest {
				closestHit = hit
				closest = distance
			}
		}

		return closestHit
	}
}

func Sphere(center Vec, radius float64, material Material) Hittable {
	return func(ray Ray, tMin, tMax float64) *Hit {
		oc := ray.Origin.Subtract(center)
		a := ray.Direction.LengthSquared()
		halfB := oc.Dot(ray.Direction)
		c := oc.LengthSquared() - radius*radius
		discriminant := halfB*halfB - a*c

		if discriminant < 0 {
			return nil
		}

		sqrtDiscriminant := math.Sqrt(discriminant)

		// Find the nearest root that lies in the acceptable range.
		root := (-halfB - sqrtDiscriminant) / a
		if root < tMin || tMax < root {
			root = (-halfB + sqrtDiscriminant) / a
			if root < tMin || tMax < root {
				return nil
			}
		}

		point := ray.At(root)

		return &Hit{
			Point:    point,
			T:        root,
			Normal:   point.Subtract(center).Multiply(1 / radius),
			Material: material,
		}
	}
}
