package main

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"image/color"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type ID [16]byte

type Server struct {
	*http.ServeMux

	clientsLock sync.RWMutex
	clients     map[ID]io.Writer // map client id -> response writer

	cancelCurrentLock sync.RWMutex
	cancelCurrent     context.CancelFunc
}

func NewServer() *Server {
	s := &Server{
		ServeMux: http.NewServeMux(),

		clients: make(map[ID]io.Writer),
	}

	s.HandleFunc("/stream", s.streamImage())
	s.HandleFunc("/change", s.changeValue())
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
