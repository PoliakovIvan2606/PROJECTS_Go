package app

import (
	"context"
	"errors"
	"fmt"
	"internet_shop/internal/config"
	"internet_shop/internal/domain/product/storage"
	"internet_shop/pkg/client/postgresql"
	"internet_shop/pkg/logging"
	"internet_shop/pkg/metric"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	_ "internet_shop/docs"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	logger *logging.Logger
	cfg *config.Config
	router *httprouter.Router
	httpServer *http.Server
	pgClient *pgxpool.Pool
}

func NewApp(logger *logging.Logger, config *config.Config) (App, error) {
	// Logger
	logger.Println("router initializing")
	router := httprouter.New()

	// Swagger
	logger.Println("swagger docs initializing")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	// Router
	logger.Info("register heartbeat init")
	metricHeandler := metric.Handler{}
	metricHeandler.Register(router)

	// PostgresDB
	pgConfig := postgresql.NewPgConfig(
		config.PostgreSQL.Username, config.PostgreSQL.Password,
		config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.Database,
	)
	pgClient, err := postgresql.NewClient(context.Background(), 5, time.Second*5, pgConfig)
	if err != nil {
		logger.Fatal(err)
	}

	// storage
	productStorage := storage.NewProductStorage(pgClient, logger)
	all, err := productStorage.All(context.Background())
	if err != nil {
		logger.Fatal(err)
	}
	logger.Fatal(all)
	

	return App{
		logger: logger,
		cfg: config,
		router: router,
		pgClient: pgClient,
	}, nil
}

func (a *App) Run() {
	a.startHTTP()
}


func (a *App) startHTTP() {
	a.logger.Info("start HTTP")
	
	var listener net.Listener

	if a.cfg.Listen.Type == config.LISTEN_TYPE_SOCK {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			a.logger.Fatal(err)
		}
		socketPath := path.Join(appDir, a.cfg.Listen.SocketFile)
		a.logger.Infof("socket path: %s", socketPath)

		a.logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			a.logger.Fatal(err)
		}
	} else {
		a.logger.Infof("bind application to host: %s and port: %s", a.cfg.Listen.BindIP, a.cfg.Listen.Port)
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
		if err != nil {
			a.logger.Fatal(err)
		}
	}

	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodOptions, http.MethodDelete},
		AllowedOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowCredentials:   true,
		AllowedHeaders:   []string{
			"Authorization",
			"Location",
			"Content-Type",
			"Origin",
			"Accept",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
		},
		OptionsPassthrough: false,
		ExposedHeaders:     []string{
			"Location",
			"Authorization",
			"Content-Disposition",
		},
		Debug: true,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler: handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	a.logger.Println("application completely inilized and started")

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warn("server shutdoun")
		default:
			a.logger.Fatal(err)
		}
	}

	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		a.logger.Fatal(err)
	}
}