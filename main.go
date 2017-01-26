package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	jsonMIME = "application/json; charset=utf-8"
	htmlMIME = "text/html; charset=utf-8"

	MsgMethodNotAllowed = "method not allowed"
)

var Users = []*User{
	&User{"admin", "admin", "1f062f19e4e581f4", true},
	&User{"demo", "demo", "756a79a17f4feaab", false},
}

// Get a user by token, returns an error if it's not a valid one.
func UserByToken(token string) (*User, error) {
	for _, user := range Users {
		if user.Token == token {
			return user, nil
		}
	}
	return nil, errors.New("invalid token")
}

// Serves an error message as JSON.
func ServeError(rw http.ResponseWriter, status int, msg string) {
	rw.Header().Set("Content-Type", jsonMIME)
	rw.WriteHeader(status)
	rw.Write([]byte(`{"error": "` + msg + `"}`))
}

// Generate a JWT token usable for further authentication.
// POST: { "username": "admin", "password": "admin" }
// Response: { "token": "1f062f19e4e581f4" }
func HandleLogin(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		ServeError(rw, http.StatusMethodNotAllowed, MsgMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		ServeError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		ServeError(rw, http.StatusBadRequest, err.Error())
		return
	}

	for _, user := range Users {
		if user.Username == data.Username && user.Password == data.Password {
			rw.Header().Set("Content-Type", jsonMIME)
			rw.Write([]byte(`{"token": "` + user.Token + `"}`))
			return
		}
	}
	ServeError(rw, http.StatusUnauthorized, "invalid username or password")
}

// Returns the authorized user.
// Response: { "username": "admin", "admin": true }
func HandleMe(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		ServeError(rw, http.StatusMethodNotAllowed, MsgMethodNotAllowed)
		return
	}

	user, err := UserByToken(req.URL.Query().Get("token"))
	if err != nil {
		ServeError(rw, http.StatusUnauthorized, err.Error())
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		ServeError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	rw.Header().Set("Content-Type", jsonMIME)
	rw.Write(res)
}

func main() {
	indexHTML, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", htmlMIME)
		rw.Write(indexHTML)
	})

	http.HandleFunc("/api/login", HandleLogin)
	http.HandleFunc("/api/me", HandleMe)

	http.ListenAndServe("0.0.0.0:8000", nil)
}
