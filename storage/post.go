package storage

import (
	"context"
	"fmt"

	"log"
	"github.com/ogeyou/go-autth.git/storage/psql"
	"github.com/ogeyou/go-autth.git/model"
)

func UserCreated(user model.User) int64 {
	// Соединение с экземпляром Postgres
	ctx := context.Background()
	dbpool := psql.Connect(ctx)
	defer dbpool.Close()

	var id int64

	_,err := dbpool.Exec(ctx, "INSERT INTO users (login, email, password) VALUES($1, $2, $3)", user.Login, user.Email, user.Password)

	if err != nil {
		log.Fatalf("Ошибка при добавлении нового пользователя. %v", err)
	}
	fmt.Printf("Новый пользователь успешно прошёл регистрацию %v\n", id)

	return id
}
