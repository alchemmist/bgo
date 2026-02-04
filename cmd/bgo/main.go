package main

import (
	"os"

	"bgo/internal/app"
)

func main() {
	code := app.Run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(code)
}
