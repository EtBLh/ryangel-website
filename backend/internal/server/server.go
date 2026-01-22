package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/ryangel/ryangel-backend/internal/config"
	"github.com/ryangel/ryangel-backend/internal/http/handlers"
	"github.com/ryangel/ryangel-backend/internal/repository"
	authsvc "github.com/ryangel/ryangel-backend/internal/services/auth"
	ebuysvc "github.com/ryangel/ryangel-backend/internal/services"
)

// Options configures the HTTP server bootstrap.
type Options struct {
	Config      *config.Config
	DB          *pgxpool.Pool
	Logger      *zap.Logger
	AuthService *authsvc.Service
	EbuyService *ebuysvc.EbuyService
}

// Server wraps the Gin engine and http.Server.
type Server struct {
	engine *gin.Engine
	http   *http.Server
}

// New constructs a Server with registered routes.
func New(opts Options) *Server {
	if strings.EqualFold(opts.Config.AppEnv, "production") {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// CORS middleware - allow all origins
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Cart-ID"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	if opts.Logger != nil {
		router.Use(ginzap.GinzapWithConfig(opts.Logger, &ginzap.Config{TimeFormat: time.RFC3339, UTC: true}))
		router.Use(ginzap.RecoveryWithZap(opts.Logger, true))
	} else {
		router.Use(gin.Logger(), gin.Recovery())
	}

	api := router.Group("/api")
	// Static media serving is handled by Nginx
	// api.Static("/media", "./media")

	healthHandler := handlers.HealthHandler{DB: opts.DB}
	healthHandler.Register(api)

	productRepo := repository.NewProductRepository(opts.DB)
	productHandler := handlers.ProductHandler{Repo: productRepo, Config: opts.Config}
	productHandler.Register(api)

	ebuyStoreRepo := repository.NewEbuyStoreRepository(opts.DB)
	ebuyStoreHandler := handlers.EbuyStoreHandler{Repo: ebuyStoreRepo}
	ebuyStoreHandler.Register(api)

	cartRepo := repository.NewCartRepository(opts.DB)
    discountRepo := repository.NewDiscountRepository(opts.DB)
	cartHandler := handlers.CartHandler{
        Repo: cartRepo, 
        EbuyStoreRepo: ebuyStoreRepo,
        DiscountRepo: discountRepo,
    }
	cartHandler.Register(api, opts.AuthService)

	orderRepo := repository.NewOrderRepository(opts.DB)
	orderHandler := handlers.OrderHandler{Orders: orderRepo}
	orderHandler.Register(api, opts.AuthService)
	orderHandler.RegisterAdmin(api, opts.AuthService)

	if opts.AuthService != nil {
		authHandler := handlers.AuthHandler{Service: opts.AuthService, Config: opts.Config}
		authHandler.RegisterAdminRoutes(api)
		authHandler.RegisterClientRoutes(api)
		authHandler.RegisterGoogleRoutes(api)
	}

	if opts.EbuyService != nil {
		ctx := context.Background()
		opts.EbuyService.StartScheduler(ctx)
	}

	httpSrv := &http.Server{
		Addr:              opts.Config.HTTPAddr(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	return &Server{engine: router, http: httpSrv}
}

// Run starts the HTTP server.
func (s *Server) Run() error {
	return s.http.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
