package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/mhernan88/dag-bisect/cmd"
	"github.com/mhernan88/dag-bisect/cmd/sessions"
	"github.com/mhernan88/dag-bisect/cmd/telemetry"
	"github.com/mhernan88/dag-bisect/data"
	"github.com/mhernan88/dag-bisect/pruners"
	"github.com/mhernan88/dag-bisect/shared"
	"github.com/mhernan88/dag-bisect/splitters"
	"github.com/urfave/cli/v2"
)

const version = "v0.1.0"

var flags = []cli.Flag{
	&cli.BoolFlag{
		Name:  "v",
		Value: false,
		Usage: "verbose - info level",
	},
	&cli.BoolFlag{
		Name:  "vv",
		Value: false,
		Usage: "verbose - debug level",
	},
	&cli.BoolFlag{
		Name:  "vvv",
		Value: false,
		Usage: "verbose - trace level",
	},
	&cli.StringFlag{
		Name:    "dag",
		Aliases: []string{"d"},
		Value:   "dag.json",
		Usage:   "filename of serialized dag",
	},
	&cli.IntFlag{
		Name:  "iteration_limit",
		Value: 99,
		Usage: "maximum iteration/recursion depth",
	},
}

func action(c *cli.Context) error {
	logFilename, err := shared.GetLogFilename()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(logFilename, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	l := slog.New(slog.NewJSONHandler(f , nil))

	fmt.Printf("dag-bisect %s\n", version)

	l.Debug("initializing splitter")
	splitter := splitters.NewDefaultSplitter()

	l.Debug("initializing pruner")
	pruner := pruners.NewDefaultPruner()

	l.Debug("loading dag")
	dag, err := data.LoadDAG(c.String("dag"))
	if err != nil {
		return err
	}
	if dag == nil {
		return fmt.Errorf("dag wil nil")
	}
	l.Info("loaded root nodes (+ additional child nodes) from dag", "n", len(dag.Nodes))
	if len(dag.Nodes) == 0 {
		return fmt.Errorf("failed to load dag")
	}

	ui := cmd.NewUI(*dag, splitter, pruner)
	return ui.Run(l)
}

func main() {
	configFolder, err := shared.GetConfigFolder()
	if err != nil {
		os.Exit(1)
	}

	if _, err := os.Stat(configFolder); os.IsNotExist(err) {
		os.MkdirAll(configFolder, os.ModePerm)
		fmt.Println("created '~/.config/dag_doctor' folder.")
	}

	app := &cli.App{
		Name:   "DAG Bisect",
		Usage:  "Recursively bisect a DAG to quickly locate data errors",
		Flags:  flags,
		Action: action,
		Commands: []*cli.Command{
			&telemetry.OptCmd,
			&sessions.SessionsCmd,
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
