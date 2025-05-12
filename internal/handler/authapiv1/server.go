package authapiv1

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/render"
	"github.com/novychok/authasvs/internal/service"
	authapiv1 "github.com/novychok/authasvs/pkg/authApi/v1"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	quietDownPeriod   = 10
	corsMaxAge        = 300
	readHeaderTimeout = 10
)

type Server struct {
	l   *slog.Logger
	cfg *Config
	h   authapiv1.ServerInterface

	authService service.Auth
}

type Config struct {
	Port              int    `mapstructure:"AUTH_API_V1_PORT"`
	ReadHeaderTimeout int    `mapstructure:"READ_HEADER_TIMEOUT"`
	QuietDownPeriod   int    `mapstructure:"QUIET_DOWN_PERIOD"`
	CorsMaxAge        int    `mapstructure:"CORS_MAX_AGE"`
	BaseDomain        string `mapstructure:"BASE_DOMAIN"`
}

func (s *Server) Start(ctx context.Context) error {
	logger := httplog.NewLogger("platform-api", httplog.Options{
		JSON:            true,
		LogLevel:        slog.LevelInfo,
		Concise:         true,
		RequestHeaders:  true,
		QuietDownRoutes: []string{"/health"},
		QuietDownPeriod: time.Duration(quietDownPeriod) * time.Second,
	})

	router := chi.NewRouter()
	router.Use(
		httplog.RequestLogger(logger),
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
		render.SetContentType(render.ContentTypeJSON),
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           corsMaxAge,
		}),
		middleware.Heartbeat("/health"),
		ContextMiddleware(),
	)

	swagger, err := authapiv1.GetSwagger()
	if err != nil {
		s.l.Error("Failed to load Swagger spec", slog.String("error", err.Error()))
		return fmt.Errorf("failed to load Swagger spec: %w", err)
	}

	swagger.Servers = nil

	oapiOpts := oapimiddleware.Options{
		SilenceServersWarning: false,
		ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			_ = json.NewEncoder(w).Encode(authapiv1.Error{Message: message})
		},
		Options: openapi3filter.Options{
			AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
				req := input.RequestValidationInput.Request
				customCtx := ContextFromRequest(req)
				switch input.SecuritySchemeName {
				case "bearerAuth":
					customCtx.Set("authRequired", true)
				case "keyAuth":
					customCtx.Set("apiKeyRequired", true)
				}
				return nil
			},
			ExcludeRequestBody:    true,
			ExcludeResponseBody:   false,
			IncludeResponseStatus: true,
			MultiError:            false,
		},
	}

	router.Use(oapimiddleware.OapiRequestValidatorWithOptions(swagger, &oapiOpts))

	router.Use(s.auth)

	authapiv1.HandlerFromMux(s.h, router)

	router.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		spec, err := authapiv1.GetSwagger()
		if err != nil {
			http.Error(w, "Could not get Swagger spec", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(spec)
	})

	router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger.json")))

	// err = chi.Walk(router, func(method string, route string, handler http.Handler,
	// 	middlewares ...func(http.Handler) http.Handler) error {
	// 	fmt.Printf("%s %s\n", method, route)
	// 	return nil
	// })
	// if err != nil {
	// 	fmt.Printf("Failed to walk routes: %s\n", err.Error())
	// }

	srv := &http.Server{
		Handler: otelhttp.NewHandler(router, "billingstix-platform-api-v1"),
		Addr:    fmt.Sprintf(":%d", s.cfg.Port),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		ReadHeaderTimeout: time.Duration(readHeaderTimeout) * time.Second,
	}

	s.l.Info("Starting Platform API server", slog.String("address", srv.Addr))
	s.l.Info("Swagger UI available", slog.String("url", fmt.Sprintf("http://localhost%s/swagger/index.html", srv.Addr)))

	return srv.ListenAndServe()
}

func NewServer(
	l *slog.Logger,
	cfg *Config,
	h authapiv1.ServerInterface,

	authService service.Auth,
) *Server {
	return &Server{
		l:   l,
		cfg: cfg,
		h:   h,

		authService: authService,
	}
}
