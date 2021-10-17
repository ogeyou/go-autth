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
	UserID uint32
	ID     string
}
type ctxKey int

const sessionKey ctxKey = 1

var (
	noAuthUrls = map[string]struct{}{
		"/user/login":  struct{}{},
		"/user/logout": struct{}{},
		"/user/reg":    struct{}{},
		"/courses":    struct{}{},
		"/":            struct{}{},
	}
)

var (
	ErrNoAuth = errors.New("No session found")
)

func AuthMiddleware(dbpool *pgxpool.Pool, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}
		sess, err := CheckSession(dbpool, r)
		if err != nil {
			http.Error(w, "Вы не авторизованы", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), sessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

func CreateSession(w http.ResponseWriter, r *http.Request, dbpol *pgxpool.Pool, UserID uint32) error {
	ctx := context.Background()
	dbpool := psql.Connect(ctx)
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
func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		return nil, ErrNoAuth
	}
	return sess, nil
}

func DestroySession(w http.ResponseWriter, r *http.Request, dbpol *pgxpool.Pool) error {
	ctx := context.Background()
	dbpool := psql.Connect(ctx)
	defer dbpool.Close()
	sess, err := SessionFromContext(r.Context())
	if err == nil {
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

func Index(w http.ResponseWriter, r *http.Request) {
	_, err := SessionFromContext(r.Context())
	if err != nil {
		http.Redirect(w, r, "/user/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/courses/", http.StatusFound)
}

func Courses(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Поздравляем, Вы зарегестрированы!")
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	dbpool := psql.Connect(ctx)
	defer dbpool.Close()
	DestroySession(w, r, dbpool)
	http.Redirect(w, r, "/", http.StatusFound)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	dbpool := psql.Connect(ctx)
	defer dbpool.Close()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintf(w, "This Login page")
	w.WriteHeader(http.StatusOK)

	DestroySession(w, r, dbpool)
	http.Redirect(w, r, "/user/login", http.StatusFound)
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
