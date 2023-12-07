package cosmos

import (
	"net/http"

	"github.com/bbengfort/cosmos/pkg/api/v1"
	"github.com/bbengfort/cosmos/pkg/auth"
	"github.com/bbengfort/cosmos/pkg/db/models"
	"github.com/bbengfort/cosmos/pkg/enums"
	"github.com/bbengfort/cosmos/pkg/jcode"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	DefaultGameSize = enums.Medium
	DefaultMaxTurns = 1000
)

func (s *Server) ListGalaxies(c *gin.Context) {
	var (
		err      error
		userID   int64
		claims   *auth.Claims
		galaxies []*models.Galaxy
	)

	if claims, err = auth.GetClaims(c); err != nil {
		log.Warn().Err(err).Msg("could not get claims from request")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete list galaxies request"))
		return
	}

	if userID, err = claims.SubjectID(); err != nil {
		log.Warn().Err(err).Msg("could not parse user ID from claims")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete list galaxies request"))
		return
	}

	if galaxies, err = models.ListActiveGalaxies(c.Request.Context(), userID); err != nil {
		log.Error().Err(err).Msg("could not fetch active galaxies from the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete list galaxies request"))
		return
	}

	c.JSON(http.StatusOK, galaxies)
}

func (s *Server) CreateGalaxy(c *gin.Context) {
	var (
		err    error
		galaxy *models.Galaxy
		userID int64
		claims *auth.Claims
		player *models.Player
	)

	if claims, err = auth.GetClaims(c); err != nil {
		log.Warn().Err(err).Msg("could not get claims to create galaxy")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create galaxy"))
		return
	}

	if userID, err = claims.SubjectID(); err != nil {
		log.Warn().Err(err).Msg("could not parse claims to create galaxy")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create galaxy"))
		return
	}

	// Parse the galaxy from the user input
	galaxy = &models.Galaxy{}
	if err = c.BindJSON(galaxy); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Set the defaults for the galaxy before saving it.
	galaxy.JoinCode = jcode.New()
	galaxy.GameState = enums.Pending
	galaxy.MaxTurns = DefaultMaxTurns
	galaxy.Turn = 0

	// Ensure there is a galaxy size
	if galaxy.Size == enums.UnknownSize {
		galaxy.Size = DefaultGameSize
	}

	// Set the max players based on game size
	switch galaxy.Size {
	case enums.Small:
		galaxy.MaxPlayers = 2
	case enums.Medium:
		galaxy.MaxPlayers = 10
	case enums.Large:
		galaxy.MaxPlayers = 20
	case enums.Galactic:
		galaxy.MaxPlayers = 50
	case enums.Cosmic:
		galaxy.MaxPlayers = 100
	}

	// Create the galaxy
	if err = models.CreateGalaxy(c.Request.Context(), galaxy); err != nil {
		log.Error().Err(err).Msg("could not create galaxy")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create galaxy"))
		return
	}

	// Create the player
	player = &models.Player{
		GalaxyID:  galaxy.ID,
		PlayerID:  userID,
		RoleID:    1,
		Faction:   enums.Harmony,
		Character: enums.Warrior,
	}

	if err = models.CreatePlayer(c.Request.Context(), player); err != nil {
		log.Error().Err(err).Msg("could not create player for galaxy")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create galaxy"))
		return
	}

	c.JSON(http.StatusCreated, galaxy)
}
