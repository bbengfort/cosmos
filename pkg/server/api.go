package server

import (
	"context"
	"time"

	"github.com/bbengfort/cosmos/pkg"
	pb "github.com/bbengfort/cosmos/pkg/pb/v1alpha"
	"github.com/rs/zerolog/log"
)

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
