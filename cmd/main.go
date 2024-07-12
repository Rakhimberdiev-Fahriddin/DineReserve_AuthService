package main

import (
	"auth-service/api"
	"auth-service/api/handler"
	"auth-service/config"
	"auth-service/logs"
	"auth-service/server"
	"auth-service/storage/postgres"
	"log"
	"sync"
)

func main() {
	logs.InitLogger()
	logs.Logger.Info("Starting the server")
	db, err := postgres.ConnectDB()
	if err != nil {
		logs.Logger.Error("Faild connect to postgres", "error", err.Error())
		log.Fatal(err)
	}
	defer db.Close()

	cfg := config.Load()
	router := api.Routes(handler.NewHandler(db, logs.Logger))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		server.ServerRun(db)
	}()

	go func() {
		defer wg.Done()
		logs.Logger.Info("server is running", "PORT", cfg.HTTP_PORT)
		err := router.Run(cfg.HTTP_PORT)
		if err != nil {
			logs.Logger.Error("Faild server is running", "error", err.Error())
			log.Fatal(err)
		}
	}()

	wg.Wait()
}
