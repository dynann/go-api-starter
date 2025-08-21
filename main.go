package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type User struct {
	Name string `json:"name"`
	Age  int32  `json:"age"`
}

var users []User

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	users = append(users, user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(users)
}

func getUserByAge(w http.ResponseWriter, r *http.Request) {
	ageStr := chi.URLParam(r, "id")
	ageInt, err := strconv.ParseInt(ageStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid age parameter", http.StatusBadRequest)
		return
	}
	var userByAge []User
	for _, u := range users {
		if int32(ageInt) == u.Age {
			userByAge = append(userByAge, u)
		}
	}

	if len(userByAge) == 0 {
		http.Error(w, "No user found with the given age", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userByAge)
}

func deleteUserByAge(w http.ResponseWriter, r *http.Request) {
	ageStr := chi.URLParam(r, "id")
	ageInt, err := strconv.ParseInt(ageStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid parameter", http.StatusBadRequest)
	}
	var newUserByAge []User
	for _, u := range users {
		if int32(ageInt) != u.Age {
			users = append(newUserByAge, u)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUserByAge)
}

func updateUserByAge(w http.ResponseWriter, r *http.Request) {
	ageStr := chi.URLParam(r, "id")
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	ageInt, err := strconv.ParseInt(ageStr, 10 , 32)
	if err != nil {
		http.Error(w, "invalid parameter", http.StatusBadRequest)
	}
   updated := false
    for i, u := range users {
        if int32(ageInt) == u.Age {
            users[i] = user // Replace user at index
            updated = true
        }
    }

    if !updated {
        http.Error(w, "No user found with the given age", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func main() {

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the API!"))
	})
	r.Get("/users", userHandler)
	r.Post("/users", createUser)
	r.Get("/users/{id}", getUserByAge)
	r.Delete("/users/{id}", deleteUserByAge)
	r.Put("/users/{id}", updateUserByAge)
	fmt.Println("Server running http://localhost:8000")
	http.ListenAndServe(":8000", r)
}
