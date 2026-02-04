package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Coordinates struct {
	Lat string
	Lon string
}

type Client struct {
	HTTPClient          *http.Client
	IPInfoURL           string
	OpenWeatherBase     string
	APIKey              string
	RequestTimeout      time.Duration
	CoordinatesPath     string
	WeatherNowPath      string
	WeatherForecastPath string
}

func NewClient(apiKey string) *Client {
	return &Client{
		HTTPClient:          &http.Client{Timeout: 15 * time.Second},
		IPInfoURL:           "http://ipinfo.io/json",
		OpenWeatherBase:     "https://api.openweathermap.org",
		APIKey:              apiKey,
		RequestTimeout:      15 * time.Second,
		CoordinatesPath:     "",
		WeatherNowPath:      "/data/2.5/weather",
		WeatherForecastPath: "/data/2.5/forecast",
	}
}

func (c *Client) GetCoordinates(ctx context.Context) (Coordinates, error) {
	client := c.client()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.IPInfoURL, nil)
	if err != nil {
		return Coordinates{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return Coordinates{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return Coordinates{}, fmt.Errorf("ipinfo status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return Coordinates{}, err
	}

	loc, ok := data["loc"].(string)
	if !ok {
		return Coordinates{}, errors.New("missing loc in ipinfo response")
	}

	parts := strings.Split(loc, ",")
	if len(parts) != 2 {
		return Coordinates{}, errors.New("invalid loc format")
	}

	return Coordinates{Lat: strings.TrimSpace(parts[0]), Lon: strings.TrimSpace(parts[1])}, nil
}

func (c *Client) GetWeatherNow(ctx context.Context, coords Coordinates) (map[string]any, error) {
	return c.getWeather(ctx, c.WeatherNowPath, coords)
}

func (c *Client) GetWeatherForecast(ctx context.Context, coords Coordinates) (map[string]any, error) {
	return c.getWeather(ctx, c.WeatherForecastPath, coords)
}

func (c *Client) getWeather(ctx context.Context, path string, coords Coordinates) (map[string]any, error) {
	if c.APIKey == "" {
		return nil, errors.New("OPEN_WEATHER_API_KEY is not set")
	}

	params := BasicRequestParams(c.APIKey)
	params["lat"] = coords.Lat
	params["lon"] = coords.Lon

	urlValues := url.Values{}
	for k, v := range params {
		urlValues.Set(k, v)
	}

	endpoint := strings.TrimSuffix(c.OpenWeatherBase, "/") + path
	if !strings.Contains(endpoint, "?") {
		endpoint += "?" + urlValues.Encode()
	} else {
		endpoint += "&" + urlValues.Encode()
	}

	client := c.client()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("openweather status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) client() *http.Client {
	if c.HTTPClient == nil {
		return &http.Client{Timeout: c.RequestTimeout}
	}
	if c.RequestTimeout > 0 {
		c.HTTPClient.Timeout = c.RequestTimeout
	}
	return c.HTTPClient
}
