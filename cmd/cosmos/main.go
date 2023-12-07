package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/bbengfort/cosmos/pkg"
	"github.com/bbengfort/cosmos/pkg/auth"
	"github.com/bbengfort/cosmos/pkg/config"
	"github.com/bbengfort/cosmos/pkg/cosmos"
	"github.com/bbengfort/cosmos/pkg/db"
	"github.com/bbengfort/cosmos/pkg/db/models"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
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
		{
			Name:     "auth:tokenkey",
			Usage:    "generate an RSA token key pair and ulid for JWT token signing",
			Category: "utility",
			Action:   authTokenKey,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o"},
					Usage:   "path to write keys out to (optional, will be saved as ulid.pem by default)",
				},
				&cli.IntFlag{
					Name:    "size",
					Aliases: []string{"s"},
					Usage:   "number of bits for the generated keys",
					Value:   4096,
				},
			},
		},
		{
			Name:     "auth:createsuperuser",
			Usage:    "create an admin user",
			Category: "utility",
			Action:   authCreateSuperUser,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Usage:    "The name of the admin user",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "email",
					Aliases:  []string{"e"},
					Usage:    "The email of the admin user",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "password",
					Aliases:  []string{"p"},
					Usage:    "The password of the admin user",
					Required: true,
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func serve(c *cli.Context) (err error) {
	var conf config.Config
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}

	var srv *cosmos.Server
	if srv, err = cosmos.New(conf); err != nil {
		return cli.Exit(err, 1)
	}

	if err = srv.Serve(); err != nil {
		return cli.Exit(err, 1)
	}
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

func authTokenKey(c *cli.Context) (err error) {
	keyID := ulid.Make()

	var key *rsa.PrivateKey
	if key, err = rsa.GenerateKey(rand.Reader, c.Int("size")); err != nil {
		return cli.Exit(err, 1)
	}

	var out string
	if out = c.String("out"); out == "" {
		out = fmt.Sprintf("%s.pem", keyID)
	}

	var f *os.File
	if f, err = os.OpenFile(out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600); err != nil {
		return cli.Exit(err, 1)
	}
	defer f.Close()

	if err = pem.Encode(f, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Printf("RSA key id %s -- saved with PEM encoding to %s\n", keyID, out)
	return nil
}

func authCreateSuperUser(c *cli.Context) (err error) {

	var conf config.Config
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}

	if err = db.Connect(conf.Database); err != nil {
		return cli.Exit(err, 1)
	}

	user := &models.User{
		Name:  sql.NullString{Valid: true, String: c.String("name")},
		Email: c.String("email"),
	}

	if user.Password, err = auth.CreateDerivedKey(c.String("password")); err != nil {
		return cli.Exit(err, 1)
	}

	// Create the user
	ctx := context.Background()
	if err = models.CreateUser(ctx, user); err != nil {
		return cli.Exit(err, 1)
	}

	// Update the user role in the database
	var tx *sqlx.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false}); err != nil {
		return cli.Exit(err, 1)
	}
	defer tx.Rollback()

	if _, err = tx.Exec("UPDATE users SET role_id=1 WHERE id=$1", user.ID); err != nil {
		return cli.Exit(err, 1)
	}

	if err = tx.Commit(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}
