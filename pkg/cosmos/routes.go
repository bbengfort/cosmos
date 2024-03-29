package cosmos

import (
	"time"

	"github.com/bbengfort/cosmos/pkg/auth"
	"github.com/bbengfort/cosmos/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup the server's middleware and routes.
func (s *Server) setupRoutes() (err error) {
	// Setup CORS configuration
	corsConf := cors.Config{
		AllowMethods:     []string{"GET", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-TOKEN"},
		AllowOrigins:     s.conf.AllowOrigins,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Application Middleware
	// NOTE: ordering is important to how middleware is handled
	middlewares := []gin.HandlerFunc{
		// Logging should be on the outside so we can record the correct latency of requests
		// NOTE: logging panics will not recover
		logger.GinLogger("cosmos"),

		// Panic recovery middleware
		gin.Recovery(),

		// CORS configuration allows the front-end to make cross-origin requests
		cors.New(corsConf),

		// Mainenance mode handling
		s.Available(),
	}

	// Add the middleware to the router
	for _, middleware := range middlewares {
		if middleware != nil {
			s.router.Use(middleware)
		}
	}

	// Create authentication middleware to add to specific routes
	authenticate := auth.Authenticate(s.auth)

	// Kubernetes liveness probes
	s.router.GET("/healthz", s.Healthz)
	s.router.GET("/livez", s.Healthz)
	s.router.GET("/readyz", s.Readyz)

	// NotFound and NotAllowed routes
	s.router.NoRoute(s.NotFound)
	s.router.NoMethod(s.NotAllowed)

	// Add the v1 API routes
	v1 := s.router.Group("/v1")
	{
		// Heartbeat route
		v1.GET("/status", s.Status)

		// Authentication routes
		v1.POST("/register", s.Register)
		v1.POST("/login", s.Login)
		v1.POST("/logout", s.Logout)
		v1.POST("/reauthenticate", s.Reauthenticate)

		// Galaxy resource
		galaxy := v1.Group("/galaxy", authenticate)
		{
			galaxy.GET("/", s.ListGalaxies, auth.Authorize("games:read"))
			galaxy.POST("/", s.CreateGalaxy, auth.Authorize("games:create"))
		}
	}

	return nil
}
