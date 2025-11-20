package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User scheme for in-database model
type User struct {
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password_hash string    `json:"password_hash"`
	Id            uuid.UUID `json:"id"`
}

// Custom Http Error
type HttpError struct {
	Message string
	Code    int
}

// Formatting Http Error
func (e *HttpError) Error() string {
	return fmt.Sprintf("%s Code: %d", e.Message, e.Code)
}

// initial database
var users []User = make([]User, 0)

// enable cors on every request
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	users = append(users, User{Name: "Alex", Email: "alexmail@google.com", Password_hash: "6u34rwuej", Id: uuid.New()}, User{Name: "John", Email: "johnmail@google.com", Password_hash: "jb84u43uifv", Id: uuid.New()}, User{Name: "Michael", Email: "michaelmail@google.com", Password_hash: "kdkm438989vjcx", Id: uuid.New()}, User{Name: "Smith", Email: "smithmail@google.com", Password_hash: "k438u9890md", Id: uuid.New()})
	port := "8080"
	fmt.Println(port)
	fmt.Println(users)

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/users", usersHandler)
	mux.HandleFunc("/users/", usersHandlerByID)

	corsMux := enableCORS(mux)

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

// "/users" handler 
func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/users" || r.URL.Path == "/users/" {
		fmt.Println("ALL USERS")
		switch r.Method {
		case http.MethodGet:
			getAllUsers(w)
		case http.MethodPost:
			// note: why uuid.Nil?
			createNewUser(w, r, uuid.Nil)
		default:
			http.Error(w, "Method"+ r.Method + "is not allowed.", http.StatusNotFound)
			return
		}
	}
}

// GET "/users"
func getAllUsers(w http.ResponseWriter) {
	data, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Error couldn't parse users", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(data))
}

// POST "/users"
func createNewUser(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	user, httpErr := getUserDataFromForm(r, id)
	if httpErr != nil {
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}
	user.Id = uuid.New()

	users = append(users, user)

	data, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error couldn't parse user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(data))
}

// GET "/users/id"
func usersHandlerByID(w http.ResponseWriter, r *http.Request) {
	parsedURL := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	fmt.Println(parsedURL, len(parsedURL))

	// if url path doesn't contain user id
	if r.URL.Path == "/users/" {
		usersHandler(w, r)
		return
	}

	if len(parsedURL) < 2 || len(parsedURL) > 2 {
		http.Error(w, "Unknown route", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(parsedURL[1])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getUserByID(w, id)
	case http.MethodDelete:
		deleteUser(w, id)
	case http.MethodPut:
		updateUserData(w,r,id)
	default:
		http.Error(w, "Unknown method.", http.StatusNotFound)
		return
	}
}
// PUT "/users/id"
func updateUserData(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
		user, httpErr := getUserDataFromForm(r, id)
		if httpErr != nil {
			http.Error(w, httpErr.Message, httpErr.Code)
			return
		}

		for idx, usr := range users {
			if usr.Id == id {
				users[idx].Name = user.Name
				users[idx].Email = user.Email
				users[idx].Password_hash = user.Password_hash
				data, err := json.Marshal(users[idx])
				if err != nil {
					http.Error(w, "Error couldn't parse user", http.StatusInternalServerError)
					return
				}
				fmt.Fprint(w, string(data))
				return
			}
		}
		http.Error(w, "User not found", http.StatusNotFound)
}

// GET "/users/id"
func getUserByID(w http.ResponseWriter, id uuid.UUID) {
	fmt.Println("ID GET", id)
	// handle GET logic
	var foundUser *User
	for _, user := range users {
		if user.Id == id {
			foundUser = &user
			break
		}
	}

	if foundUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	data, err := json.Marshal(foundUser)
	if err != nil {
		http.Error(w, "Error couldn't parse user", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(data))
}

// DELETE "/users/id"
func deleteUser(w http.ResponseWriter, id uuid.UUID) {
	for i, user := range users {
		if user.Id == id {
			users = append(users[:i], users[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

// fetch user data from form by id
func getUserDataFromForm(r *http.Request, id uuid.UUID) (User, *HttpError) {
	err := r.ParseForm()
	if err != nil {
		return User{}, &HttpError{Message: "Error parsing form data", Code: http.StatusBadRequest}
	}

	name := r.PostFormValue("name")
	if len(name) < 3 {
		return User{}, &HttpError{Message: "Name is too short", Code: http.StatusBadRequest}
	}
	email := r.PostFormValue("email")
	_, err = mail.ParseAddress(email)
	if err != nil {
		fmt.Println(err)
		return User{}, &HttpError{Message: "Email is not correct", Code: http.StatusBadRequest}
	}

	// check email unique
	for _, user := range users {
		if user.Email == email && user.Id != id {
			return User{}, &HttpError{Message: "Email already exists in database", Code: http.StatusBadRequest}
		}
	}

	password := r.PostFormValue("password")
	if len(password) < 8 {
		return User{}, &HttpError{Message: "Password is too short", Code: http.StatusBadRequest}
	}

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, &HttpError{Message: "Couldn't hash password", Code: http.StatusInternalServerError}
	}

	user := User{Name: name, Email: email, Password_hash: string(hashed_password)}
	fmt.Println("Get user", user)
	return user, nil
}
