package telegram

import "TGBot/client/telegram"

type Dispatcher struct {
	tg     *telegram.Client
	offset int
}
