package view

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"bgo/internal/weather"
)

func TestSelectASCIIArtAndColorDayNight(t *testing.T) {
	original := nowFunc
	defer func() { nowFunc = original }()

	nowFunc = func() time.Time {
		return time.Date(2026, 2, 4, 12, 0, 0, 0, time.UTC)
	}
	art, _ := SelectASCIIArtAndColor(800)
	if art != clearSunny {
		t.Fatalf("expected sunny art during day")
	}

	nowFunc = func() time.Time {
		return time.Date(2026, 2, 4, 2, 0, 0, 0, time.UTC)
	}
	art, _ = SelectASCIIArtAndColor(800)
	if art != clearNight {
		t.Fatalf("expected night art during night")
	}
}

func TestPrintWeatherNow(t *testing.T) {
	buf := &bytes.Buffer{}
	weatherNow := weather.WeatherNow{
		Temp:               "10°C",
		FeelsLike:          "8°C",
		Humidity:           "60%",
		WeatherID:          802,
		WeatherDescription: "scattered clouds",
		Location:           "Testville",
	}

	PrintWeatherNow(buf, weatherNow)
	if !strings.Contains(buf.String(), "Testville") {
		t.Fatalf("expected location in output")
	}
}
