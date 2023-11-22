package main

import (
	"log"
	"os"
	"text/tabwriter"

	"github.com/bbengfort/cosmos/pkg"
	"github.com/bbengfort/cosmos/pkg/config"
	"github.com/joho/godotenv"
	confire "github.com/rotationalio/confire/usage"
	"github.com/urfave/cli/v2"
)

func main() {
	// If a dotenv file exists load it for configuration
	godotenv.Load()

	// Create a multi-command CLI application
	app := cli.NewApp()
	app.Name = "cosmos"
	app.Version = pkg.Version()
	app.Usage = "cosmos API service and utilities"
	app.Flags = []cli.Flag{}
	app.Commands = []*cli.Command{
		{
			Name:   "serve",
			Usage:  "serve the cosmos API service",
			Action: serve,
		},
		{
			Name:     "config",
			Usage:    "print cosmos configuration guide",
			Category: "utility",
			Action:   usage,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "print in list mode instead of table mode",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func serve(c *cli.Context) error {
	return nil
}

func usage(c *cli.Context) (err error) {
	tabs := tabwriter.NewWriter(os.Stdout, 1, 0, 4, ' ', 0)
	format := confire.DefaultTableFormat
	if c.Bool("list") {
		format = confire.DefaultListFormat
	}

	var conf config.Config
	if err := confire.Usagef(config.Prefix, &conf, tabs, format); err != nil {
		return cli.Exit(err, 1)
	}
	tabs.Flush()
	return nil
}
