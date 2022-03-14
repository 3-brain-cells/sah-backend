package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/3-brain-cells/sah-backend/db/mongo"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// APIServer is a struct that bundles together the various server-wide
// resources used at runtime that each have
// a lifecycle of initialization, connection, and disconnection
type APIServer struct {
	dbProvider *mongo.Provider
	logger     zerolog.Logger
}

// NewAPIServer initializes the struct and all constituent components
func NewAPIServer(logger zerolog.Logger) (*APIServer, error) {

	// Initialize the MongoDB handler
	dbProvider, err := mongo.NewProvider(logger)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize MongoDB handler")
	}

	return &APIServer{
		dbProvider: dbProvider,
		logger:     logger,
	}, nil
}

// Connect initializes the struct and all constituent components
func (a *APIServer) Connect(ctx context.Context) error {
	// Connect to the MongoDB database
	a.logger.Info().Msg("initializing MongoDB database provider")
	err := a.dbProvider.Connect(ctx)
	if err != nil {
		return errors.Wrap(err, "could not disconnect to the database")
	}
	a.logger.Info().Msg("successfully connected to and pinged the database")

	return nil
}

// Disconnect initializes the struct and all constituent components
func (a *APIServer) Disconnect(ctx context.Context) error {
	err := a.dbProvider.Disconnect(ctx)
	if err != nil {
		return errors.Wrap(err, "could not disconnect from the database")
	}
	a.logger.Info().Msg("disconnected from the database")

	return nil
}

// Serve runs the main API server until it's cancelled for some reason,
// in which case it attempts to gracefully shutdown.
// This function blocks.
func (a *APIServer) Serve(ctx context.Context, port int) {
	router := a.routes()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal().Err(err).Msg("an error occurred when serving the HTTP server")
		}
	}()
	a.logger.Info().Int("port", port).Msg("API server started")

	<-ctx.Done()
	a.logger.Info().Msg("API server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		a.logger.Fatal().Err(err).Msg("API server shutdown failed")
	}
	a.logger.Info().Msg("API server exited properly")
}

func (a *APIServer) routes() *chi.Mux {
	// Approach from:
	// https://itnext.io/structuring-a-production-grade-rest-api-in-golang-c0229b3feedc
	// https://itnext.io/how-i-pass-around-shared-resources-databases-configuration-etc-within-golang-projects-b27af4d8e8a
	router := chi.NewRouter()
	router.Use(
		middleware.Recoverer,      // Recover from panics without crashing the server
		hlog.NewHandler(a.logger), // Attach a logger to each request
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			// Log API request calls once they complete:
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status", status).
				Int("bytes_out", size).
				Dur("duration", duration).
				Str("ip", r.RemoteAddr).
				Str("user_agent", r.Header.Get("User-Agent")).
				Msg("handled HTTP request")
		}),
		hlog.RequestIDHandler("req_id", "X-Request-Id"), // Attach a unique request ID to each incoming request
		middleware.RedirectSlashes,                      // Redirect slashes to no slash URL versions
		render.SetContentType(render.ContentTypeJSON),   // Set content-type headers to application/json
		middleware.Compress(5),                          // Compress results, mostly gzipping assets and json
		middleware.NoCache,                              // Prevent clients from caching the results
		a.corsMiddleware(),                              // Create cors middleware from go-chi/cors
	)

	// ==============================
	// Add all routes to the API here
	// ==============================
	router.Route("/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			// Can be used for health checks
			r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(204)
			})

			// r.Mount("/auth", apiAuth.Routes(a.casProvider, a.dbProvider, a.jwtManager))
		})

		// Protected routes
		// r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens,
		// sending appropriate status codes upon failure.
		// Note that this does not perform *authorization* checks involving perms;
		// if needed, use auth.AdminAuthenticator to use Permissions.AdminAccess

		// r.Mount("/announcements", announcements.Routes(a.dbProvider))
		// })
	})

	return router
}

func (a *APIServer) corsMiddleware() func(http.Handler) http.Handler {
	// See if the CORS_ALLOWED_ORIGINS environment variable was set
	allowedOrigins := "*"
	if value, ok := os.LookupEnv("CORS_ALLOWED_ORIGINS"); ok {
		allowedOrigins = value
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{allowedOrigins},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
		MaxAge:           300,
	})
}
