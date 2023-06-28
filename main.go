package main

import (
	"fmt"
	"os"

	"github.com/application-research/edge-vertex/cmd"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
)

var (
	log = logging.Logger("api")
)

var Commit string
var Version string

func main() {

	commands := cmd.SetupCommands()

	app := &cli.App{
		Commands: commands,
		Usage:    "An application that aggregates available contents from Edge-URIDs",
		Version:  fmt.Sprintf("%s+git.%s\n", Version, Commit),
		// Flags:    cmd.CLIConnectFlags,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
