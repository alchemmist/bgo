package view

import (
	"fmt"
	"io"
	"strings"
	"time"

	"bgo/internal/weather"
)

var nowFunc = time.Now

func PrintWeatherNow(out io.Writer, weather weather.WeatherNow) {
	asciiArt, color := SelectASCIIArtAndColor(weather.WeatherID)

	timeLine := nowFunc().Format("15:04 PM")

	col1 := []string{
		fmt.Sprintf("%s", timeLine),
		fmt.Sprintf("temperature: %s", weather.Temp),
		fmt.Sprintf("humidity: %s", weather.Humidity),
	}

	col2 := []string{
		capitalize(weather.WeatherDescription),
		fmt.Sprintf("feels like: %s", weather.FeelsLike),
		"source: OpenWeather",
	}

	fmt.Fprintf(out, "%s%s%s\n", colorize(fmt.Sprintf("%s :sun_behind_small_cloud:", weather.Location), color), resetColor(), "")

	left := strings.Split(asciiArt, "\n")
	for i := 0; i < 3; i++ {
		a := ""
		if i < len(left) {
			a = left[i]
		}
		c1 := ""
		if i < len(col1) {
			c1 = col1[i]
		}
		c2 := ""
		if i < len(col2) {
			c2 = col2[i]
		}

		line := fmt.Sprintf("%-10s  %-28s  %s", a, c1, c2)
		fmt.Fprintf(out, "%s%s%s\n", colorize(line, color), resetColor(), "")
	}
}

func PrintWeatherForecast(out io.Writer, rows []weather.ForecastRow, withTime bool) {
	if withTime {
		fmt.Fprintf(out, "%-12s %-6s %-14s %-16s %-10s\n", "Date", "Time", "Temperature", "Feels Like", "Humidity")
		for _, row := range rows {
			fmt.Fprintf(out, "%-12s %-6s %-14s %-16s %-10s\n", row.Date, row.Time, row.Temp, row.FeelsLike, row.Humidity)
		}
		return
	}

	fmt.Fprintf(out, "%-12s %-14s %-16s %-10s\n", "Date", "Temperature", "Feels Like", "Humidity")
	for _, row := range rows {
		fmt.Fprintf(out, "%-12s %-14s %-16s %-10s\n", row.Date, row.Temp, row.FeelsLike, row.Humidity)
	}
}

func SelectASCIIArtAndColor(weatherID int) (string, string) {
	weatherTypeID := weatherID / 100
	weatherStateID := weatherID % 100
	now := nowFunc()

	if weatherID == 800 {
		sunrise := time.Date(now.Year(), now.Month(), now.Day(), SunriseTime, 0, 0, 0, now.Location())
		sunset := time.Date(now.Year(), now.Month(), now.Day(), SunsetTime, 0, 0, 0, now.Location())
		if now.After(sunrise) && now.Before(sunset) {
			return clearSunny, palette["YELLOW"]
		}
		return clearNight, palette["PURPLE"]
	}

	switch weatherTypeID {
	case 2:
		return thunderstorm, palette["RED"]
	case 3:
		return drizzle, palette["LIGHT_BLUE"]
	case 5:
		return rain, palette["BLUE"]
	case 6:
		return snow, palette["WHITE"]
	case 7:
		return fog, palette["DARK_GRAY"]
	case 8:
		if weatherStateID < 3 {
			return partialClouds, palette["LIGHT_GRAY"]
		}
		return clouds, palette["GRAY"]
	default:
		return everythingElse, palette["LIGHT_GRAY"]
	}
}

func colorize(text string, color string) string {
	if color == "" {
		return text
	}
	if strings.HasPrefix(color, "#") && len(color) == 7 {
		r := hexToInt(color[1:3])
		g := hexToInt(color[3:5])
		b := hexToInt(color[5:7])
		return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s", r, g, b, text)
	}
	switch strings.ToLower(color) {
	case "red":
		return "\x1b[31m" + text
	case "green":
		return "\x1b[32m" + text
	default:
		return text
	}
}

func resetColor() string {
	return "\x1b[0m"
}

func hexToInt(hex string) int {
	n := 0
	for _, c := range hex {
		n <<= 4
		switch {
		case c >= '0' && c <= '9':
			n += int(c - '0')
		case c >= 'a' && c <= 'f':
			n += int(c-'a') + 10
		case c >= 'A' && c <= 'F':
			n += int(c-'A') + 10
		}
	}
	return n
}

func capitalize(value string) string {
	if value == "" {
		return value
	}
	runes := []rune(value)
	runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
	return string(runes)
}
