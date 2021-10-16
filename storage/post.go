package storage

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/ogeyou/go-autth.git/model"
	"github.com/ogeyou/go-autth.git/storage/psql"
	"golang.org/x/crypto/argon2"
)
// Добавляем соль к пароль
// На время иссследований
func hashPass(plainPassword, salt string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	res := []byte(salt)
	return append(res, hashedPass...)
}
// Добавляем соль к пароль
// На время иссследований
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
// Добавляем соль к пароль
// На время иссследований
var (
	sizes       = []uint{80, 160, 320}
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func UserCreated(user model.User) int64 {
	// Соединение с экземпляром Postgres
	ctx := context.Background()
	dbpool := psql.Connect(ctx)
	defer dbpool.Close()

	var id int64

	salt := RandStringRunes(8)
	
	pass := hashPass(user.Password, salt)

	const sql = "INSERT INTO users (login, email, password) VALUES($1, $2, $3)"

	_, err := dbpool.Exec(ctx, sql, user.Login, user.Email, pass)

	if err != nil {
		fmt.Println(err)
	}
	
	fmt.Printf("Новый пользователь успешно прошёл регистрацию %v\n", id)

	return id
}

