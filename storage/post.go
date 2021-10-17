package storage

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"

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

	salt := RandStringRunes(8)
	var id int64
	pass := hashPass(user.Password, salt)

	const sql = "INSERT INTO users (login, email, password) VALUES($1, $2, $3)"

	result, err := dbpool.Exec(ctx, sql, user.Login, user.Email, pass)
	var w http.ResponseWriter
	if err != nil {
		log.Println("Ощибка при добавлении нового пользователя в базу данных", err)
		http.Error(w, "Ощибка на сервере при добавлении нового пользователя в базу данных", http.StatusInternalServerError)
		return 0
	}

	affected := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "Looks like user exists", http.StatusBadRequest)
		return 1
	}

	fmt.Printf("%T\n", result)
	fmt.Printf("Новый пользователь успешно прошёл регистрацию %v\n", id)

	return id
}
