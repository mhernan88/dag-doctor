package main

import (
    "os"
    "fmt"
    "log"

    "github.com/urfave/cli/v2"
    "github.com/sirupsen/logrus"
)

const version = "v0.1.0"

var flags = []cli.Flag{
    &cli.BoolFlag{
        Name: "verbose",
        Aliases: []string{"v"},
        Value: false,
    },
    &cli.StringFlag{
        Name: "pipeline",
        Aliases: []string{"p"},
        Value: "dag.json",
    },
    &cli.StringFlag{
        Name: "catalog",
        Aliases: []string{"c"},
        Value: "catalog.json",
    },
}

func action(c *cli.Context) error {
    l := logrus.New()
    if c.Bool("verbose") {
        l.SetLevel(logrus.TraceLevel)
    }

    fmt.Printf("dag-bisect %s\n", version)
    return nil
}

func main() {
    app := &cli.App{
        Name: "DAG Bisect",
        Usage: "Recursively bisect a DAG to quickly locate data errors",
        Flags: flags,
        Action: action,
    }

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}
