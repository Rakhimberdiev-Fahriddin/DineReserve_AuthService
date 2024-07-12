package handler

import (
	"auth-service/storage/postgres"
	"database/sql"
	"log/slog"
)

type Handler struct {
	UserRepo *postgres.UserRepo
	Logger *slog.Logger
}

func NewHandler(db *sql.DB, logger *slog.Logger) *Handler {
	return &Handler{
		UserRepo: postgres.NewUserRepo(db),
		Logger: logger,
	}
}