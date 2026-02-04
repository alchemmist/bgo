package app

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"bgo/internal/weather"
)

type fakeClient struct {
	coords   weather.Coordinates
	now      map[string]any
	forecast map[string]any
}

func (f fakeClient) GetCoordinates(ctx context.Context) (weather.Coordinates, error) {
	return f.coords, nil
}

func (f fakeClient) GetWeatherNow(ctx context.Context, coords weather.Coordinates) (map[string]any, error) {
	return f.now, nil
}

func (f fakeClient) GetWeatherForecast(ctx context.Context, coords weather.Coordinates) (map[string]any, error) {
	return f.forecast, nil
}

func TestRunnerNow(t *testing.T) {
	out := &bytes.Buffer{}
	in := strings.NewReader("n\n")

	runner := &Runner{
		Client: fakeClient{
			coords: weather.Coordinates{Lat: "1", Lon: "2"},
			now: map[string]any{
				"main": map[string]any{
					"temp":       10.0,
					"feels_like": 8.0,
					"humidity":   60.0,
				},
				"weather": []any{map[string]any{"id": 800.0, "description": "clear"}},
				"name":    "Testville",
			},
		},
		Out:    out,
		In:     in,
		ErrOut: out,
	}

	code := runner.Run([]string{"now"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(out.String(), "Testville") {
		t.Fatalf("expected output to include location")
	}
}

func TestRunnerInvalidDays(t *testing.T) {
	out := &bytes.Buffer{}
	runner := &Runner{
		Client: fakeClient{},
		Out:    out,
		In:     strings.NewReader("n\n"),
		ErrOut: out,
	}

	code := runner.Run([]string{"forecast", "-d", "8"})
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
}

func TestRunnerFullInfo(t *testing.T) {
	out := &bytes.Buffer{}
	runner := &Runner{
		Client: fakeClient{
			coords: weather.Coordinates{Lat: "1", Lon: "2"},
			now: map[string]any{
				"main": map[string]any{
					"temp":       10.0,
					"feels_like": 8.0,
					"humidity":   60.0,
				},
				"weather": []any{map[string]any{"id": 800.0, "description": "clear"}},
				"name":    "Testville",
			},
		},
		Out:    out,
		In:     strings.NewReader("n\n"),
		ErrOut: out,
	}

	code := runner.Run([]string{"now", "--full-info"})
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(out.String(), "Testville") {
		t.Fatalf("expected json output to include location")
	}
}
