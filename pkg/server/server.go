package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/bbengfort/cosmos/pkg"
	pb "github.com/bbengfort/cosmos/pkg/pb/v1alpha"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// Server implements the GRPC Cosmos Service
type Server struct {
	pb.UnimplementedCosmosServer
	srv   *grpc.Server
	echan chan error
}

// New creates a new Cosmos server.
func New() *Server {
	return &Server{
		echan: make(chan error, 1),
	}
}

// Serve GRPC requests.
func (s *Server) Serve(addr string) (err error) {
	// Initialize the gRPC server
	s.srv = grpc.NewServer()
	pb.RegisterCosmosServer(s.srv, s)

	// Catch OS signals for graceful shutdowns
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		s.echan <- s.Shutdown()
	}()

	// Listen for TCP requests on the specified address and port
	var sock net.Listener
	if sock, err = net.Listen("tcp", addr); err != nil {
		return fmt.Errorf("could not listen on %q", addr)
	}
	defer sock.Close()

	// Run the server
	go func() {
		log.Info().Str("listen", addr).Str("version", pkg.Version()).Msg("server started")
		if err := s.srv.Serve(sock); err != nil {
			s.echan <- err
		}
	}()

	// Listen for any errors that might have occurred and wait for all go routines to finish
	if err = <-s.echan; err != nil {
		return err
	}
	return nil
}

// Shutdown the Cosmos server gracefully.
func (s *Server) Shutdown() (err error) {
	return nil
}
