package rtgo

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"image/color"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/davherrmann/rtgo/external/colormind"
	"github.com/davherrmann/rtgo/raytracing"
)

type ID [16]byte

type Server struct {
	*http.ServeMux

	camera      raytracing.Camera
	world       raytracing.Hittable
	clientsLock sync.RWMutex
	clients     map[ID]io.Writer // map client id -> response writer

	cancelCurrentLock sync.RWMutex
	cancelCurrent     context.CancelFunc

	done func() <-chan struct{}
}

func NewServer(camera raytracing.Camera, world raytracing.Hittable) *Server {
	s := &Server{
		ServeMux: http.NewServeMux(),

		clients: make(map[ID]io.Writer),
		world:   world,
		camera:  camera,
	}

	s.HandleFunc("/stream", s.streamImage())
	s.HandleFunc("/change", s.changeValue())
	s.HandleFunc("/randomize", s.randomizeColors())
	s.HandleFunc("/", http.FileServer(http.FS(os.DirFS("assets"))).ServeHTTP)

	return s
}

func (s *Server) ListenAndServe(port string) {
	httpServer := http.Server{
		Addr:    ":" + port,
		Handler: s,
	}

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	s.done = func() <-chan struct{} {
		return ctx.Done()
	}

	log.Println("listening on port " + port)
	go httpServer.ListenAndServe()

	<-ctx.Done()
	stop()

	log.Println("shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Println(err)
	}
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

		select {
		case <-ctx.Done():
		case <-s.done():
		}
	}
}

func (s *Server) drawForAllListeners(ctx context.Context) {
	options := raytracing.RenderOptions{
		ResolutionX: 400,
		ResolutionY: 300,
	}
	raytracing.Render(ctx, s.world, s.camera, options, func(x, y int, color color.RGBA) {
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

		zoomPercent, _ := strconv.Atoi(r.FormValue("zoom"))
		zoom := float64(zoomPercent)/200 + 0.2

		s.camera = GenerateCamera(angleInRadians, zoom, 400, 300)

		s.drawForAllListeners(ctx)
	}
}

func (s *Server) randomizeColors() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		randomColorPalette, err := colormind.FetchRandomPalette(ctx)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			log.Printf("error fetching random color palette: %v", err)
			return
		}

		s.world = GenerateWorld(randomColorPalette)
		s.drawForAllListeners(ctx)
	}
}
