package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/bbengfort/cosmos/pkg"
	pb "github.com/bbengfort/cosmos/pkg/pb/v1alpha"
	cli "github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	client pb.CosmosClient
)

func main() {
	app := cli.NewApp()
	app.Name = "cosmos"
	app.Usage = "interact with the cosmos game server"
	app.Version = pkg.Version()
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "no-secure",
			Aliases: []string{"S"},
			Usage:   "don't connect with TLS (e.g. for development)",
			EnvVars: []string{"COSMOS_NOTLS", "COSMOS_CLIENT_INSECURE"},
			Value:   false,
		},
		&cli.StringFlag{
			Name:    "endpoint",
			Aliases: []string{"u"},
			Usage:   "endpoint to connect to the Cosmos service on",
			EnvVars: []string{"COSMOS_CLIENT_URL"},
			Value:   "localhost:8088",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:     "login",
			Usage:    "login to the cosmos service",
			Category: "client",
			Before:   initClient,
			Action:   login,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "username",
					Aliases: []string{"u"},
					Usage:   "the username to login with",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "the password to login with",
				},
			},
		},
	}

	app.Run(os.Args)
}

//===========================================================================
// Client Actions
//===========================================================================

func login(c *cli.Context) (err error) {
	auth := &pb.Auth{
		Username: c.String("username"),
		Password: c.String("password"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var tokens *pb.AuthToken
	if tokens, err = client.Login(ctx, auth); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(tokens)
}

//===========================================================================
// Helper Functions
//===========================================================================

// helper function to create the GRPC client with default options
func initClient(c *cli.Context) (err error) {
	var opts []grpc.DialOption
	if c.Bool("no-secure") {
		opts = append(opts, grpc.WithInsecure())
	} else {
		config := &tls.Config{}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	}

	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(c.String("endpoint"), opts...); err != nil {
		return cli.Exit(err, 1)
	}
	client = pb.NewCosmosClient(cc)
	return nil
}

// helper function to print JSON response and exit
func printJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Println(string(data))
	return nil
}
