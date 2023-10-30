package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
)

var (
	username = "tim"
	password = "password"
	token    = ""
)

func generateHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", hash)
}

// curl -X POST -d "username=tim&password=password" http://localhost:8080/login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	user := r.FormValue("username")
	pass := r.FormValue("password")

	if user == username && generateHash(pass) == generateHash(password) {
		token = generateHash(username + password)

		w.Write([]byte(fmt.Sprintf("Login successful! Token: %s", token)))
	} else {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	}
}

// curl -X GET -H "Authorization: your-secret-token" http://localhost:8080/protected
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	requestedToken := r.Header.Get("Authorization")
	if requestedToken != token {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("You have access to the protected resource!"))
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/protected", protectedHandler)

	fmt.Println("Starting the server on :8080...")
	http.ListenAndServe(":8080", nil)
}
