package main

type CameraRay func(u, v float64) Ray

func Camera(up, from, lookAt Vec, width, height int) CameraRay {
	aspectRatio := float64(width) / float64(height)

	// camera & viewport
	viewportHeight := 2.0
	viewportWidth := aspectRatio * viewportHeight

	w := from.Subtract(lookAt).Normalized()
	u := up.Cross(w).Normalized()
	v := w.Cross(u)

	// ray origin
	origin := from
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
