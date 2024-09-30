package telegram

// structure of telegram API getUpdates chapter "getting updates"
type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

// for more info see "making requests" chapter of telegram API
type UpdateResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}
