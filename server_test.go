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
	camera := raytracing.Camera{
		Up:     Vec{X: 1},
		From:   Vec{Z: -1},
		LookAt: Vec{},
		Zoom:   1,
	}
	world := raytracing.World()
	server := rtgo.NewServer(camera, world)

	test := httptest.NewServer(server)

	res, _ := http.Get(test.URL + "/")

	if res.StatusCode != http.StatusOK {
		t.Fail()
	}
}
