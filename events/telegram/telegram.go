package telegram

import (
	"fmt"

	"TGBot/client/telegram"
	"TGBot/events"
	"TGBot/storage"
)

type Dispatcher struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

func newDispatcher(client *telegram.Client, storage storage.Storage) *Dispatcher {
	return &Dispatcher{
		tg:      client,
		storage: storage,
	}
}

func (d *Dispatcher) Fetch(limit int) ([]events.Event, error) {
	update, err := d.tg.Updates(d.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events %w", err)
	}

	fetchResult := make([]events.Event, 0, len(update))

	for _, fetchedEvent := range update {
		fetchResult = append(fetchResult, toEvent(fetchedEvent))
	}
}

// func event is used to convert update into event
func toEvent(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}
