// "/users" handler
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/advn1/backend/global"
	"github.com/advn1/backend/internal/jsonutil"
	httperror "github.com/advn1/backend/internal/models/http_error"
	"github.com/advn1/backend/internal/models/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ALL USERS")
	switch r.Method {
	case http.MethodGet:
		GetAllUsers(w)
	case http.MethodPost:
		// note: why uuid.Nil?
		CreateNewUser(w, r, uuid.Nil)
	default:
		w.Header().Set("Content-Type", "application/json")
		jsonutil.JSONError(w, "Method"+ r.Method + "is not allowed.", http.StatusNotFound)
		return
	}
}

// GET "/users"
func GetAllUsers(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(global.Users)
	if err != nil {
		jsonutil.JSONError(w, "Error couldn't parse users", http.StatusInternalServerError)
		return
	}
	
	fmt.Fprint(w, string(data))
}

// POST "/users"
func CreateNewUser(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	w.Header().Set("Content-Type", "application/json")

	user, httpErr := GetUserDataFromForm(r, id)
	if httpErr != nil {
		jsonutil.JSONError(w, httpErr.Message, httpErr.Code)
		return
	}
	user.Id = uuid.New()

	global.Users = append(global.Users, user)

	data, err := json.Marshal(user)
	if err != nil {
		jsonutil.JSONError(w, "Error couldn't parse user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(data))
}

// GET "/users/id"
func UsersHandlerByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parsedURL := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	fmt.Println(parsedURL, len(parsedURL))

	// if url path doesn't contain user id
	if r.URL.Path == "/users/" {
		UsersHandler(w, r)
		return
	}

	if len(parsedURL) < 2 || len(parsedURL) > 2 {
		jsonutil.JSONError(w, "Unknown route", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(parsedURL[1])
	if err != nil {
		jsonutil.JSONError(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetUserByID(w, id)
	case http.MethodDelete:
		DeleteUser(w, id)
	case http.MethodPut:
		UpdateUserData(w,r,id)
	default:
		jsonutil.JSONError(w, "Unknown method.", http.StatusNotFound)
		return
	}
}
// PUT "/users/id"
func UpdateUserData(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
		w.Header().Set("Content-Type", "application/json")

		user, httpErr := GetUserDataFromForm(r, id)
		if httpErr != nil {
			jsonutil.JSONError(w, httpErr.Message, httpErr.Code)
			return
		}

		for idx, usr := range global.Users {
			if usr.Id == id {
				global.Users[idx].Name = user.Name
				global.Users[idx].Email = user.Email
				global.Users[idx].Password_hash = user.Password_hash
				data, err := json.Marshal(global.Users[idx])
				if err != nil {
					jsonutil.JSONError(w, "Error couldn't parse user", http.StatusInternalServerError)
					return
				}
				fmt.Fprint(w, string(data))
				return
			}
		}
		jsonutil.JSONError(w, "User not found", http.StatusNotFound)
}

// GET "/users/id"
func GetUserByID(w http.ResponseWriter, id uuid.UUID) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("ID GET", id)
	// handle GET logic
	var foundUser *user.User
	for _, user := range global.Users {
		if user.Id == id {
			foundUser = &user
			break
		}
	}

	if foundUser == nil {
		jsonutil.JSONError(w, "User not found", http.StatusNotFound)
		return
	}

	data, err := json.Marshal(foundUser)
	if err != nil {
		jsonutil.JSONError(w, "Error couldn't parse user", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(data))
}

// DELETE "/users/id"
func DeleteUser(w http.ResponseWriter, id uuid.UUID) {
	for i, user := range global.Users {
		if user.Id == id {
			global.Users = append(global.Users[:i], global.Users[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	jsonutil.JSONError(w, "User not found", http.StatusNotFound)
}

func GetUserDataFromForm(r *http.Request, id uuid.UUID) (user.User, *httperror.HttpError) {
	err := r.ParseForm()
	if err != nil {
		return user.User{}, &httperror.HttpError{Message: "Error parsing form data", Code: http.StatusBadRequest}
	}

	name := r.PostFormValue("name")
	if len(name) < 3 {
		return user.User{}, &httperror.HttpError{Message: "Name is too short", Code: http.StatusBadRequest}
	}
	email := r.PostFormValue("email")
	_, err = mail.ParseAddress(email)
	if err != nil {
		fmt.Println(err)
		return user.User{}, &httperror.HttpError{Message: "Email is not correct", Code: http.StatusBadRequest}
	}

	// check email unique
	for _, usr := range global.Users {
		if usr.Email == email && usr.Id != id {
			return user.User{}, &httperror.HttpError{Message: "Email already exists in database", Code: http.StatusBadRequest}
		}
	}

	password := r.PostFormValue("password")
	if len(password) < 8 {
		return user.User{}, &httperror.HttpError{Message: "Password is too short", Code: http.StatusBadRequest}
	}

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user.User{}, &httperror.HttpError{Message: "Couldn't hash password", Code: http.StatusInternalServerError}
	}

	user := user.User{Name: name, Email: email, Password_hash: string(hashed_password)}
	fmt.Println("Get user", user)
	return user, nil
}