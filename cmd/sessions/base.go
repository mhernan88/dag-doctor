package sessions

import "github.com/urfave/cli/v2"

var SessionsCmd = cli.Command{
	Name: "session",
	Usage: "session commands",
	Subcommands: []*cli.Command{
		&ListCmd,
	},
}
