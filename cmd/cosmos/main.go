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
	"github.com/joho/godotenv"
	cli "github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	client pb.CosmosClient
)

func main() {
	// Load the dotenv file
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "cosmos"
	app.Usage = "interact with the cosmos game server"
	app.Version = pkg.Version()
	app.Before = initClient
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
			Value:   "localhost:10001",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:     "login",
			Usage:    "login to the cosmos service",
			Category: "client",
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
		{
			Name:     "status",
			Usage:    "check the status of the cosmos service",
			Category: "client",
			Action:   status,
			Flags: []cli.Flag{
				&cli.Uint64Flag{
					Name:    "attempts",
					Aliases: []string{"a"},
					Usage:   "specify a number of attempts to send to the server (optional)",
				},
				&cli.TimestampFlag{
					Name:    "last-checked-at",
					Aliases: []string{"t", "timestamp"},
					Usage:   "specify a last checked at timestamp (optional)",
					Layout:  time.RFC3339,
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

func status(c *cli.Context) (err error) {
	req := &pb.HealthCheck{
		Attempts: uint32(c.Uint64("attempts")),
	}

	if ts := c.Timestamp("last-checked-at"); ts != nil && !ts.IsZero() {
		req.LastCheckedAt = timestamppb.New(*ts)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *pb.SystemStatus
	if rep, err = client.Status(ctx, req); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

//===========================================================================
// Helper Functions
//===========================================================================

// helper function to create the GRPC client with default options
func initClient(c *cli.Context) (err error) {
	var opts []grpc.DialOption

	var creds credentials.TransportCredentials
	if c.Bool("no-secure") {
		creds = insecure.NewCredentials()
	} else {
		config := &tls.Config{}
		creds = credentials.NewTLS(config)
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))

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
