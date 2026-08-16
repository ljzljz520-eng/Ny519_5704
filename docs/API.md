# Sine3D HTTP API

The service represents a sine curve as `(x, amplitude * sin(frequency * x + phase), 0)`. Angles in the curve formula are radians. View rotation angles are degrees.

## Start

~~~sh
GOTOOLCHAIN=local CGO_ENABLED=0 go run ./cmd/sine3d -addr :8080
~~~

Pass `-fixture path/to/preset.yaml` to replace the embedded deterministic YAML preset.

## Generate coordinates

`POST /api/v1/coordinates`

~~~json
{
  "amplitude": 1,
  "frequency": 1,
  "phase": 0,
  "start": 0,
  "end": 6.283185307179586,
  "points": 101
}
~~~

The successful response contains the normalized function input, the sampling step, and an ordered `points` array.

~~~json
{
  "function": {
    "amplitude": 1,
    "frequency": 1,
    "phase": 0,
    "start": 0,
    "end": 6.283185307179586,
    "points": 101
  },
  "step": 0.06283185307179587,
  "points": [
    {"x": 0, "y": 0, "z": 0},
    {"x": 0.06283185307179587, "y": 0.06279051952931337, "z": 0}
  ]
}
~~~

## Download coordinates

`POST /api/v1/coordinates/download` accepts the same body and returns the same JSON schema with `Content-Disposition: attachment; filename="sine3d-coordinates.json"`.

## Rotate coordinates

`POST /api/v1/coordinates/rotate`

~~~json
{
  "curve": {
    "amplitude": 1,
    "frequency": 1,
    "phase": 0,
    "start": 0,
    "end": 6.283185307179586,
    "points": 101
  },
  "view": {"yaw": 30, "pitch": 25, "roll": 0}
}
~~~

The API applies yaw around Y, pitch around X, then roll around Z. The response has `function`, `view`, and the rotated `points` array.

## Preset and health

`GET /api/v1/presets/default` returns the YAML-backed curve and view. `GET /healthz` returns `{"status":"ok"}`. Browser clients may call all endpoints cross-origin.

## Browser consumption

~~~js
const response = await fetch("http://localhost:8080/api/v1/coordinates/rotate", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify(request)
});
if (!response.ok) throw new Error(await response.text());
const data = await response.json();
const chartPoints = data.points.map(({ x, y, z }) => [x, y, z]);
~~~

The standalone request and representative response files in `docs/examples` can be used as frontend fixtures.

## Validation

Amplitude must be non-negative, frequency must be positive, `end` must be greater than `start`, every number must be finite, and `points` must be between 2 and 100000. Invalid JSON returns 400; invalid parameters return 422.
