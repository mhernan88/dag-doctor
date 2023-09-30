package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/mhernan88/dag-bisect/cmd"
	"github.com/mhernan88/dag-bisect/cmd/sessions"
	"github.com/mhernan88/dag-bisect/cmd/telemetry"
	"github.com/mhernan88/dag-bisect/models"
	"github.com/mhernan88/dag-bisect/db"
	"github.com/mhernan88/dag-bisect/pruners"
	"github.com/mhernan88/dag-bisect/shared"
	"github.com/mhernan88/dag-bisect/splitters"
	"github.com/urfave/cli/v2"
)

const version = "v0.1.0"

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "dag",
		Aliases: []string{"d"},
		Value:   "dag.json",
		Usage:   "filename of serialized dag",
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
	dag, err := models.LoadDAG(c.String("dag"))
	if err != nil {
		return err
	}

	l.Info("loaded root nodes (+ additional child nodes) from dag", "n", len(dag.Nodes))
	if len(dag.Nodes) == 0 {
		return fmt.Errorf("loaded an empty dag!")
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
		fmt.Printf("created '%s' folder.\n", configFolder)
	}

	dbHandle, err := db.Connect()
	if err != nil {
		os.Exit(1)
	}
	err = db.CreateTables(dbHandle, false)
	if err != nil {
		os.Exit(1)
	}
	defer dbHandle.Close()

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
