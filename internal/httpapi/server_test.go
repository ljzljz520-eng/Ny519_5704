package httpapi

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/sine3d/internal/curve"
	"example.com/sine3d/internal/preset"
)

func TestCoordinateWorkflowIncludesSamplingIntervalEnd(t *testing.T) {
	handler := testHandler(t)
	body := `{"amplitude":1,"frequency":1,"phase":0,"start":0,"end":6.283185307179586,"points":101}`
	response := serve(handler, http.MethodPost, "/api/v1/coordinates", body)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	var result curve.CoordinateSet
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(result.Points) != 101 {
		t.Fatalf("len(Points) = %d, want 101", len(result.Points))
	}
	last := result.Points[len(result.Points)-1]
	if math.Abs(last.X-2*math.Pi) > 1e-12 {
		t.Fatalf("last X = %.16g, want %.16g", last.X, 2*math.Pi)
	}
}

func TestDownloadWorkflow(t *testing.T) {
	handler := testHandler(t)
	body := `{"amplitude":2,"frequency":1,"phase":0,"start":0,"end":1.5707963267948966,"points":3}`
	response := serve(handler, http.MethodPost, "/api/v1/coordinates/download", body)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if got := response.Header().Get("Content-Disposition"); got != `attachment; filename="sine3d-coordinates.json"` {
		t.Fatalf("Content-Disposition = %q", got)
	}
	var result curve.CoordinateSet
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(result.Points) != 3 {
		t.Fatalf("len(Points) = %d, want 3", len(result.Points))
	}
}

func TestRotationWorkflow(t *testing.T) {
	handler := testHandler(t)
	body := `{"curve":{"amplitude":1,"frequency":1,"phase":0,"start":1,"end":2,"points":3},"view":{"yaw":0,"pitch":0,"roll":90}}`
	response := serve(handler, http.MethodPost, "/api/v1/coordinates/rotate", body)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	var result curve.RotatedCoordinateSet
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(result.Points) != 3 {
		t.Fatalf("len(Points) = %d, want 3", len(result.Points))
	}
	if math.Abs(result.Points[0].X+math.Sin(1)) > 1e-12 || math.Abs(result.Points[0].Y-1) > 1e-12 {
		t.Fatalf("first rotated point = %+v", result.Points[0])
	}
}

func TestPresetWorkflow(t *testing.T) {
	handler := testHandler(t)
	response := serve(handler, http.MethodGet, "/api/v1/presets/default", "")
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	var result preset.Preset
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result.Name != "unit-circle-period" || result.Curve.Points != 101 {
		t.Fatalf("preset = %+v", result)
	}
}

func TestRequestErrors(t *testing.T) {
	handler := testHandler(t)
	tests := []struct {
		name   string
		body   string
		status int
	}{
		{name: "unknown field", body: `{"amplitude":1,"frequency":1,"start":0,"end":1,"points":2,"extra":true}`, status: http.StatusBadRequest},
		{name: "invalid interval", body: `{"amplitude":1,"frequency":1,"start":2,"end":1,"points":2}`, status: http.StatusUnprocessableEntity},
		{name: "multiple objects", body: `{} {}`, status: http.StatusBadRequest},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := serve(handler, http.MethodPost, "/api/v1/coordinates", test.body)
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d, body = %s", response.Code, test.status, response.Body.String())
			}
		})
	}
}

func TestConcurrentCoordinateRequests(t *testing.T) {
	handler := testHandler(t)
	const workers = 8
	ready := make(chan struct{}, workers)
	start := make(chan struct{})
	results := make(chan *httptest.ResponseRecorder, workers)
	for index := 0; index < workers; index++ {
		go func() {
			ready <- struct{}{}
			<-start
			results <- serve(handler, http.MethodPost, "/api/v1/coordinates", `{"amplitude":1,"frequency":1,"phase":0,"start":0,"end":1,"points":5}`)
		}()
	}
	for index := 0; index < workers; index++ {
		<-ready
	}
	close(start)
	for index := 0; index < workers; index++ {
		response := <-results
		if response.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
		}
	}
}

func testHandler(t *testing.T) http.Handler {
	t.Helper()
	value, err := preset.LoadDefault()
	if err != nil {
		t.Fatalf("LoadDefault() error = %v", err)
	}
	return New(value)
}

func serve(handler http.Handler, method string, path string, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if strings.TrimSpace(body) != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}
