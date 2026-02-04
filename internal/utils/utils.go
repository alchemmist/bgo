package utils

import (
	"bufio"
	"encoding/json"
	"io"
	"math"
	"os"
	"strings"
)

func RoundJSON(value any) any {
	switch v := value.(type) {
	case map[string]any:
		for key, val := range v {
			v[key] = RoundJSON(val)
		}
		return v
	case []any:
		for i, val := range v {
			v[i] = RoundJSON(val)
		}
		return v
	case float64:
		return math.Round(v)
	default:
		return value
	}
}

func PrintJSON(out io.Writer, data map[string]any) error {
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		_, werr := out.Write([]byte("invalid json\n"))
		if werr != nil {
			return werr
		}
		return err
	}
	_, err = out.Write(append(raw, '\n'))
	return err
}

func LoadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, "\"'")
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}
