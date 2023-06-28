package cmd

import (
	"fmt"

	"github.com/application-research/edge-vertex/core"
	"github.com/application-research/edge-vertex/util"
	"github.com/urfave/cli/v2"
)

func SetupCommands() []*cli.Command {
	var commands []*cli.Command

	/* daemon command */
	commands = append(commands, &cli.Command{
		Name:    "daemon",
		Aliases: []string{"d"},
		Usage:   "run the delta-importer daemon to continuously import deals",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "ddm-api",
				Usage:       "address of ddm instance",
				DefaultText: "http://127.0.0.1:1415",
				Value:       "http://127.0.0.1:1415",
				EnvVars:     []string{"DDM_URL"},
			},
			&cli.StringFlag{
				Name:     "ddm-token",
				Usage:    "ddm auth token",
				Required: true,
				EnvVars:  []string{"DDM_TOKEN"},
			},
			&cli.StringFlag{
				Name:        "edge-file",
				Usage:       "filename containing edge addresses",
				DefaultText: "edges.json",
				Value:       "edges.json",
				EnvVars:     []string{"EDGE_FILE"},
			},
			&cli.IntFlag{
				Name:        "interval",
				Usage:       "interval in seconds between each run",
				DefaultText: "300",
				Value:       300,
				EnvVars:     []string{"INTERVAL"},
			},
		},

		Action: func(cctx *cli.Context) error {
			logo := `Edge Vertex`
			fmt.Println(util.Purple + logo + util.Reset)
			fmt.Printf("\n--\n")
			fmt.Println("API is available at" + util.Red + " 127.0.0.1:" + cctx.String("port") + util.Reset)

			interval := cctx.Int("interval")
			edgeListFilename := cctx.String("edge-file")

			daemon := core.NewEdgeDaemon(interval, edgeListFilename)

			daemon.Run()
			return nil
		},
	})

	return commands
}
