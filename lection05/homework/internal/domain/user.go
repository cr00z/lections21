package domain

type User struct {
	ID       int64  `json:"_"`
	Username string `json:"username"`
	Password string `json:"password"`
}
