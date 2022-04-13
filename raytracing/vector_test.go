package raytracing_test

import (
	"math"
	"testing"

	"github.com/davherrmann/rtgo/raytracing"
)

type Vec = raytracing.Vec

type Comparable interface {
	float64 | Vec
}

func requireEqual[T Comparable](t *testing.T, a, b T) {
	t.Helper()

	maxDelta := 1e-8

	switch a := any(a).(type) {
	case Vec:
		deltaX := math.Abs(a.X - any(b).(Vec).X)
		deltaY := math.Abs(a.Y - any(b).(Vec).Y)
		deltaZ := math.Abs(a.Z - any(b).(Vec).Z)

		if deltaX > maxDelta || deltaY > maxDelta || deltaZ > maxDelta {
			t.Errorf("vector %#v and %#v expected to be equal, but are not", a, b)
		}
	case float64:
		delta := math.Abs(a - any(b).(float64))
		if delta > maxDelta {
			t.Errorf("float64 %#v and %#v expected to be equal, but are not", a, b)
		}
	}
}

func TestVectorAdd(t *testing.T) {
	tests := []struct {
		A        Vec
		B        Vec
		Expected Vec
	}{
		{A: Vec{X: 1}, B: Vec{X: 1}, Expected: Vec{X: 2}},
		{A: Vec{X: -1}, B: Vec{X: 1}, Expected: Vec{}},
		{A: Vec{X: 1, Y: 1, Z: 1}, B: Vec{X: 2, Y: 3, Z: 4}, Expected: Vec{X: 3, Y: 4, Z: 5}},
	}

	for _, test := range tests {
		actual := test.A.Add(test.B)
		requireEqual(t, actual, test.Expected)
	}
}

func TestVectorMultiply(t *testing.T) {
	tests := []struct {
		A        Vec
		B        float64
		Expected Vec
	}{
		{A: Vec{X: 1}, B: 1, Expected: Vec{X: 1}},
		{A: Vec{X: 1}, B: 2, Expected: Vec{X: 2}},
		{A: Vec{X: 2}, B: -2, Expected: Vec{X: -4}},
		{A: Vec{X: 3, Y: 4, Z: 5}, B: -0.5, Expected: Vec{X: -1.5, Y: -2, Z: -2.5}},
	}

	for _, test := range tests {
		actual := test.A.Multiply(test.B)
		requireEqual(t, actual, test.Expected)
	}
}

func TestVectorDot(t *testing.T) {
	tests := []struct {
		A        Vec
		B        Vec
		Expected float64
	}{
		{A: Vec{X: 1, Y: 2, Z: 3}, B: Vec{X: 4, Y: 5, Z: 6}, Expected: 32},
		{A: Vec{X: -1, Y: 2, Z: -3}, B: Vec{X: 4, Y: 5, Z: 6}, Expected: -12},
	}

	for _, test := range tests {
		actual := test.A.Dot(test.B)
		requireEqual(t, actual, test.Expected)
	}
}
