package main

import (
	"context"
	"encoding/binary"
	"image/color"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/google/uuid"
)

type Server struct {
	*http.ServeMux

	clientsLock sync.RWMutex
	clients     map[uuid.UUID]io.Writer // map client id -> response writer

	cancelCurrentLock sync.RWMutex
	cancelCurrent     context.CancelFunc
}

func NewServer() *Server {
	s := &Server{
		ServeMux: http.NewServeMux(),

		clients: make(map[uuid.UUID]io.Writer),
	}

	s.HandleFunc("/stream", s.streamImage())
	s.HandleFunc("/change", s.changeValue())
	s.HandleFunc("/", http.FileServer(http.FS(os.DirFS("assets"))).ServeHTTP)

	return s
}

func (s *Server) streamImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := uuid.New()

		s.clientsLock.Lock()
		s.clients[id] = w
		s.clientsLock.Unlock()

		defer func() {
			s.clientsLock.Lock()
			delete(s.clients, id)
			s.clientsLock.Unlock()
		}()

		s.drawForAllListeners(ctx, 0)

		<-ctx.Done()
	}
}

func (s *Server) drawForAllListeners(ctx context.Context, angle float64) {
	draw(ctx, 400, 300, angle, func(x, y int, color color.RGBA) {
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

		s.drawForAllListeners(ctx, angleInRadians)
	}
}
