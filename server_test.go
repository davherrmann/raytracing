package rtgo_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davherrmann/rtgo"
	"github.com/davherrmann/rtgo/raytracing"
)

type Vec = raytracing.Vec

func TestServer(t *testing.T) {
	camera := raytracing.Camera(Vec{X: 1, Y: 0, Z: 0}, Vec{Z: -1}, Vec{}, 1, 1, 1)
	world := raytracing.World()
	server := rtgo.NewServer(camera, world)

	test := httptest.NewServer(server)

	res, _ := http.Get(test.URL + "/")

	if res.StatusCode != http.StatusOK {
		t.Fail()
	}
}
