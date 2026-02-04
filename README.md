# BGO

![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white) ![Build](https://img.shields.io/badge/build-make%20build-4c9a2a) ![Tests](https://img.shields.io/badge/tests-go%20test-007ec6) ![License](https://img.shields.io/badge/license-MIT-6c757d)

A focused command-line weather client that delivers current conditions and forecasts fast. Built for reliability, clear output, and disciplined behavior under failure conditions. This project is intentionally compact, test-covered, and ready to merge into a main branch without cleanup debt.

## Demo

<video src="media/example_usage.mp4" controls muted playsinline width="800"></video>

## Features
- Current weather and multi-day forecast
- Optional time-granular forecast output
- High-precision mode or rounded output
- Full API response printing for debugging
- Clean ANSI output with ASCII weather art

## Requirements
- Go 1.22+
- `OPEN_WEATHER_API_KEY` in your environment (or a `.env` file in the project root). See `.env.example`

## Quick Start
```bash
make build
./bin/bgo now
```

## Usage
```bash
# current weather
./bin/bgo now

# 3-day forecast
./bin/bgo forecast -d 3

# forecast with time
./bin/bgo forecast -d 2 --with-time

# high precision
./bin/bgo now --high-precision

# show full API response (rounded unless --high-precision)
./bin/bgo now --full-info
```

## Development
```bash
make fmt
make test
make test-watch
```
