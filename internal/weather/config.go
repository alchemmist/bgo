package weather

const language = "en"

func BasicRequestParams(apiKey string) map[string]string {
	return map[string]string{
		"units": "metric",
		"appid": apiKey,
		"lang":  language,
	}
}
