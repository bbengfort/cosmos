package auth

import (
	"net/http"
	"regexp"
	"time"

	api "github.com/bbengfort/cosmos/pkg/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	authorization      = "Authorization"
	ContextUserClaims  = "user_claims"
	ContextAccessToken = "access_token"
	ContextRequestID   = "request_id"
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
)

// used to extract the access token from the header
var (
	bearer = regexp.MustCompile(`^\s*[Bb]earer\s+([a-zA-Z0-9_\-\.]+)\s*$`)
)

func Authenticate(issuer *ClaimsIssuer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err         error
			accessToken string
			claims      *Claims
		)

		// Fetch access token from the request, if no access token is available, reject.
		if accessToken, err = GetAccessToken(c); err != nil {
			log.Debug().Err(err).Msg("no access token in authenticated request")
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(ErrAuthRequired))
			return
		}

		if claims, err = issuer.Verify(accessToken); err != nil {
			log.Warn().Err(err).Msg("invalid access token in request")
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(ErrAuthRequired))
			return
		}

		// Add claims to context fo ruse in downstream processing
		c.Set(ContextUserClaims, claims)
		c.Next()
	}
}

func Authorize(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := GetClaims(c)
		if err != nil {
			log.Warn().Err(err).Msg("no claims in request")
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(ErrNotAuthorized))
			return
		}

		if !claims.HasAllPermissions(permissions...) {
			log.Warn().Err(err).Msg("user does not have required permissions")
			c.AbortWithStatusJSON(http.StatusForbidden, api.ErrorResponse(ErrNotAuthorized))
			return
		}

		c.Next()
	}
}

func Reauthenticate(issuer *ClaimsIssuer) gin.HandlerFunc {
	reauthenticate := func(c *gin.Context) (err error) {
		// Get the refresh token from the cookies or the headers of the request.
		var refreshToken string
		if refreshToken, err = GetRefreshToken(c); err != nil {
			// If there is no refresh token, return no error.
			return nil
		}

		// Check to ensure the refresh token is still valid.
		// NOTE: this will also validate the not before and not after claims
		if _, err = issuer.Verify(refreshToken); err != nil {
			return err
		}

		// Get the access tokens to fetch the claims to refresh auth
		var oldAccessToken string
		if oldAccessToken, err = GetAccessToken(c); err != nil {
			return err
		}

		// Parse the claims from the old access token
		var claims *Claims
		if claims, err = issuer.Parse(oldAccessToken); err != nil {
			return nil
		}

		// Create new access and refresh tokens
		var accessToken, newRefreshToken string
		if accessToken, newRefreshToken, err = issuer.CreateTokens(claims); err != nil {
			return err
		}

		// Set the new access and refresh cookies
		if err = SetAuthCookies(c, accessToken, newRefreshToken, issuer.conf.CookieDomain); err != nil {
			return err
		}

		return nil
	}

	return func(c *gin.Context) {
		if err := reauthenticate(c); err != nil {
			log.Debug().Err(err).Msg("could not reauthenticate request")
		}
		c.Next()
	}
}

// GetAccessToken retrieves the bearer token from the authorization header and parses it
// to return only the JWT access token component of the header. Alternatively, if the
// authorization header is not present, then the token is fetched from cookies. If the
// header is missing or the token is not available, an error is returned.
//
// NOTE: the authorization header takes precedence over access tokens in cookies.
func GetAccessToken(c *gin.Context) (tks string, err error) {
	// Attempt to get the access token from the header.
	if header := c.GetHeader(authorization); header != "" {
		match := bearer.FindStringSubmatch(header)
		if len(match) == 2 {
			return match[1], nil
		}
		return "", ErrParseBearer
	}

	// Attempt to get the access token from cookies.
	var cookie string
	if cookie, err = c.Cookie(AccessTokenCookie); err == nil {
		// If the error is nil, that means we were able to retrieve the access token cookie
		return cookie, nil
	}
	return "", ErrNoAuthorization
}

// GetRefreshToken retrieves the refresh token from the cookies in the request. If the
// cookie is not present or expired then an error is returned.
func GetRefreshToken(c *gin.Context) (tks string, err error) {
	if tks, err = c.Cookie(RefreshTokenCookie); err != nil {
		return "", ErrNoRefreshToken
	}
	return tks, nil
}

func GetClaims(c *gin.Context) (*Claims, error) {
	claims, exists := c.Get(ContextUserClaims)
	if !exists {
		return nil, ErrNoClaims
	}
	return claims.(*Claims), nil
}

// SetAuthCookies is a helper function to set authentication cookies on a gin request.
// The access token cookie (access_token) is an http only cookie that expires when the
// access token expires. The refresh token cookie is not an http only cookie (it can be
// accessed by client-side scripts) and it expires when the refresh token expires. Both
// cookies require https and will not be set (silently) over http connections.
func SetAuthCookies(c *gin.Context, accessToken, refreshToken, domain string) (err error) {
	// Parse access token to get expiration time
	var accessExpires time.Time
	if accessExpires, err = ExpiresAt(accessToken); err != nil {
		return err
	}

	// Set the access token cookie: httpOnly is true; cannot be accessed by Javascript
	accessMaxAge := int((time.Until(accessExpires.Add(600 * time.Second))).Seconds())
	c.SetCookie(AccessTokenCookie, accessToken, accessMaxAge, "/", domain, true, true)

	// Parse refresh token to get expiration time
	var refreshExpires time.Time
	if refreshExpires, err = ExpiresAt(refreshToken); err != nil {
		return err
	}

	// Set the refresh token cookie: httpOnly is false; can be accessed by Javascript
	refreshMaxAge := int((time.Until(refreshExpires.Add(600 * time.Second))).Seconds())
	c.SetCookie(RefreshTokenCookie, refreshToken, refreshMaxAge, "/", domain, true, false)
	return nil
}

// ClearAuthCookies is a helper function to clear authentication cookies on a gin
// request to effectively log out a user.
func ClearAuthCookies(c *gin.Context, domain string) {
	c.SetCookie(AccessTokenCookie, "", -1, "/", domain, true, true)
	c.SetCookie(RefreshTokenCookie, "", -1, "/", domain, true, false)
}
