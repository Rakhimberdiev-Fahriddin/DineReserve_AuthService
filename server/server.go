package server

import (
	"auth-service/config"
	pb "auth-service/generated/auth_service"
	"auth-service/logs"
	"auth-service/service"
	"auth-service/storage/postgres"
	"database/sql"
	"log"
	"net"

	"google.golang.org/grpc"
)

func ServerRun(db *sql.DB) {
	logs.InitLogger()
	cfg := config.Load()
	listener, err := net.Listen("tcp", cfg.GRPC_PORT)
	if err != nil {
		logs.Logger.Error("Error create to new listener", "error", err.Error())
		log.Fatal(err)
	}

	s := grpc.NewServer()
	service := service.NewAuthService(postgres.NewUserRepo(db))

	pb.RegisterAuthServiceServer(s, service)

	logs.Logger.Info("server is running ", "PORT", cfg.GRPC_PORT)
	log.Printf("server is running on %v...", listener.Addr())
	if err := s.Serve(listener); err != nil {
		logs.Logger.Error("Faild server is running", "error", err.Error())
		log.Fatal(err)
	}
}
