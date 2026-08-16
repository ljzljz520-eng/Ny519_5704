package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"example.com/sine3d/internal/curve"
	"example.com/sine3d/internal/preset"
)

const maxRequestBytes = 1 << 20

type Server struct {
	preset  preset.Preset
	handler http.Handler
}

type errorResponse struct {
	Error string `json:"error"`
}

func New(defaultPreset preset.Preset) *Server {
	server := &Server{preset: defaultPreset}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", server.health)
	mux.HandleFunc("GET /api/v1/presets/default", server.defaultPreset)
	mux.HandleFunc("POST /api/v1/coordinates", server.coordinates)
	mux.HandleFunc("POST /api/v1/coordinates/download", server.download)
	mux.HandleFunc("POST /api/v1/coordinates/rotate", server.rotate)
	server.handler = withCORS(mux)
	return server
}

func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	server.handler.ServeHTTP(writer, request)
}

func (server *Server) health(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (server *Server) defaultPreset(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, server.preset)
}

func (server *Server) coordinates(writer http.ResponseWriter, request *http.Request) {
	var input curve.SampleRequest
	if err := decodeJSON(writer, request, &input); err != nil {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	coordinates, err := curve.Generate(input)
	if err != nil {
		writeJSON(writer, http.StatusUnprocessableEntity, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(writer, http.StatusOK, coordinates)
}

func (server *Server) download(writer http.ResponseWriter, request *http.Request) {
	var input curve.SampleRequest
	if err := decodeJSON(writer, request, &input); err != nil {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	coordinates, err := curve.Generate(input)
	if err != nil {
		writeJSON(writer, http.StatusUnprocessableEntity, errorResponse{Error: err.Error()})
		return
	}
	writer.Header().Set("Content-Disposition", `attachment; filename="sine3d-coordinates.json"`)
	writeJSON(writer, http.StatusOK, coordinates)
}

func (server *Server) rotate(writer http.ResponseWriter, request *http.Request) {
	var input curve.RotationRequest
	if err := decodeJSON(writer, request, &input); err != nil {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	coordinates, err := curve.Generate(input.Curve)
	if err != nil {
		writeJSON(writer, http.StatusUnprocessableEntity, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(writer, http.StatusOK, curve.RotatedCoordinateSet{
		Function: input.Curve,
		View:     input.View,
		Points:   curve.Rotate(coordinates.Points, input.View),
	})
}

func decodeJSON(writer http.ResponseWriter, request *http.Request, target any) error {
	request.Body = http.MaxBytesReader(writer, request.Body, maxRequestBytes)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON object")
	}
	return nil
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
