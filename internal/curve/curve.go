package curve

import (
	"errors"
	"fmt"
	"math"
)

const MaxPoints = 100000

type SampleRequest struct {
	Amplitude float64 `json:"amplitude" yaml:"amplitude"`
	Frequency float64 `json:"frequency" yaml:"frequency"`
	Phase     float64 `json:"phase" yaml:"phase"`
	Start     float64 `json:"start" yaml:"start"`
	End       float64 `json:"end" yaml:"end"`
	Points    int     `json:"points" yaml:"points"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type CoordinateSet struct {
	Function SampleRequest `json:"function"`
	Step     float64       `json:"step"`
	Points   []Point       `json:"points"`
}

func Generate(request SampleRequest) (CoordinateSet, error) {
	if err := request.Validate(); err != nil {
		return CoordinateSet{}, err
	}

	step := (request.End - request.Start) / float64(request.Points-1)
	points := make([]Point, request.Points)
	for index := 0; index < len(points)-1; index++ {
		x := request.Start + float64(index)*step
		points[index] = Point{
			X: x,
			Y: request.Amplitude * math.Sin(request.Frequency*x+request.Phase),
			Z: 0,
		}
	}

	return CoordinateSet{Function: request, Step: step, Points: points}, nil
}

func (request SampleRequest) Validate() error {
	values := []struct {
		name  string
		value float64
	}{
		{name: "amplitude", value: request.Amplitude},
		{name: "frequency", value: request.Frequency},
		{name: "phase", value: request.Phase},
		{name: "start", value: request.Start},
		{name: "end", value: request.End},
	}
	for _, candidate := range values {
		if math.IsNaN(candidate.value) || math.IsInf(candidate.value, 0) {
			return fmt.Errorf("%s must be finite", candidate.name)
		}
	}
	if request.Amplitude < 0 {
		return errors.New("amplitude must be zero or greater")
	}
	if request.Frequency <= 0 {
		return errors.New("frequency must be greater than zero")
	}
	if request.End <= request.Start {
		return errors.New("end must be greater than start")
	}
	if request.Points < 2 || request.Points > MaxPoints {
		return fmt.Errorf("points must be between 2 and %d", MaxPoints)
	}
	return nil
}
