package main

import (
	"context"
	"fmt"

	"net/http"

	"github.com/ogeyou/go-autth.git/handlers"
	"github.com/ogeyou/go-autth.git/storage/psql"
)

func main() {
	// Соединение с экземпляром Postgres
	ctx := context.Background()
	dbpool := psql.Connect(ctx)
	defer dbpool.Close()
	
	router := http.NewServeMux()

	fmt.Println("Start server ...")
	router.HandleFunc("/courses/", handlers.Courses)

	router.HandleFunc("/user/reg", handlers.UserRegistration)
	router.HandleFunc("/user/login", handlers.UserLogin)
	router.HandleFunc("/user/logout", handlers.UserLogout)

	router.HandleFunc("/", handlers.Index)
	http.Handle("/", handlers.AuthMiddleware(dbpool, router))

	http.ListenAndServe(":8081", nil)
}
