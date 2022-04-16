package raytracing

type Camera struct {
	Up     Vec
	From   Vec
	LookAt Vec

	Zoom float64
}

func (c *Camera) RayCaster(aspectRatio float64) RayCaster {
	// camera & viewport
	viewportHeight := c.Zoom
	viewportWidth := aspectRatio * viewportHeight

	w := c.From.Subtract(c.LookAt).Normalized()
	u := c.Up.Cross(w).Normalized()
	v := w.Cross(u)

	// ray origin
	origin := c.From
	horizontal := u.Multiply(viewportWidth)
	vertical := v.Multiply(viewportHeight)
	lowerLeftCorner := origin. //
					Subtract(horizontal.Multiply(0.5)). //
					Subtract(vertical.Multiply(0.5)).   //
					Subtract(w)

	return func(u, v float64) Ray {
		rayDirection := lowerLeftCorner. //
							Add(horizontal.Multiply(u)). //
							Add(vertical.Multiply(v)).   //
							Subtract(origin)

		return Ray{origin, rayDirection}
	}
}
