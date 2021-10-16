package main

import (
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	fmt.Println("Start server ...")

	router.HandleFunc("/api/registration", handlers.UserRegistration)
	// router.HandleFunc("/api/login", handlers.UserLogin)
	http.ListenAndServe(":8081", router)
}
