package main

import (
	"log"
	"os"

	"github.com/bbengfort/cosmos/pkg"
	"github.com/joho/godotenv"
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
	app.Before = configure
	app.Commands = []*cli.Command{
		{
			Name:   "serve",
			Usage:  "serve the cosmos API service",
			Action: serve,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func configure(c *cli.Context) error {
	return nil
}

func serve(c *cli.Context) error {
	return nil
}
