package main

import (
	"flag"
	"log"
	"microservices/pkg/store/postgres"
	"microservices/services/contract/internal/delivery/http"
	"microservices/services/contract/internal/repository"
	"microservices/services/contract/internal/usecase"
	"os"
)

func main() {
	dbConnCfg := postgres.ConnConfig{}
	httpServerCfg := http.ServerConfig{}

	flag.IntVar(&httpServerCfg.Port, "http-port", 4040, "HTTP server port")
	flag.StringVar(&httpServerCfg.ReadTimeout, "http-read-timeout", "10s", "HTTP read timeout")
	flag.StringVar(&httpServerCfg.WriteTimeout, "http-write-timeout", "30s", "HTTP write timeout")
	flag.StringVar(&httpServerCfg.IdleTimeout, "http-idle-timeout", "1m", "HTTP idle timeout")

	flag.IntVar(&dbConnCfg.Port, "pg-port", 5432, "Postgres port")
	flag.StringVar(&dbConnCfg.Host, "pg-host", "localhost", "Postgres host")
	flag.StringVar(&dbConnCfg.User, "pg-user", os.Getenv("POSTGRE_USER"), "Postgres user")
	flag.StringVar(&dbConnCfg.Password, "pg-password", os.Getenv("POSTGRE_PASSWORD"), "Postgres password")
	flag.StringVar(&dbConnCfg.DbName, "pg-db-name", os.Getenv("POSTGRE_DB_NAME"), "Postgres DB name")
	flag.IntVar(&dbConnCfg.MaxOpenConns, "pg-max-open-conns", 15, "Postgres max open connections")
	flag.StringVar(&dbConnCfg.MaxIdleTime, "pg-max-idle-time", "15m", "Postgres max connection idle time")
	flag.Parse()

	db, err := postgres.OpenDB(dbConnCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Pool.Close()

	log.Print("database connection pool established")

	userRepository := repository.NewRepo(db.Pool)
	userService := usecase.New(userRepository)

	httpServer := http.NewHttpServer(http.NewRouter(userService).GetRoutes(), httpServerCfg)

	err = httpServer.Serve()
	if err != nil {
		log.Fatal("Failed to start HTTP server")
	}

}
