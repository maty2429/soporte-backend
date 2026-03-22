package routes

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	gormrepo "soporte/internal/adapters/repository/repository"
	"soporte/internal/application/services"
	"soporte/internal/config"
	"soporte/internal/delivery/http/handlers"
	docs "soporte/internal/delivery/http/handlers/docs"
	"soporte/internal/delivery/http/middlewares"
)

func NewRouter(log *slog.Logger, cfg config.Config, db *gorm.DB, startedAt time.Time) *gin.Engine {
	if config.IsProduction || cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Silencia el banner "[GIN-debug] [WARNING] Running in debug mode" y los prints de rutas.
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	router := gin.New()

	// Restaura stdout para que gin.Logger() muestre los requests con colores.
	gin.DefaultWriter = os.Stdout

	trustedProxies := cfg.Security.TrustedProxies
	if len(trustedProxies) == 0 {
		trustedProxies = []string{}
	}
	if err := router.SetTrustedProxies(trustedProxies); err != nil {
		panic("invalid trusted proxies: " + err.Error())
	}

	router.Use(
		middlewares.RequestID(),
		middlewares.Recovery(log),
		gzip.Gzip(gzip.DefaultCompression),
		middlewares.CORS(cfg.Security),
		middlewares.SecurityHeaders(),
		middlewares.RequestSizeLimit(cfg.Security.RequestSizeLimitByte),
		middlewares.ContentTypeJSON(),
		middlewares.RateLimit(cfg.Security),
		middlewares.QueryTimeout(cfg.Database.QueryTimeout),
		middlewares.Metrics(),
		gin.Logger(),
	)

	healthHandler := handlers.NewHealthHandler(cfg, db, startedAt)
	docsHandler := docs.NewDocsHandler()

	solicitanteRepo := gormrepo.NewSolicitanteRepository(db)
	catalogoRepo := gormrepo.NewCatalogoRepository(db)

	tecnicoRepo := gormrepo.NewTecnicoRepository(db)

	solicitanteHandler := handlers.NewSolicitanteHandler(services.NewSolicitanteService(solicitanteRepo))
	tecnicoHandler := handlers.NewTecnicoHandler(services.NewTecnicoService(tecnicoRepo))
	catalogoHandler := handlers.NewCatalogoHandler(services.NewCatalogoService(catalogoRepo))
	servicioHandler := handlers.NewServicioHandler(services.NewServicioService(gormrepo.NewServicioRepository(db)))
	catalogoFallaHandler := handlers.NewCatalogoFallaHandler(services.NewCatalogoFallaService(gormrepo.NewCatalogoFallaRepository(db)))
	ticketHandler := handlers.NewTicketHandler(services.NewTicketService(gormrepo.NewTicketRepository(db), solicitanteRepo, catalogoRepo))

	v1 := router.Group("/api/v1")

	registerInfraRoutes(router, v1, healthHandler, docsHandler, cfg)
	registerSolicitanteRoutes(v1, solicitanteHandler)
	registerTecnicoRoutes(v1, tecnicoHandler)
	registerTicketRoutes(v1, ticketHandler)
	registerServicioRoutes(v1, servicioHandler)
	registerCatalogoFallaRoutes(v1, catalogoFallaHandler)
	registerCatalogoRoutes(v1, catalogoHandler)

	registerRootRoute(router, cfg)

	return router
}
