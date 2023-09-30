package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mhernan88/dag-bisect/cmd/sessions"
	"github.com/mhernan88/dag-bisect/cmd/telemetry"
	"github.com/mhernan88/dag-bisect/db"
	"github.com/mhernan88/dag-bisect/shared"
	"github.com/urfave/cli/v2"
)

const version = "v0.1.0"


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
