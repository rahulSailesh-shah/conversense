package service

import (
	"github.com/jackc/pgx/v5/pgxpool"
)


type Service struct {}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{}
}
