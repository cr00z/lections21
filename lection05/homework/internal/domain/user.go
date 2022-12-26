package domain

type User struct {
	ID       int64  `json:"_" db:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
