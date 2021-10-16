package handlers

import (
	// "fmt"
	// "log"http"

	"fmt"
	"net/http"
)

func UserRegistration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	UserCreated()
}



func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintf(w, "This Login page")
	w.WriteHeader(http.StatusOK)
}