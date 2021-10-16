import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("api/registration", handlers.UserRegistration)

	http.ListenAndServe(":8081", router)
	
	fmt.Println("Start server ...")
}
