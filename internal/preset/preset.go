package preset

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"

	"example.com/sine3d/internal/curve"
	"gopkg.in/yaml.v3"
)

type Preset struct {
	Name  string              `json:"name" yaml:"name"`
	Curve curve.SampleRequest `json:"curve" yaml:"curve"`
	View  curve.Rotation      `json:"view" yaml:"view"`
}

//go:embed default_curve.yaml
var defaultFixture []byte

func Load(reader io.Reader) (Preset, error) {
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)
	var value Preset
	if err := decoder.Decode(&value); err != nil {
		return Preset{}, fmt.Errorf("decode preset: %w", err)
	}
	if value.Name == "" {
		return Preset{}, fmt.Errorf("preset name is required")
	}
	if err := value.Curve.Validate(); err != nil {
		return Preset{}, fmt.Errorf("validate preset curve: %w", err)
	}
	return value, nil
}

func LoadDefault() (Preset, error) {
	return Load(bytes.NewReader(defaultFixture))
}

func LoadFile(path string) (Preset, error) {
	file, err := os.Open(path)
	if err != nil {
		return Preset{}, fmt.Errorf("open preset: %w", err)
	}
	defer file.Close()
	return Load(file)
}
