package handlers

import (
	// "fmt"
	// "log"http"

	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ogeyou/go-autth.git/model"
	"github.com/ogeyou/go-autth.git/storage"
)

func UserRegistration(w http.ResponseWriter, r *http.Request) {

	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Printf("")
	}

	insertID := storage.UserCreated(user)

	fmt.Println(insertID)

}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintf(w, "This Login page")
	w.WriteHeader(http.StatusOK)
}
