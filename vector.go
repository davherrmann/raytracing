package main

import "math"

type Vector struct {
	X float64
	Y float64
	Z float64
}

func (v Vector) Multiply(m float64) Vector {
	return Vector{
		X: v.X * m,
		Y: v.Y * m,
		Z: v.Z * m,
	}
}

func (a Vector) Add(b Vector) Vector {
	return Vector{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

func (a Vector) Subtract(b Vector) Vector {
	return Vector{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func (v Vector) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v Vector) Normalized() Vector {
	length := v.Length()
	return v.Multiply(1 / length)
}

func (v Vector) Reflect(normal Vector) Vector {
	return v.Subtract(normal.Multiply(2 * v.Dot(normal)))
}

func (v Vector) Refract(normal Vector, refractionRatio float64) Vector {
	cosTheta := math.Min(v.Multiply(-1).Dot(normal), 1)
	rOutPerp := v.Add(normal.Multiply(cosTheta)).Multiply(refractionRatio)
	rOutParallel := normal.Multiply(-math.Sqrt(math.Abs(1 - rOutPerp.LengthSquared())))

	return rOutPerp.Add(rOutParallel)
}

func (a Vector) Dot(b Vector) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vector) Cross(b Vector) Vector {
	return Vector{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}
