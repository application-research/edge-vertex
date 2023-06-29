package cmd

import (
	"fmt"

	"github.com/application-research/edge-vertex/core"
	"github.com/application-research/edge-vertex/util"
	log "github.com/sirupsen/logrus"
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
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "set to enable debug logging output",
			},
		},

		Action: func(cctx *cli.Context) error {
			logo := `𝔼𝕕𝕘𝕖 𝕍𝕖𝕣𝕥𝕖𝕩`
			fmt.Println(util.Purple + logo + util.Reset)
			fmt.Printf("\n--\n")
			fmt.Println("running content aggregation every" + util.Red + cctx.String("interval") + util.Reset + "seconds")

			interval := cctx.Int("interval")
			edgeListFilename := cctx.String("edge-file")
			debug := cctx.Bool("debug")
			ddmUrl := cctx.String("ddm-api")
			ddmKey := cctx.String("ddm-token")

			if debug {
				log.SetLevel(log.DebugLevel)
			}

			daemon := core.NewEdgeDaemon(interval, edgeListFilename, ddmUrl, ddmKey)

			daemon.Run()
			return nil
		},
	})

	return commands
}
