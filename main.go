package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mhernan88/dag-bisect/cmd"
	"github.com/mhernan88/dag-bisect/data"
	"github.com/mhernan88/dag-bisect/pruners"
	"github.com/mhernan88/dag-bisect/splitters"
	"github.com/sirupsen/logrus"
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
	l := logrus.New()
	l.SetLevel(logrus.WarnLevel)

	if c.Bool("v") {
		l.SetLevel(logrus.InfoLevel)
		l.Info("logging set to INFO level")
	}

	if c.Bool("vv") {
		l.SetLevel(logrus.DebugLevel)
		l.Info("logging set to DEBUG level")
	}

	if c.Bool("vvv") {
		l.SetLevel(logrus.TraceLevel)
		l.Info("logging set to TRACE level")
	}

	fmt.Printf("dag-bisect %s\n", version)

	l.Debug("initializing splitter")
	splitter := splitters.NewDefaultSplitter(
		c.Int("iteration_limit"),
		l,
	)

	l.Debug("initializing pruner")
	pruner := pruners.NewDefaultPruner(
		c.Int("iteration_limit"),
		l,
	)

	l.Debug("loading dag")
	dag, err := data.LoadDAG(c.String("dag"))
	if err != nil {
		return err
	}
	if dag == nil {
		return fmt.Errorf("dag wil nil")
	}
	l.Infof("loaded %d root nodes (+ additional child nodes) from dag", len(dag))
	if len(dag) == 0 {
		return fmt.Errorf("failed to load dag")
	}
	l.Tracef("dag: %v", dag)

	ui := cmd.NewUI(dag, splitter, pruner, c.Int("iteration_limit"), l)
	return ui.Run()
}

func main() {
	app := &cli.App{
		Name:   "DAG Bisect",
		Usage:  "Recursively bisect a DAG to quickly locate data errors",
		Flags:  flags,
		Action: action,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
