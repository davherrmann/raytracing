package main

var Black = Color{}

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
