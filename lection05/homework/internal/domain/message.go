package domain

type Message struct {
	ID         int64  `json:"id"`
	FromUserID int64  `json:"from_user"`
	ToUserID   int64  `json:"to_user"`
	Text       string `json:"text"`
}
