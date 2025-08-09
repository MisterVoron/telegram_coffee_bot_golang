package types

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

type From struct {
	ID int64 `json:"id"`
}

type Chat struct {
	ID int64 `json:"id"`
}
