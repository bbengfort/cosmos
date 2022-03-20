package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bbengfort/cosmos/pkg"
	"github.com/bbengfort/cosmos/pkg/config"
	"github.com/bbengfort/cosmos/pkg/server"
	"github.com/joho/godotenv"
	"github.com/segmentio/ksuid"
	cli "github.com/urfave/cli/v2"
)

func main() {
	// Load the dotenv file
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "celeste"
	app.Usage = "management commands for a cosmos game server"
	app.Version = pkg.Version()
	app.Commands = []*cli.Command{
		{
			Name:     "serve",
			Usage:    "run a cosmos game server",
			Category: "server",
			Action:   serve,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "addr",
					Aliases: []string{"a"},
					Usage:   "specify the address and port to bind the server on",
					Value:   ":10001",
				},
			},
		},
		{
			Name:     "tokenkey",
			Usage:    "generate an RSA token key pair and ksuid for JWT token signing",
			Category: "admin",
			Action:   generateTokenKey,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o"},
					Usage:   "directory to write keys to (optional, will be saved as ksuid.pem)",
				},
				&cli.IntFlag{
					Name:    "size",
					Aliases: []string{"s"},
					Usage:   "number of bits for the generated keys",
					Value:   4096,
				},
			},
		},
	}

	app.Run(os.Args)
}

//===========================================================================
// Server Actions
//===========================================================================

func serve(c *cli.Context) (err error) {
	var conf config.Config
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}

	// Update the configuration from the CLI flags
	if addr := c.String("addr"); addr != "" {
		conf.BindAddr = addr
	}

	var srv *server.Server
	if srv, err = server.New(conf); err != nil {
		return cli.Exit(err, 1)
	}

	if err = srv.Serve(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func generateTokenKey(c *cli.Context) (err error) {
	var keyid ksuid.KSUID
	if keyid, err = ksuid.NewRandom(); err != nil {
		return cli.Exit(err, 1)
	}

	out := fmt.Sprintf("%s.pem", keyid)
	if path := c.String("out"); path != "" {
		out = filepath.Join(path, out)
	}

	var key *rsa.PrivateKey
	if key, err = rsa.GenerateKey(rand.Reader, c.Int("size")); err != nil {
		return cli.Exit(err, 1)
	}

	var f *os.File
	if f, err = os.OpenFile(out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600); err != nil {
		return cli.Exit(err, 1)
	}

	if err = pem.Encode(f, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Printf("RSA key id: %s -- saved with PEM encoding to %s\n", keyid, out)
	return nil
}
