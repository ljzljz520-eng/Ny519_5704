package curve

import (
	"math"
	"testing"
)

func TestGenerateReturnsOrderedSamples(t *testing.T) {
	request := SampleRequest{
		Amplitude: 2,
		Frequency: 1,
		Phase:     0,
		Start:     0,
		End:       math.Pi,
		Points:    5,
	}
	result, err := Generate(request)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if len(result.Points) != request.Points {
		t.Fatalf("len(Points) = %d, want %d", len(result.Points), request.Points)
	}
	if result.Step != math.Pi/4 {
		t.Fatalf("Step = %v, want %v", result.Step, math.Pi/4)
	}
	for index := 1; index < len(result.Points)-1; index++ {
		if result.Points[index].X <= result.Points[index-1].X {
			t.Fatalf("Points[%d].X = %v, previous = %v", index, result.Points[index].X, result.Points[index-1].X)
		}
	}
	if math.Abs(result.Points[2].Y-2) > 1e-12 {
		t.Fatalf("Points[2].Y = %v, want 2", result.Points[2].Y)
	}
}

func TestSampleRequestValidation(t *testing.T) {
	valid := SampleRequest{Amplitude: 1, Frequency: 1, Start: 0, End: 1, Points: 2}
	tests := []struct {
		name    string
		request SampleRequest
	}{
		{name: "negative amplitude", request: SampleRequest{Amplitude: -1, Frequency: 1, Start: 0, End: 1, Points: 2}},
		{name: "zero frequency", request: SampleRequest{Amplitude: 1, Frequency: 0, Start: 0, End: 1, Points: 2}},
		{name: "reversed interval", request: SampleRequest{Amplitude: 1, Frequency: 1, Start: 1, End: 0, Points: 2}},
		{name: "too few points", request: SampleRequest{Amplitude: 1, Frequency: 1, Start: 0, End: 1, Points: 1}},
		{name: "non-finite phase", request: SampleRequest{Amplitude: 1, Frequency: 1, Phase: math.Inf(1), Start: 0, End: 1, Points: 2}},
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid request error = %v", err)
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.request.Validate(); err == nil {
				t.Fatal("Validate() error = nil")
			}
		})
	}
}

func TestRotateAroundZAxis(t *testing.T) {
	result := Rotate([]Point{{X: 1, Y: 0, Z: 0}}, Rotation{Roll: 90})
	if math.Abs(result[0].X) > 1e-12 || math.Abs(result[0].Y-1) > 1e-12 || math.Abs(result[0].Z) > 1e-12 {
		t.Fatalf("Rotate() = %+v, want {X:0 Y:1 Z:0}", result[0])
	}
}
