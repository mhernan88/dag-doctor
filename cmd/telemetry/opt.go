package telemetry

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/mhernan88/dag-bisect/shared"
	"github.com/urfave/cli/v2"
)

type TelemetryStatus struct {
	Opt string `json:"opt"`
	Date int64 `json:"time"`
}

func Opt(ctx *cli.Context) error {
	status := ctx.Args().Get(0)
	if (status != "in") && (status != "out") {
		fmt.Println("ERROR: opt must be followed by either 'in' or 'out'")
	}

	opts := slog.HandlerOptions{
		AddSource: true,
		Level: slog.LevelDebug,
	}

	logFilename, err := shared.GetLogFilename()
	if err != nil {
		return err
	}
	logfile, err := os.OpenFile(logFilename, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		slog.Error("error opening log file", err)
	}
	defer logfile.Close()
	l := slog.New(slog.NewJSONHandler(logfile, &opts))

	optFilename, err := shared.GetOptFilename()
	if err != nil {
		return err
	}
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

	l.Info(fmt.Sprintf("telemetry sharing opt-%s", status))
	fmt.Printf("saved opt-%s status to '%s'\n", status, optFilename)
	return nil
}

var OptCmd = cli.Command{
	Name: "opt",
	Usage: "opt in/out of telemetry",
	Action: Opt,
}
