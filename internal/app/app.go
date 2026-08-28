package app

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sakamoto-max/diablo/internal/config"
	"github.com/sakamoto-max/diablo/internal/database"
	"github.com/sakamoto-max/diablo/internal/handlers"
	"github.com/sakamoto-max/diablo/internal/middleware"
	"github.com/sakamoto-max/diablo/internal/pkg/token"
	"github.com/sakamoto-max/diablo/internal/repository"
	"github.com/sakamoto-max/diablo/internal/router"
	"github.com/sakamoto-max/diablo/internal/services"
)

type app struct {
	httpServer *http.Server
	pool       *pgxpool.Pool
}

func NewApp(config *config.Config) *app {

	if config.Primary == "production" {
		err := database.Migrate(config)
		if err != nil {
			log.Fatalf("failed to migrate database : %v", err)
		}
	}

	pool, err := database.NewPgPool(config)
	if err != nil {
		log.Fatalf("failed to create postgres pool : %v", err)
	}

	log.Println("created pg pool")

	db := repository.NewDb(pool)

	service := services.NewService(db)

	handlers := handlers.NewHandlers(service)
	log.Println("created handlers")

	middlewares := middleware.NewMiddlewares()

	router := router.NewRouter(handlers, middlewares)
	log.Println("created router")

	token.Init(config)

	httpServer := http.Server{
		Addr:    ":" + config.Http.Port,
		Handler: router,
	}

	return &app{
		httpServer: &httpServer,
		pool:       pool,
	}
}

func (a *app) Run() {

	go a.StartHttpServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	sig := <-sigChan
	log.Printf("shutdown signal received : %v", sig)

	a.ShutDown()
}

func (a *app) StartHttpServer() {
	log.Printf("http server has started on %v", a.httpServer.Addr)
	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed to start http server : %v", err)
	}
}

func (a *app) ShutDown() {
	a.httpServer.Close()
	log.Println("http server is closed")

	a.pool.Close()
	log.Println("pg pool is closed")

	log.Println("server has shutdown")
}

//{
// "path" : "diablod/internal/app",
// "type" : "dir"
// "contents" : []byte()
//}
//
