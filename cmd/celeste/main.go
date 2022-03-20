package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bbengfort/cosmos/pkg"
	"github.com/bbengfort/cosmos/pkg/config"
	"github.com/bbengfort/cosmos/pkg/db/schema"
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
			Name:     "migrate",
			Usage:    "migrate the database to the latest schema version",
			Category: "database",
			Action:   migrate,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "dsn",
					Aliases: []string{"d", "db"},
					Usage:   "database dsn to connect to the database on",
					EnvVars: []string{"DATABASE_URL", "COSMOS_DATABASE_URL"},
				},
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "force the latest schema version to be applied",
				},
				&cli.BoolFlag{
					Name:    "drop",
					Aliases: []string{"D"},
					Usage:   "drop the database schema before migrating (force must be true)",
				},
			},
		},
		{
			Name:     "schema",
			Usage:    "get the current version of the database schema",
			Category: "database",
			Action:   schemaVersion,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "dsn",
					Aliases: []string{"d", "db"},
					Usage:   "database dsn to connect to the database on",
					EnvVars: []string{"DATABASE_URL", "COSMOS_DATABASE_URL"},
				},
				&cli.BoolFlag{
					Name:    "verify",
					Aliases: []string{"v"},
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

//===========================================================================
// Database Actions
//===========================================================================

func migrate(c *cli.Context) (err error) {
	if err = schema.Configure(c.String("dsn")); err != nil {
		return cli.Exit(err, 1)
	}
	defer schema.Close()

	if c.Bool("drop") {
		if !c.Bool("force") {
			return cli.Exit("cannot drop without forcing", 1)
		}
		if err = schema.Drop(); err != nil {
			return cli.Exit(err, 1)
		}
	}

	if c.Bool("force") {
		if err = schema.Force(); err != nil {
			return cli.Exit(err, 1)
		}
	} else {
		if err = schema.Migrate(); err != nil {
			return cli.Exit(err, 1)
		}
	}
	return nil
}

func schemaVersion(c *cli.Context) (err error) {
	defer schema.Close()
	if c.Bool("verify") {
		if err = schema.Verify(c.String("dsn")); err != nil {
			return cli.Exit(err, 1)
		}
	}

	var vers *schema.Version
	if vers, err = schema.CurrentVersion(c.String("dsn")); err != nil {
		return cli.Exit(err, 1)
	}
	return printJSON(vers)
}

//===========================================================================
// Administrative Actions
//===========================================================================

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

//===========================================================================
// Helper Functions
//===========================================================================

// helper function to print JSON response and exit
func printJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Println(string(data))
	return nil
}
