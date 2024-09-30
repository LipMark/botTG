package telegram

// structure of telegram API getUpdates chapter "getting updates"
type Update struct {
	ID      int    `json:"update_id"`
	Message string `json:"message"`
}

// for more info see "making requests" chapter of telegram API
type UpdateResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}
