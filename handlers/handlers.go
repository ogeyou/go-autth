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
	ErrNoAuth = errors.New("No session found")
)

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		return nil, ErrNoAuth
	}
	return sess, nil
}
func Index(w http.ResponseWriter, r *http.Request) {
	_, err := SessionFromContext(r.Context())
	if err != nil {
		http.Redirect(w, r, "/user/login", http.StatusFound)
		return
	}
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

func CreateSession(w http.ResponseWriter, r *http.Request, dbpol *pgxpool.Pool, UserID uint32) error {
	ctx := context.Background()
	dbpool := psql.Connect(ctx)
	defer dbpool.Close()
	fmt.Println("Смотри, значение передается или нет куки", UserID)
	sessID := storage.RandStringRunes(32)
	_, err := dbpool.Exec(ctx, "insert into sessions(id, user_id) VALUES($1, $2);", sessID, UserID)
	if err != nil{
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

func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintf(w, "This Login page")
	w.WriteHeader(http.StatusOK)
}

func Courses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintf(w, "Поздравляем, Вы зарегестрированы!")
	w.WriteHeader(http.StatusOK)
}	

