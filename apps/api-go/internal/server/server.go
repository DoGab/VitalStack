package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	DefaultReadHeaderTimeout = 30 * time.Second
	// DefaultMaxMultipartMemory defines the maximum size of a multipart form request body (10MB)
	DefaultMaxMultipartMemory = 10 << 20
)

// Controller is an interface for controllers
type Controller interface {
	Register(api huma.API)
}

type (
	// Option is a function that configures the server
	Option func(*Server)
	// ShutdownFunc is a function that shuts down the server
	ShutdownFunc func(context.Context) error
)

// Server is a server type. It is used to setup a new api server,
// register controllers and start the server.
type Server struct {
	name    string
	version string

	srv    *http.Server
	router *gin.Engine

	huma huma.API
}

// NewServer creates a new server instance
func NewServer(addr string, opts ...Option) (*Server, ShutdownFunc) {
	router := gin.New()
	router.Use(gin.Recovery())

	// Allow large request bodies for image uploads (10MB)
	router.MaxMultipartMemory = DefaultMaxMultipartMemory

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
			"https://localhost:5173",
			"https://127.0.0.1:5173",
		},
		AllowOriginFunc: func(origin string) bool {
			// Allow any local network origin for mobile testing
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	s := &Server{
		srv: &http.Server{
			Addr:              addr,
			Handler:           router,
			ReadHeaderTimeout: DefaultReadHeaderTimeout,
		},
		router:  router,
		name:    "MacroGuard API",
		version: "1.0.0",
	}

	s.huma = newHumaAPI(s)

	for _, opt := range opts {
		opt(s)
	}

	s.registerDiagnosticEndpoints()

	shutdownFunc := func(ctx context.Context) error {
		slog.Info("MacroGuard API shutting down...")
		if err := s.srv.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	}

	return s, shutdownFunc
}

func newHumaAPI(s *Server) huma.API {
	config := huma.DefaultConfig(s.name, s.version)
	config.SchemasPath = "/schemas"
	config.Servers = []*huma.Server{
		{
			URL: "/",
		},
	}
	return humagin.New(s.router, config)
}

// RegisterAPI registers the API with all provided controllers
func (s *Server) RegisterAPI(controllers ...Controller) {
	for _, ctrl := range controllers {
		ctrl.Register(s.huma)
	}
}

// Serve starts the server
func (s *Server) Serve(ctx context.Context) error {
	go func() {
		err := s.srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	return nil
}

func (s *Server) registerDiagnosticEndpoints() {
	s.router.GET("/api/health", s.healthCheckHandler)
}

func (s *Server) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
