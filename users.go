package main

type User struct {
	Username string `json:"username"`
	Password string `json:"-"`
	Token    string `json:"-"`
	Admin    bool   `json:"admin"`
}
