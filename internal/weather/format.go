package weather

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type WeatherNow struct {
	Temp               string
	FeelsLike          string
	Humidity           string
	WeatherID          int
	WeatherDescription string
	Location           string
}

type ForecastRow struct {
	Date      string
	Time      string
	Temp      string
	FeelsLike string
	Humidity  string
}

func ParseWeatherNow(response map[string]any) (WeatherNow, error) {
	mainMap, ok := response["main"].(map[string]any)
	if !ok {
		return WeatherNow{}, errors.New("missing main")
	}

	weatherList, ok := response["weather"].([]any)
	if !ok || len(weatherList) == 0 {
		return WeatherNow{}, errors.New("missing weather")
	}

	weatherItem, ok := weatherList[0].(map[string]any)
	if !ok {
		return WeatherNow{}, errors.New("invalid weather item")
	}

	temp, ok := toFloat(mainMap["temp"])
	if !ok {
		return WeatherNow{}, errors.New("invalid temp")
	}
	feelsLike, ok := toFloat(mainMap["feels_like"])
	if !ok {
		return WeatherNow{}, errors.New("invalid feels_like")
	}
	humidity, ok := toFloat(mainMap["humidity"])
	if !ok {
		return WeatherNow{}, errors.New("invalid humidity")
	}

	weatherID, ok := toInt(weatherItem["id"])
	if !ok {
		return WeatherNow{}, errors.New("invalid weather id")
	}

	description, _ := weatherItem["description"].(string)
	location, _ := response["name"].(string)

	return WeatherNow{
		Temp:               formatNumber(temp) + "°C",
		FeelsLike:          formatNumber(feelsLike) + "°C",
		Humidity:           formatNumber(humidity) + "%",
		WeatherID:          weatherID,
		WeatherDescription: description,
		Location:           location,
	}, nil
}

func ParseWeatherForecast(response map[string]any, days int, withTime bool, highPrecision bool) ([]ForecastRow, error) {
	listRaw, ok := response["list"].([]any)
	if !ok || len(listRaw) == 0 {
		return nil, errors.New("missing list")
	}

	if withTime {
		return parseForecastWithTime(listRaw, days)
	}

	return parseForecastDaily(listRaw, days, highPrecision)
}

func parseForecastDaily(listRaw []any, days int, highPrecision bool) ([]ForecastRow, error) {
	rows := []ForecastRow{}

	firstItem, ok := listRaw[0].(map[string]any)
	if !ok {
		return rows, errors.New("invalid forecast item")
	}

	lastDate, err := getDate(firstItem)
	if err != nil {
		return rows, err
	}

	j := 0
	for d := 0; d < days; d++ {
		params := map[string]float64{"temp": 0, "feels_like": 0, "humidity": 0}
		count := 0

		for j < len(listRaw) {
			item, ok := listRaw[j].(map[string]any)
			if !ok {
				return rows, errors.New("invalid forecast item")
			}
			date, err := getDate(item)
			if err != nil {
				return rows, err
			}
			if date != lastDate {
				break
			}

			mainMap, ok := item["main"].(map[string]any)
			if !ok {
				return rows, errors.New("missing main in forecast")
			}

			temp, _ := toFloat(mainMap["temp"])
			feelsLike, _ := toFloat(mainMap["feels_like"])
			humidity, _ := toFloat(mainMap["humidity"])

			params["temp"] += temp
			params["feels_like"] += feelsLike
			params["humidity"] += humidity
			count++
			j++
		}

		if count == 0 {
			break
		}

		avgTemp := params["temp"] / float64(count)
		avgFeels := params["feels_like"] / float64(count)
		avgHumidity := params["humidity"] / float64(count)

		tempStr := formatNumberPrecision(avgTemp, highPrecision) + " °C"
		feelsStr := formatNumberPrecision(avgFeels, highPrecision) + " °C"
		humStr := formatNumberPrecision(avgHumidity, highPrecision) + "%"

		dateForRow := lastDate
		if j-1 >= 0 {
			if item, ok := listRaw[j-1].(map[string]any); ok {
				if date, err := getDate(item); err == nil {
					dateForRow = date
				}
			}
		}

		rows = append(rows, ForecastRow{
			Date:      dateForRow,
			Temp:      tempStr,
			FeelsLike: feelsStr,
			Humidity:  humStr,
		})

		if j < len(listRaw) {
			if item, ok := listRaw[j].(map[string]any); ok {
				if date, err := getDate(item); err == nil {
					lastDate = date
				}
			}
		} else {
			break
		}
	}

	return rows, nil
}

func parseForecastWithTime(listRaw []any, days int) ([]ForecastRow, error) {
	rows := []ForecastRow{}
	lastDate := ""
	count := 0

	for i := 0; i < len(listRaw); i++ {
		if count == days+1 {
			break
		}

		item, ok := listRaw[i].(map[string]any)
		if !ok {
			return rows, errors.New("invalid forecast item")
		}

		date, tm, err := getDateTime(item)
		if err != nil {
			return rows, err
		}

		mainMap, ok := item["main"].(map[string]any)
		if !ok {
			return rows, errors.New("missing main in forecast")
		}

		temp, _ := toFloat(mainMap["temp"])
		feelsLike, _ := toFloat(mainMap["feels_like"])
		humidity, _ := toFloat(mainMap["humidity"])

		if date == lastDate {
			date = ""
		} else {
			count++
			lastDate = date
		}

		if count != days+1 {
			rows = append(rows, ForecastRow{
				Date:      date,
				Time:      tm,
				Temp:      formatNumber(temp) + " °C",
				FeelsLike: formatNumber(feelsLike) + " °C",
				Humidity:  formatNumber(humidity) + "%",
			})
		}
	}

	return rows, nil
}

func getDate(item map[string]any) (string, error) {
	dt, ok := item["dt_txt"].(string)
	if !ok {
		return "", errors.New("missing dt_txt")
	}
	parts := strings.Split(dt, " ")
	if len(parts) == 0 {
		return "", errors.New("invalid dt_txt")
	}
	return parts[0], nil
}

func getDateTime(item map[string]any) (string, string, error) {
	dt, ok := item["dt_txt"].(string)
	if !ok {
		return "", "", errors.New("missing dt_txt")
	}
	parts := strings.Split(dt, " ")
	if len(parts) < 2 {
		return "", "", errors.New("invalid dt_txt")
	}
	tm := parts[1]
	if len(tm) >= 3 {
		tm = tm[:len(tm)-3]
	}
	return parts[0], tm, nil
}

func formatNumber(value float64) string {
	if value == float64(int64(value)) {
		return fmt.Sprintf("%.0f", value)
	}
	return fmt.Sprintf("%.2f", value)
}

func formatNumberPrecision(value float64, highPrecision bool) string {
	if highPrecision {
		return fmt.Sprintf("%.2f", value)
	}
	return fmt.Sprintf("%.0f", value)
}

func toFloat(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

func toInt(value any) (int, bool) {
	switch v := value.(type) {
	case float64:
		return int(v), true
	case int:
		return v, true
	case int64:
		return int(v), true
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}
		return i, true
	default:
		return 0, false
	}
}
