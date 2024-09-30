package telegramclient

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

// structure of incoming message
type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

// user who sent this message (TG API)
type From struct {
	Username string `json:"username"`
}

// conversations the message belongs to (TG API)
type Chat struct {
	ID int `json:"id"`
}
