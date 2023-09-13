package telemetry

import (
	"fmt"
	"log/slog"
	"time"
	"encoding/json"
	"os"

	"github.com/urfave/cli/v2"
)

const configFolder = "~/.config/dag_doctor"
var logFilename = fmt.Sprintf("%s/log.json", configFolder)
var optFilename = fmt.Sprintf("%s/opt.json", configFolder)

type TelemetryStatus struct {
	Opt string `json:"opt"`
	Date int64 `json:"time"`
}

func Opt(ctx *cli.Context) error {
	status := ctx.Args().Get(0)
	if (status != "in") && (status != "out") {
		fmt.Println("ERROR: opt must be followed by either 'in' or 'out'")
	}

	if _, err := os.Stat(configFolder); os.IsNotExist(err) {
		os.MkdirAll(configFolder, os.ModePerm)
		fmt.Println("created '~/.config/dag_doctor' folder.")
	}

	opts := slog.HandlerOptions{
		AddSource: true,
		Level: slog.LevelDebug,
	}

	logfile, err := os.OpenFile(logFilename, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		slog.Error("error opening log file", err)
	}
	defer logfile.Close()
	l := slog.New(slog.NewJSONHandler(logfile, &opts))

	f, err := os.OpenFile(optFilename, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		l.Error("error opening opt file", err)
	}
	defer f.Close()

	var telemetry []TelemetryStatus
	decoder := json.NewDecoder(f)
	decoder.Decode(telemetry)
	telemetry = append(telemetry, TelemetryStatus{Opt: status, Date: time.Now().Unix()})
	encoder := json.NewEncoder(f)
	encoder.Encode(telemetry)

	l.Debug(fmt.Sprintf("telemetry sharing opt-%s", status))
	fmt.Printf("saved opt-%s status to '%s'\n", status, optFilename)
	return nil
}

var OptCmd = cli.Command{
	Name: "opt",
	Usage: "opt in/out of telemetry",
	Action: Opt,
}
