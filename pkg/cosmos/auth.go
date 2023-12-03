package cosmos

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bbengfort/cosmos/pkg/api/v1"
	"github.com/bbengfort/cosmos/pkg/auth"
	"github.com/bbengfort/cosmos/pkg/db"
	"github.com/bbengfort/cosmos/pkg/db/models"
	"github.com/gin-gonic/gin"
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
		err    error
		in     *api.LoginRequest
		out    *api.LoginReply
		user   *models.User
		claims *auth.Claims
	)

	in = &api.LoginRequest{}
	if err = c.BindJSON(in); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Fetch the user from the database to authenticate (username is the user's email)
	if user, err = models.GetUser(c.Request.Context(), in.Username); err != nil {
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
	if claims, err = auth.NewClaimsForUser(c.Request.Context(), user); err != nil {
		log.Error().Err(err).Msg("could not create claims for user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("authentication failed"))
		return
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
	var (
		err           error
		in            *api.ReauthenticateRequest
		out           *api.LoginReply
		userID        int64
		user          *models.User
		accessToken   string
		refreshClaims *auth.Claims
		accessClaims  *auth.Claims
		claims        *auth.Claims
	)

	// Attempt to bind the request from the user
	in = &api.ReauthenticateRequest{}
	if err = c.BindJSON(in); err != nil || in.RefreshToken == "" {
		// If we couldn't get the refresh token from the request, attempt to get it
		// from the cookies in the header of the request.
		if in.RefreshToken, err = auth.GetRefreshToken(c); err != nil || in.RefreshToken == "" {
			log.Debug().Err(err).Msg("could not get refresh token from request")
			c.JSON(http.StatusBadRequest, api.ErrorResponse("no reauthentication credentials"))
			return
		}
	}

	// Check to ensure the refresh token is still valid.
	// NOTE: this will also validate the not before and not after claims
	if refreshClaims, err = s.auth.Verify(in.RefreshToken); err != nil {
		log.Debug().Err(err).Msg("invalid refresh token")
		c.JSON(http.StatusForbidden, api.ErrorResponse("reauthentication failed"))
		return
	}

	// Fetch the access token from the request
	if accessToken, err = auth.GetAccessToken(c); err != nil || accessToken == "" {
		log.Debug().Err(err).Msg("no access token in reauthenticate request")
		c.JSON(http.StatusForbidden, api.ErrorResponse("reauthentication failed"))
		return
	}

	// Get the access token claims
	if accessClaims, err = s.auth.Parse(accessToken); err != nil {
		log.Debug().Err(err).Msg("invalid access token")
		c.JSON(http.StatusForbidden, api.ErrorResponse("reauthentication failed"))
		return
	}

	// Ensure the access and refresh token match
	if accessClaims.ID != refreshClaims.ID || accessClaims.Subject != refreshClaims.Subject {
		log.Debug().Msg("access token claims do not match refresh token claims")
		c.JSON(http.StatusForbidden, api.ErrorResponse("reauthentication failed"))
		return
	}

	if userID, err = refreshClaims.SubjectID(); err != nil {
		log.Error().Err(err).Str("subject", refreshClaims.Subject).Msg("could not parse user ID from refresh claims")
		c.JSON(http.StatusForbidden, api.ErrorResponse("reauthentication failed"))
		return
	}

	// Fetch the user to get the most up to date claims (do not rely on old claims)
	if user, err = models.GetUser(c.Request.Context(), userID); err != nil {
		log.Error().Err(err).Msg("could not create access and refresh tokens for user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("reauthentication failed"))
		return
	}

	// The user has been reauthenticated at this point: create access and refresh tokens
	if claims, err = auth.NewClaimsForUser(c.Request.Context(), user); err != nil {
		log.Error().Err(err).Msg("could not create claims for user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("reauthentication failed"))
		return
	}

	out = &api.LoginReply{}
	if out.AccessToken, out.RefreshToken, err = s.auth.CreateTokens(claims); err != nil {
		log.Error().Err(err).Msg("could not create access and refresh tokens for user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("reauthentication failed"))
		return
	}

	// Update the last login timestamp for user tracking
	if err = user.LoggedIn(c.Request.Context()); err != nil {
		log.Error().Err(err).Msg("could not update last login timestamp")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("reauthentication failed"))
		return
	}

	// Return the new access and refresh tokens and set cookies as needed
	auth.SetAuthCookies(c, out.AccessToken, out.RefreshToken, s.conf.Auth.CookieDomain)
	c.JSON(http.StatusOK, out)
}
