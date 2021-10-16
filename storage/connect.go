package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)


func Connect(ctx context.Context) *pgxpool.Pool {

	url := "postgres://postgres:password@localhost:5432/university?sslmode=disable"

	dbpool, err := pgxpool.Connect(ctx, url)

	if err != nil {
		log.Println("Ошибка подключения к базе данных Postgresql")
	}

	return dbpool
}