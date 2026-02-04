package weather

import "testing"

func TestParseWeatherNow(t *testing.T) {
	response := map[string]any{
		"main": map[string]any{
			"temp":       10.4,
			"feels_like": 8.7,
			"humidity":   65.0,
		},
		"weather": []any{
			map[string]any{
				"id":          802.0,
				"description": "scattered clouds",
			},
		},
		"name": "Moscow",
	}

	weather, err := ParseWeatherNow(response)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if weather.Location != "Moscow" {
		t.Fatalf("expected location Moscow, got %q", weather.Location)
	}
	if weather.WeatherID != 802 {
		t.Fatalf("expected weather id 802, got %d", weather.WeatherID)
	}
	if weather.Temp == "" || weather.FeelsLike == "" || weather.Humidity == "" {
		t.Fatalf("expected formatted values, got %+v", weather)
	}
}

func TestParseWeatherForecastDaily(t *testing.T) {
	response := map[string]any{
		"list": []any{
			map[string]any{
				"dt_txt": "2026-02-04 00:00:00",
				"main":   map[string]any{"temp": 10.0, "feels_like": 9.0, "humidity": 60.0},
			},
			map[string]any{
				"dt_txt": "2026-02-04 03:00:00",
				"main":   map[string]any{"temp": 12.0, "feels_like": 11.0, "humidity": 66.0},
			},
			map[string]any{
				"dt_txt": "2026-02-05 00:00:00",
				"main":   map[string]any{"temp": 8.0, "feels_like": 7.0, "humidity": 70.0},
			},
			map[string]any{
				"dt_txt": "2026-02-05 03:00:00",
				"main":   map[string]any{"temp": 6.0, "feels_like": 5.0, "humidity": 74.0},
			},
		},
	}

	rows, err := ParseWeatherForecast(response, 2, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0].Date != "2026-02-04" {
		t.Fatalf("expected date 2026-02-04, got %q", rows[0].Date)
	}
}

func TestParseWeatherForecastWithTime(t *testing.T) {
	response := map[string]any{
		"list": []any{
			map[string]any{
				"dt_txt": "2026-02-04 00:00:00",
				"main":   map[string]any{"temp": 10.0, "feels_like": 9.0, "humidity": 60.0},
			},
			map[string]any{
				"dt_txt": "2026-02-04 03:00:00",
				"main":   map[string]any{"temp": 12.0, "feels_like": 11.0, "humidity": 66.0},
			},
			map[string]any{
				"dt_txt": "2026-02-05 00:00:00",
				"main":   map[string]any{"temp": 8.0, "feels_like": 7.0, "humidity": 70.0},
			},
		},
	}

	rows, err := ParseWeatherForecast(response, 1, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected rows, got 0")
	}
	if rows[0].Time == "" {
		t.Fatalf("expected time to be set")
	}
}
