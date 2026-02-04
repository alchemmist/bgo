package weather

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestClientGetCoordinates(t *testing.T) {
	client := NewClient("key")
	client.HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != "http://ipinfo.io/json" {
			return &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(bytes.NewBufferString(""))}, nil
		}
		payload, _ := json.Marshal(map[string]any{"loc": "10.1,20.2"})
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBuffer(payload))}, nil
	})}

	coords, err := client.GetCoordinates(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if coords.Lat != "10.1" || coords.Lon != "20.2" {
		t.Fatalf("unexpected coordinates: %+v", coords)
	}
}

func TestClientGetWeatherNow(t *testing.T) {
	client := NewClient("key")
	client.OpenWeatherBase = "https://example.test"
	client.WeatherNowPath = "/now"
	client.HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/now" {
			return &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(bytes.NewBufferString(""))}, nil
		}
		query := req.URL.Query()
		if query.Get("appid") != "key" || query.Get("lat") != "55" || query.Get("lon") != "37" {
			return &http.Response{StatusCode: http.StatusBadRequest, Body: io.NopCloser(bytes.NewBufferString(""))}, nil
		}
		payload, _ := json.Marshal(map[string]any{"ok": true})
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBuffer(payload))}, nil
	})}

	response, err := client.GetWeatherNow(context.Background(), Coordinates{Lat: "55", Lon: "37"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response["ok"] != true {
		t.Fatalf("expected ok response")
	}
}
