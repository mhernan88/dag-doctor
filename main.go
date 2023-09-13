package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mhernan88/dag-bisect/data"
	"github.com/mhernan88/dag-bisect/tui"
	"github.com/rivo/tview"

	// "github.com/mhernan88/dag-bisect/pruners"
	// "github.com/mhernan88/dag-bisect/splitters"
	// "github.com/mhernan88/dag-bisect/tui"
	"github.com/sirupsen/logrus"
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
	&cli.IntFlag{
		Name:  "iteration_limit",
		Value: 99,
		Usage: "maximum iteration/recursion depth",
	},
}

func action(c *cli.Context) error {
	l := logrus.New()
	l.SetLevel(logrus.ErrorLevel)

	// splitter := splitters.NewDefaultSplitter(
	// 	c.Int("iteration_limit"),
	// 	l,
	// )
	//
	// pruner := pruners.NewDefaultPruner(
	// 	c.Int("iteration_limit"),
	// 	l,
	// )

	dag, err := data.LoadDAG(c.String("dag"))
	if err != nil {
		return err
	}
	if dag == nil {
		return fmt.Errorf("dag wil nil")
	}
	if len(dag.Nodes) == 0 {
		return fmt.Errorf("failed to load dag")
	}

	// ui := tview.NewBox().SetBorder(true).SetTitle("Hello World")
	grid := tui.CreateLayout()
	err = tview.NewApplication().SetRoot(grid, true).SetFocus(grid).Run()
	if err != nil {
		return err
	}
	return nil
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
