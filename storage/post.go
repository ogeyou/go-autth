package storage

import (
	"bytes"
	"context"
	"fmt"

	"log"
	"math/rand"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/ogeyou/go-autth.git/model"
	"github.com/ogeyou/go-autth.git/storage/psql"
	"golang.org/x/crypto/argon2"
)

var (
	DBPass []byte
	UserID uint32
	Email  string
)

// Добавляем соль к пароль
// На время иссследований
func HashPass(plainPassword, salt string) []byte {
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

	pass := HashPass(user.Password, salt)

	var id int64

	const sql = "INSERT INTO users (login, password) VALUES($1, $2) RETURNING id"

	err := dbpool.QueryRow(ctx, sql, user.Login, pass).Scan(&id)
	var w http.ResponseWriter
	if err != nil {
		log.Println("Ощибка при добавлении нового пользователя в базу данных", err)
		http.Error(w, "Ощибка на сервере при добавлении нового пользователя в базу данных", http.StatusInternalServerError)
	}

	UserID := id

	return UserID
}

// Получаю данные при проверке логина пользователем из базы данных
func UserProtected(login string, password string) uint32 {
	// Соединение с экземпляром Postgres
	ctx := context.Background()
	dbpool := psql.Connect(ctx)
	defer dbpool.Close()

	
	var user model.User

	const sql = "SELECT id, password FROM users WHERE login = $1"

	rows, err := dbpool.Query(ctx, sql, login)

	if err != nil{
		log.Print("Err select db", err)
	}
	res := []model.User{}

	for rows.Next() {

		err = rows.Scan(&UserID, &DBPass)

		if err == pgx.ErrNoRows {
			log.Println("Ошибка при получении данных", err)
		} else if err != nil {
			fmt.Println("Данные успешно получены из функции дб")
		}

		salt := string(DBPass[0:8])
		
		var w http.ResponseWriter
		if !bytes.Equal(HashPass(user.Password, salt), DBPass) {
			http.Error(w, "Bad pass", http.StatusBadRequest)
		}

		CoursesBook := model.User{
			Id:       int64(UserID),
			Password: string(DBPass),
		}

		res = append(res, CoursesBook)
		if UserID == 0{
			panic(UserID)
		}
	}

	return UserID
}
