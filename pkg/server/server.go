package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/bbengfort/cosmos/pkg"
	"github.com/bbengfort/cosmos/pkg/auth"
	"github.com/bbengfort/cosmos/pkg/config"
	"github.com/bbengfort/cosmos/pkg/db"
	"github.com/bbengfort/cosmos/pkg/db/schema"
	pb "github.com/bbengfort/cosmos/pkg/pb/v1alpha"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// Server implements the GRPC Cosmos Service
type Server struct {
	pb.UnimplementedCosmosServer
	conf    config.Config
	tokens  *auth.TokenManager
	srv     *grpc.Server
	started time.Time
	echan   chan error
}

// New creates a new Cosmos server.
func New(conf config.Config) (s *Server, err error) {
	// Load the default configuration from the environment
	if conf.IsZero() {
		if conf, err = config.New(); err != nil {
			return nil, err
		}
	}

	// Set the global level
	zerolog.SetGlobalLevel(conf.GetLogLevel())

	// Set human readable logging if specified
	if conf.ConsoleLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Create the server and prepare to serve
	s = &Server{conf: conf, echan: make(chan error, 1)}

	// Create the gRPC server options
	opts := make([]grpc.ServerOption, 0, 2)
	opts = append(opts, s.UnaryInterceptors())
	opts = append(opts, s.StreamInterceptors())

	// Create and register the gRPC server
	s.srv = grpc.NewServer(opts...)
	pb.RegisterCosmosServer(s.srv, s)

	// At this point, if we're in maintenance mode, we're done setting up
	if conf.Maintenance {
		return s, nil
	}

	// Perform all setup and configuration that is required when not in maintenance mode
	if s.tokens, err = auth.New(s.conf.Auth.TokenKeys, s.conf.Auth.Audience); err != nil {
		return nil, err
	}
	return s, nil
}

// Serve GRPC requests.
func (s *Server) Serve() (err error) {
	// Catch OS signals for graceful shutdowns
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		s.echan <- s.Shutdown()
	}()

	if s.conf.Maintenance {
		log.Warn().Msg("starting cosmos server in maintenance mode")
	} else {
		// Wait for the database to be at the correct schema
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		err = schema.Wait(ctx, s.conf.Database.URL)
		cancel()
		if err != nil {
			return err
		}

		// Connect to the database
		if err = db.Connect(s.conf.Database); err != nil {
			return err
		}
	}

	// Set the started timestamp for uptime requests
	s.started = time.Now()

	// Listen for TCP requests on the specified address and port
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.conf.BindAddr); err != nil {
		log.Error().Err(err).Str("bindaddr", s.conf.BindAddr).Msg("could not listen on addr")
		return err
	}

	// Run the gRPC server
	go s.Run(sock)
	log.Info().Str("listen", s.conf.BindAddr).Str("version", pkg.Version()).Msg("cosmos server started")

	// Listen for any errors that might have occurred and wait for all go routines to finish
	if err = <-s.echan; err != nil {
		return err
	}
	return nil
}

// Run the gRPC server. This method is extracted from the Serve method so that it can be
// run in its own go routine and allow tests to Run a bufcon server without starting a
// live server with all of the various go routines and channels running and open.
func (s *Server) Run(sock net.Listener) {
	defer sock.Close()
	if err := s.srv.Serve(sock); err != nil {
		s.echan <- err
	}
}

// Shutdown the Cosmos server gracefully.
func (s *Server) Shutdown() (err error) {
	log.Info().Msg("gracefully shutting down the cosmos server")
	s.srv.GracefulStop()

	log.Debug().Msg("successful shutdown of cosmos server")
	return nil
}
