package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"bgo/internal/utils"
	"bgo/internal/view"
	"bgo/internal/weather"
)

type Runner struct {
	Client WeatherClient
	Out    io.Writer
	In     io.Reader
	ErrOut io.Writer
}

type WeatherClient interface {
	GetCoordinates(ctx context.Context) (weather.Coordinates, error)
	GetWeatherNow(ctx context.Context, coords weather.Coordinates) (map[string]any, error)
	GetWeatherForecast(ctx context.Context, coords weather.Coordinates) (map[string]any, error)
}

func Run(args []string, out io.Writer, errOut io.Writer) int {
	utils.LoadDotEnv(".env")

	runner := &Runner{
		Client: weather.NewClient(os.Getenv("OPEN_WEATHER_API_KEY")),
		Out:    out,
		In:     os.Stdin,
		ErrOut: errOut,
	}

	return runner.Run(args)
}

func (r *Runner) Run(args []string) int {
	fs := flag.NewFlagSet("bgo", flag.ContinueOnError)
	fs.SetOutput(r.ErrOut)

	days := fs.Int("days", 5, "set how long the forecast you want to see (from 1 to 5 days)")
	fs.IntVar(days, "d", 5, "set how long the forecast you want to see (from 1 to 5 days)")
	highPrecision := fs.Bool("high-precision", false, "use this field for show value with max precision")
	fullInfo := fs.Bool("full-info", false, "use this field for show all information")
	withTime := fs.Bool("with-time", false, "use this field for show forecast with time")

	fs.Usage = func() {
		view.PrintHelp(r.ErrOut)
	}

	command := "now"
	flagArgs := args
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		command = args[0]
		flagArgs = args[1:]
	}

	if err := fs.Parse(flagArgs); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}

	if command == "now" && fs.NArg() > 0 {
		command = fs.Arg(0)
	}

	if *days < 1 || *days > 5 {
		fmt.Fprintln(r.Out, "Oops, I don't know what to do! Use -h to see the usage.")
		return 2
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	coords, err := r.Client.GetCoordinates(ctx)
	if err != nil {
		handleError(r.In, r.Out, err, "Oops! Something went wrong, couldn't get coordinates")
		return 1
	}

	switch command {
	case "now":
		response, err := r.Client.GetWeatherNow(ctx, coords)
		if err != nil {
			handleError(r.In, r.Out, err, "Oops! Please check your internet connection")
			return 1
		}
		if !*highPrecision {
			response = utils.RoundJSON(response).(map[string]any)
		}
		if *fullInfo {
			_ = utils.PrintJSON(r.Out, response)
			return 1
		}
		weatherNow, err := weather.ParseWeatherNow(response)
		if err != nil {
			handleError(r.In, r.Out, err, "Oops! Something went wrong, couldn't parse the response")
			return 1
		}
		view.PrintWeatherNow(r.Out, weatherNow)

	case "forecast":
		response, err := r.Client.GetWeatherForecast(ctx, coords)
		if err != nil {
			handleError(r.In, r.Out, err, "Oops! Please check your internet connection")
			return 1
		}
		if !*highPrecision {
			response = utils.RoundJSON(response).(map[string]any)
		}
		if *fullInfo {
			_ = utils.PrintJSON(r.Out, response)
			return 1
		}
		rows, err := weather.ParseWeatherForecast(response, *days, *withTime, *highPrecision)
		if err != nil {
			handleError(r.In, r.Out, err, "Oops! Something went wrong, couldn't parse the response")
			return 1
		}
		view.PrintWeatherForecast(r.Out, rows, *withTime)

	default:
		fmt.Fprintln(r.Out, "Oops, I don't know what to do! Use -h to see the usage.")
		return 2
	}

	return 0
}
