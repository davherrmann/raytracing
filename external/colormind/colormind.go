package colormind

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/davherrmann/rtgo/raytracing"
)

type ColorMindRequest struct {
	Model string `json:"model"`
}

type ColorMindColor [3]float64

type ColorMindResponse struct {
	Result []ColorMindColor `json:"result"`
}

func FetchRandomPalette(ctx context.Context) ([]raytracing.Color, error) {
	encoded, _ := json.Marshal(ColorMindRequest{
		Model: "default",
	})

	req, _ := http.NewRequest("POST", "http://colormind.io/api/", bytes.NewReader(encoded))
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	decoded := ColorMindResponse{}
	err = json.NewDecoder(res.Body).Decode(&decoded)
	if err != nil {
		return nil, err
	}

	colors := make([]raytracing.Color, len(decoded.Result))
	for i, color := range decoded.Result {
		colors[i] = raytracing.Color{
			R: color[0] / 0xff,
			G: color[1] / 0xff,
			B: color[2] / 0xff,
		}
	}

	return colors, nil
}
