package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"image/color"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type ID [16]byte

type Server struct {
	*http.ServeMux

	world     Hittable
	viewAngle float64

	clientsLock sync.RWMutex
	clients     map[ID]io.Writer // map client id -> response writer

	cancelCurrentLock sync.RWMutex
	cancelCurrent     context.CancelFunc
}

func NewServer(world Hittable) *Server {
	s := &Server{
		ServeMux: http.NewServeMux(),

		clients: make(map[ID]io.Writer),
		world:   world,
	}

	s.HandleFunc("/stream", s.streamImage())
	s.HandleFunc("/change", s.changeValue())
	s.HandleFunc("/randomize", s.randomizeColors())
	s.HandleFunc("/", http.FileServer(http.FS(os.DirFS("assets"))).ServeHTTP)

	return s
}

func generateClientID() ID {
	id := [16]byte{}
	rand.Read(id[:])
	return id
}

func (s *Server) streamImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientID := generateClientID()

		s.clientsLock.Lock()
		s.clients[clientID] = w
		s.clientsLock.Unlock()

		defer func() {
			s.clientsLock.Lock()
			delete(s.clients, clientID)
			s.clientsLock.Unlock()
		}()

		s.drawForAllListeners(ctx)

		<-ctx.Done()
	}
}

func (s *Server) drawForAllListeners(ctx context.Context) {
	draw(ctx, s.world, 400, 300, s.viewAngle, func(x, y int, color color.RGBA) {
		// prevent concurrent write while iterating clients
		s.clientsLock.RLock()
		defer s.clientsLock.RUnlock()

		// prevent concurrent write on responses
		s.cancelCurrentLock.Lock()
		defer s.cancelCurrentLock.Unlock()

		for _, w := range s.clients {
			// format: XX YY R G B (little endian, 8 bits per character)
			binary.Write(w, binary.LittleEndian, uint16(x))
			binary.Write(w, binary.LittleEndian, uint16(y))
			binary.Write(w, binary.LittleEndian, uint8(color.R))
			binary.Write(w, binary.LittleEndian, uint8(color.G))
			binary.Write(w, binary.LittleEndian, uint8(color.B))
		}
	})
}

func (s *Server) changeValue() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		s.cancelCurrentLock.Lock()
		if s.cancelCurrent != nil {
			s.cancelCurrent()
		}
		ctx, s.cancelCurrent = context.WithCancel(ctx)
		s.cancelCurrentLock.Unlock()

		angleInDegrees, _ := strconv.Atoi(r.FormValue("angle"))
		angleInRadians := float64(angleInDegrees) / 180 * math.Pi
		s.viewAngle = angleInRadians

		s.drawForAllListeners(ctx)
	}
}

func (s *Server) randomizeColors() http.HandlerFunc {
	type ColorMindRequest struct {
		Model string `json:"model"`
	}

	type ColorMindColor [3]float64

	type ColorMindResponse struct {
		Result []ColorMindColor `json:"result"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		encoded, _ := json.Marshal(ColorMindRequest{
			Model: "default",
		})

		req, _ := http.NewRequest("POST", "http://colormind.io/api/", bytes.NewReader(encoded))
		req = req.WithContext(ctx)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("error fetching random color: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		decoded := ColorMindResponse{}
		err = json.NewDecoder(res.Body).Decode(&decoded)
		if err != nil {
			log.Printf("error decoding random color response %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		colors := make([]Color, len(decoded.Result))
		for i, color := range decoded.Result {
			colors[i] = Color{color[0] / 0xff, color[1] / 0xff, color[2] / 0xff}
		}
		s.world = generateWorld(colors)
		s.drawForAllListeners(ctx)
	}
}
