package main

import (
	"encoding/binary"
	"image/color"
	"net/http"
	"os"
)

type Server struct {
	*http.ServeMux
}

func NewServer() *Server {
	s := &Server{
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/stream", s.streamImage())
	s.HandleFunc("/", http.FileServer(http.FS(os.DirFS("assets"))).ServeHTTP)

	return s
}

func (s *Server) streamImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		draw(400, 300, func(x, y int, color color.RGBA) {
			// format: XX YY R G B (little endian, 8 bits per character)
			binary.Write(w, binary.LittleEndian, uint16(x))
			binary.Write(w, binary.LittleEndian, uint16(y))
			binary.Write(w, binary.LittleEndian, uint8(color.R))
			binary.Write(w, binary.LittleEndian, uint8(color.G))
			binary.Write(w, binary.LittleEndian, uint8(color.B))
		})
	}
}
