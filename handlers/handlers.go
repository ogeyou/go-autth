package handlers

import (
	// "fmt"
	// "log"http"

	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ogeyou/go-autth.git/model"
	"github.com/ogeyou/go-autth.git/storage"
	"github.com/ogeyou/go-autth.git/storage/psql"
)

type Session struct {
	UserID uint32 `json:"user_id,omitempty"`
	ID     string `json:"id,omitempty"`
}
type ctxKey int

const sessionKey ctxKey = 1

var (
	noAuthUrls = map[string]struct{}{
		"/user/login":  struct{}{},
		"/user/logout": struct{}{},
		"/user/reg":    struct{}{},
		"/":            struct{}{},
	}
)

var (
	ErrNoAuth = errors.New("No session found")
)

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		return nil, ErrNoAuth
	}
	return sess, nil
}

func CheckSession(dbpool *pgxpool.Pool, r *http.Request) (*Session, error) {

	ctx := context.Background()

	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		log.Println("CheckSession no cookie")
		return nil, ErrNoAuth
	}

	sess := &Session{}
	row := dbpool.QueryRow(ctx, `SELECT user_id FROM sessions WHERE id = $1;`, sessionCookie.Value)
	if row != nil {
		fmt.Println(row)
	}
	err = row.Scan(&sess.UserID)
	if err == pgx.ErrNoRows {
		log.Println("CheckSession no rows")
		return nil, ErrNoAuth
	} else if err != nil {
		log.Println("CheckSession err:", err)
		return nil, err
	}

	sess.ID = sessionCookie.Value
	return sess, nil
}

func CreateSession(w http.ResponseWriter, r *http.Request, dbpool *pgxpool.Pool, UserID uint32) error {

	ctx := context.Background()
	defer dbpool.Close()
	fmt.Println("Смотри, значение передается или нет куки", UserID)
	sessID := storage.RandStringRunes(32)
	_, err := dbpool.Exec(ctx, "insert into sessions(id, user_id) VALUES($1, $2);", sessID, UserID)
	if err != nil {
		fmt.Println(err)
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessID,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	return nil
}

func DestroySession(w http.ResponseWriter, r *http.Request, dbpool *pgxpool.Pool) error {

	ctx := context.Background()
	defer dbpool.Close()
	sess, err := SessionFromContext(r.Context())
	if err != nil {
		fmt.Println("EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE", err)
	} else {
		dbpool.Exec(ctx, "DELETE FROM sessions WHERE user_id = $1", sess.ID)
		fmt.Println("Данные куки удалены из базы данных")
	}
	cookie := http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	return nil
}
func AuthMiddleware(dbpool *pgxpool.Pool, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}
		sess, err := CheckSession(dbpool, r)
		if err != nil {
			log.Println("Сообщение от гуфера - Вы не прошлт проверку Авторизации")
			http.Error(w, "Вы не авторизованы", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), sessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func Index(w http.ResponseWriter, r *http.Request) {

	_, err := SessionFromContext(r.Context())
	if err != nil {
		http.Redirect(w, r, "/user/login/", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/courses/", http.StatusFound)
}

func Courses(w http.ResponseWriter, r *http.Request) {

	log.Println(http.StatusOK)

	fmt.Fprintf(w, "Поздравляем, Вы зарегестрированы!")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx := context.Background()
	dbpool := psql.Connect(ctx)
	defer dbpool.Close()



	var user model.User

	// err := json.NewDecoder(r.Body).Decode(&user.Login)


	LoginName, err := storage.UserProtected(user.Login)
	
	if err != nil {
		log.Println("Ошибка при получении данных", err)
	} else if err == nil {
		fmt.Println("Данные успешно получены из функции дб")
	}

	fmt.Println(LoginName)

	// salt := string(storage.DBPass[0:8])

	// if !bytes.Equal(storage.HashPass(user.Password, salt), storage.DBPass) {
	// 	http.Error(w, "Bad pass", http.StatusBadRequest)
	// 	return
	// }

	CreateSession(w, r, dbpool, storage.UserID)
	http.Redirect(w, r, "/courses", http.StatusFound)
}

func UserRegistration(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	dbpool := psql.Connect(ctx)
	defer dbpool.Close()

	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Printf("")
	}

	UserID := storage.UserCreated(user)
	fmt.Println("Уже чуть ближе к успеху ", UserID)

	CreateSession(w, r, dbpool, uint32(UserID))
	http.Redirect(w, r, "/courses", http.StatusFound)

}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx := context.Background()
	dbpool := psql.Connect(ctx)
	defer dbpool.Close()
	DestroySession(w, r, dbpool)
	http.Redirect(w, r, "/", http.StatusFound)
}
