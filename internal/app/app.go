package app

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/sakamoto-max/diablo/internal/config"
	"github.com/sakamoto-max/diablo/internal/database"
	"github.com/sakamoto-max/diablo/internal/handlers"
	"github.com/sakamoto-max/diablo/internal/repository"
	"github.com/sakamoto-max/diablo/internal/router"
	"github.com/sakamoto-max/diablo/internal/services"
)

type app struct {
	httpServer *http.Server
	dbConn     *sql.DB
	dbDown     string
}

func NewApp(config *config.Config) *app {

	dbConn, err := database.New()
	if err != nil {
		log.Fatalf("failed to create database connection : %v", err)
	}

	log.Println("created db Conn")

	if config.Db.Up == "yes" {
		err := database.CreateTables(dbConn)
		if err != nil {
			log.Fatalf("failed to create tables : %v", err)
		}
		log.Println("created the tables")
	}


	repo := repository.New(dbConn)

	service := services.NewService(repo)

	handlers := handlers.NewHandlers(service)
	log.Println("created handlers")

	router := router.NewRouter(handlers)

	httpServer := http.Server{
		Addr:    ":" + config.Http.Port,
		Handler: router,
	}

	return &app{
		httpServer: &httpServer,
		dbConn:     dbConn,
		dbDown:     config.Db.Down,
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

	if a.dbDown == "yes" {
		database.DropTables(a.dbConn)
		log.Println("db tables are dropped")
	}

	a.dbConn.Close()
	log.Println("db is closed")

	log.Println("server has shutdown")
}
