package utils

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRoundJSON(t *testing.T) {
	input := map[string]any{
		"a": 1.4,
		"b": []any{2.6, map[string]any{"c": 3.2}},
	}

	output := RoundJSON(input).(map[string]any)
	if output["a"].(float64) != 1 {
		t.Fatalf("expected a rounded to 1, got %v", output["a"])
	}
	list := output["b"].([]any)
	if list[0].(float64) != 3 {
		t.Fatalf("expected b[0] rounded to 3, got %v", list[0])
	}
}

func TestLoadDotEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("FOO=bar\n# comment\n"), 0644); err != nil {
		t.Fatalf("write env: %v", err)
	}

	os.Unsetenv("FOO")
	LoadDotEnv(path)

	if os.Getenv("FOO") != "bar" {
		t.Fatalf("expected FOO to be bar")
	}
}

func TestPrintJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	err := PrintJSON(buf, map[string]any{"a": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatalf("expected output")
	}
}
