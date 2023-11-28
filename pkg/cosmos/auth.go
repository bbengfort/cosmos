package cosmos

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/bbengfort/cosmos/pkg/api/v1"
	"github.com/bbengfort/cosmos/pkg/auth"
	"github.com/bbengfort/cosmos/pkg/db"
	"github.com/bbengfort/cosmos/pkg/db/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) Register(c *gin.Context) {
	var (
		err error
		in  *api.RegisterRequest
		out *api.RegisterReply
	)

	in = &api.RegisterRequest{}
	if err = c.BindJSON(in); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	if err = in.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	user := &models.User{
		Name:  sql.NullString{Valid: true, String: in.Name},
		Email: in.Email,
	}

	if user.Password, err = auth.CreateDerivedKey(in.Password); err != nil {
		log.Warn().Err(err).Msg("could not create derived key for password")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete registration"))
		return
	}

	if err = models.CreateUser(c.Request.Context(), user); err != nil {
		if errors.Is(db.Check(err), db.ErrAlreadyExists) {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("user already exists"))
			return
		}

		log.Error().Err(err).Msg("could not create new user in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete registration"))
		return
	}

	out = &api.RegisterReply{
		ID:    user.ID,
		Name:  user.Name.String,
		Email: user.Email,
	}

	role, _ := user.Role(c.Request.Context())
	out.Role = role.Title

	log.Info().Int64("user_id", out.ID).Str("email", out.Email).Msg("new user registered")
	c.JSON(http.StatusCreated, out)
}

func (s *Server) Login(c *gin.Context) {
	var (
		err  error
		in   *api.LoginRequest
		out  *api.LoginReply
		user *models.User
	)

	in = &api.LoginRequest{}
	if err = c.BindJSON(in); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Fetch the user from the database to authenticate
	if user, err = models.GetUser(c.Request.Context(), in.Email); err != nil {
		if errors.Is(db.Check(err), db.ErrNotFound) {
			c.JSON(http.StatusForbidden, api.ErrorResponse("authentication failed"))
			return
		}

		log.Error().Err(err).Msg("could not fetch user from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("authentication failed"))
		return
	}

	// Authenticate the user with their password
	var verified bool
	if verified, err = auth.VerifyDerivedKey(user.Password, in.Password); err != nil {
		log.Error().Err(err).Msg("could not verify derived key")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("authentication failed"))
		return
	}

	// Wrong password
	if !verified {
		c.JSON(http.StatusForbidden, api.ErrorResponse("authentication failed"))
		return
	}

	// The user has been authenticated at this point: create access and refresh tokens
	claims := &auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: strconv.FormatInt(user.ID, 36),
		},
		Name:  user.Name.String,
		Email: user.Email,
	}

	role, _ := user.Role(c.Request.Context())
	claims.Role = role.Title

	perms, _ := user.Permissions(c.Request.Context())
	claims.Permissions = make([]string, 0, len(perms))
	for _, perm := range perms {
		claims.Permissions = append(claims.Permissions, perm.Title)
	}

	out = &api.LoginReply{}
	if out.AccessToken, out.RefreshToken, err = s.auth.CreateTokens(claims); err != nil {
		log.Error().Err(err).Msg("could not create access and refresh tokens for user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("authentication failed"))
		return
	}

	// Update the last login timestamp for user tracking
	if err = user.LoggedIn(c.Request.Context()); err != nil {
		log.Error().Err(err).Msg("could not update last login timestamp")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("authentication failed"))
		return
	}

	// Set credentials on cookies for web based applications
	auth.SetAuthCookies(c, out.AccessToken, out.RefreshToken, s.conf.Auth.CookieDomain)
	c.JSON(http.StatusOK, out)
}

func (s *Server) Logout(c *gin.Context) {
	auth.ClearAuthCookies(c, s.conf.Auth.CookieDomain)
	c.JSON(http.StatusOK, &api.Reply{Success: true})
}

func (s *Server) Reauthenticate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not implemented yet"))
}
