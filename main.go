package main

import (
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/ogeyou/go-autth.git/handlers"
)

func main() {
	router := mux.NewRouter()
	fmt.Println("Start server ...")
	router.HandleFunc("/api/registration", handlers.UserRegistration)
http.ListenAndServe(":8081", router)

	

}
