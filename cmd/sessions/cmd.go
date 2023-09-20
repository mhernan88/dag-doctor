package sessions

import "github.com/urfave/cli/v2"

var SessionsDevCmd = cli.Command{
	Name: "dev",
	Usage: "session dev commands",
	Subcommands: []*cli.Command{
		// &UpdateSessionCmd,
	},
}

var SessionsCmd = cli.Command{
	Name: "session",
	Usage: "session commands",
	Subcommands: []*cli.Command{
		&IterSessionCmd,
		&ListSessionsCmd,
		&NewSessionCmd,
		&SessionsDevCmd,
	},
}

