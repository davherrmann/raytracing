<img width="604" alt="image" src="https://user-images.githubusercontent.com/2004131/165345254-7ea83893-440c-475d-b884-53d7169cefaa.png">

### About

This is a pure Go implementation (CPU-based) of [Raytracing in One Weekend](https://raytracing.github.io) with some additions:

- live streaming the rendering to the browser and rendering to a canvas element
- change camera angle and distance from the browser, immediately cancelling old renders
- fetching random color palettes from colormind.io

### Run locally

1. Clone the repository.
2. `cd` into the cloned directory.
3. Run `go run cmd/rtgo/main.go`.
4. Open http://localhost:8080 in a browser.

### Usage of Go standard library

- flag is used for parsing CLI parameters
- crypto/rand for random client ids
- image/color for working with RGBA colors
- encoding/binary for little endian encoding of image bytes
- sync for R/W mutexes
- net/http for http server
- net/http/httptest for integration tests
- context for cancelling expensive renders
