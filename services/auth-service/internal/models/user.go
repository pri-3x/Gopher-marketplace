package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // The "-" tag means this field won't be included in JSON output
	Email    string `json:"email"`
}
