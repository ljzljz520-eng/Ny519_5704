package preset

import (
	"path/filepath"
	"testing"
)

func TestLoadDefault(t *testing.T) {
	value, err := LoadDefault()
	if err != nil {
		t.Fatalf("LoadDefault() error = %v", err)
	}
	if value.Name != "unit-circle-period" {
		t.Fatalf("Name = %q, want %q", value.Name, "unit-circle-period")
	}
	if value.Curve.Points != 101 {
		t.Fatalf("Curve.Points = %d, want 101", value.Curve.Points)
	}
}

func TestLoadFile(t *testing.T) {
	value, err := LoadFile(filepath.Join("testdata", "classroom.yaml"))
	if err != nil {
		t.Fatalf("LoadFile() error = %v", err)
	}
	if value.Curve.Amplitude != 2.5 || value.View.Yaw != 45 {
		t.Fatalf("LoadFile() = %+v", value)
	}
}
