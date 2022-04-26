<img width="604" alt="image" src="https://user-images.githubusercontent.com/2004131/165345254-7ea83893-440c-475d-b884-53d7169cefaa.png">

This is a pure Go implementation (CPU-based) of [Raytracing in One Weekend](https://raytracing.github.io).

### Usage of Go standard library

- flag is used for parsing CLI parameters
- crypto/rand for random client ids
- image/color for working with RGBA colors
- encoding/binary for little endian encoding of image bytes
- sync for R/W mutexes
- net/http for http server
- net/http/httptest for integration tests
- context for cancelling expensive renders
