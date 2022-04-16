package raytracing

import "math"

type Hit struct {
	Point     Vec
	Normal    Vec
	T         float64
	FrontFace bool
	Material  Material
}

type Hittable interface {
	Hit(ray Ray, tMin, tMax float64) *Hit
}

type World struct {
	Objects []Hittable
}

func (w World) Hit(ray Ray, tMin, tMax float64) *Hit {
	var closestHit *Hit
	closest := tMax

	for _, hittable := range w.Objects {
		hit := hittable.Hit(ray, tMin, tMax)
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

type Sphere struct {
	Center   Vec
	Radius   float64
	Material Material
}

func (s Sphere) Hit(ray Ray, tMin, tMax float64) *Hit {
	oc := ray.Origin.Subtract(s.Center)
	a := ray.Direction.LengthSquared()
	halfB := oc.Dot(ray.Direction)
	c := oc.LengthSquared() - s.Radius*s.Radius
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

	// normal and front face
	normal := point.Subtract(s.Center).Multiply(1 / s.Radius)
	frontFace := ray.Direction.Dot(normal) < 0
	if !frontFace {
		normal = normal.Multiply(-1)
	}

	return &Hit{
		Point:     point,
		T:         root,
		Normal:    normal,
		Material:  s.Material,
		FrontFace: frontFace,
	}
}
