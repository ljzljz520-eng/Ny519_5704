package curve

import "math"

type Rotation struct {
	Yaw   float64 `json:"yaw" yaml:"yaw"`
	Pitch float64 `json:"pitch" yaml:"pitch"`
	Roll  float64 `json:"roll" yaml:"roll"`
}

type RotationRequest struct {
	Curve SampleRequest `json:"curve"`
	View  Rotation      `json:"view"`
}

type RotatedCoordinateSet struct {
	Function SampleRequest `json:"function"`
	View     Rotation      `json:"view"`
	Points   []Point       `json:"points"`
}

func Rotate(points []Point, rotation Rotation) []Point {
	yaw := rotation.Yaw * math.Pi / 180
	pitch := rotation.Pitch * math.Pi / 180
	roll := rotation.Roll * math.Pi / 180
	cy, sy := math.Cos(yaw), math.Sin(yaw)
	cp, sp := math.Cos(pitch), math.Sin(pitch)
	cr, sr := math.Cos(roll), math.Sin(roll)

	rotated := make([]Point, len(points))
	for index, point := range points {
		x1 := cy*point.X + sy*point.Z
		y1 := point.Y
		z1 := -sy*point.X + cy*point.Z
		x2 := x1
		y2 := cp*y1 - sp*z1
		z2 := sp*y1 + cp*z1
		rotated[index] = Point{
			X: cr*x2 - sr*y2,
			Y: sr*x2 + cr*y2,
			Z: z2,
		}
	}
	return rotated
}
