package domain

type Message struct {
	ID         int64
	FromUserID int64
	ToUserID   int64
	Text       string
}
