package main

import (
	"fmt"
	"net/http"

	"github.com/advn1/backend/global"
	"github.com/advn1/backend/internal/handlers"
	"github.com/advn1/backend/internal/middleware"
	"github.com/advn1/backend/internal/models/user"
	"github.com/google/uuid"
)

func main() {
	mux := http.NewServeMux()
	
	global.Users = append(global.Users, user.User{Name: "Alex", Email: "alexmail@google.com", Password_hash: "6u34rwuej", Id: uuid.New()}, user.User{Name: "John", Email: "johnmail@google.com", Password_hash: "jb84u43uifv", Id: uuid.New()}, user.User{Name: "Michael", Email: "michaelmail@google.com", Password_hash: "kdkm438989vjcx", Id: uuid.New()}, user.User{Name: "Smith", Email: "smithmail@google.com", Password_hash: "k438u9890md", Id: uuid.New()})
	port := "8080"
	fmt.Println(port)
	fmt.Println(global.Users)

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/users", handlers.UsersHandler)
	mux.HandleFunc("/users/", handlers.UsersHandlerByID)

	corsMux := middleware.EnableCORS(mux)

	fmt.Println("Listening to port:", port)
	err := http.ListenAndServe("localhost:"+port, corsMux)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// "/" handler
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Text")
}