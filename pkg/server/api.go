package server

import (
	"context"
	"time"

	"github.com/bbengfort/cosmos/pkg"
	"github.com/bbengfort/cosmos/pkg/auth"
	pb "github.com/bbengfort/cosmos/pkg/pb/v1alpha"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Login(ctx context.Context, in *pb.Auth) (out *pb.AuthToken, err error) {
	var user *auth.User
	if user, err = auth.Authenticate(ctx, in.Username, in.Password); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	var accessToken *jwt.Token
	if accessToken, err = s.tokens.CreateAccessToken(user.Email); err != nil {
		log.Error().Err(err).Msg("could not create access token")
		return nil, status.Error(codes.FailedPrecondition, "unable to authenticate user")
	}

	var refreshToken *jwt.Token
	if refreshToken, err = s.tokens.CreateRefreshToken(accessToken); err != nil {
		log.Error().Err(err).Msg("could not create refresh token")
		return nil, status.Error(codes.FailedPrecondition, "unable to authenticate user")
	}

	out = &pb.AuthToken{}
	if out.AccessToken, err = s.tokens.Sign(accessToken); err != nil {
		log.Error().Err(err).Msg("could not sign access token")
		return nil, status.Error(codes.FailedPrecondition, "unable to authenticate user")
	}
	if out.RefreshToken, err = s.tokens.Sign(refreshToken); err != nil {
		log.Error().Err(err).Msg("could not sign refresh token")
		return nil, status.Error(codes.FailedPrecondition, "unable to authenticate user")
	}

	return out, nil
}

const (
	statusOk          = "ok"
	statusMaintenance = "maintenance"
)

func (s *Server) Status(ctx context.Context, in *pb.HealthCheck) (out *pb.SystemStatus, err error) {
	out = &pb.SystemStatus{
		Status:  statusOk,
		Version: pkg.Version(),
		Uptime:  time.Since(s.started).String(),
	}

	if s.conf.Maintenance {
		out.Status = statusMaintenance
	}

	log.Info().Uint32("attempts", in.Attempts).Time("last_checked_at", in.LastCheckedAt.AsTime()).Msg("health check")
	return out, nil
}
