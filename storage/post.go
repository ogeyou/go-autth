package database

import (
	"context"
	"fmt"

	"log"
)

func UserCreated(user model.User) int64 {
	// Соединение с экземпляром Postgres
	ctx := context.Background()
	dbpool := database.Connect(ctx)
	defer dbpool.Close()

	var id int64

	err := dbpool.QueryRow(ctx, `INSERT INTO users (login, email, password) VALUES($1, $2, $3) returning login;`, user.Login, user.Email, user.Password).Scan(&id)

	if err != nil {
		log.Fatalf("Ошибка при добавлении нового пользователя. %v", err)
	}
	fmt.Printf("Новый пользователь успешно прошёл регистрацию %v\n", id)

	return id
}
