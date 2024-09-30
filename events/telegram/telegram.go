package telegram

import (
	"errors"
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

// info from tg messenger
type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func newDispatcher(client *telegram.Client, storage storage.Storage) *Dispatcher {
	return &Dispatcher{
		tg:      client,
		storage: storage,
	}
}

func (d *Dispatcher) Fetch(limit int) ([]events.Event, error) {
	updates, err := d.tg.Updates(d.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events %w", err)
	}

	fetchResult := make([]events.Event, 0, len(updates))

	for _, fetchedEvent := range updates {
		fetchResult = append(fetchResult, toEvent(fetchedEvent))
	}

	d.offset = updates[len(updates)-1].ID + 1

	return fetchResult, nil
}

func (d *Dispatcher) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return d.processMessage(event)
	default:
		return fmt.Errorf("can't process message %w", ErrUnknownEventType)
	}
}

// func processMessage processing an actual message, not an update.
func (d *Dispatcher) processMessage(event events.Event) error {
	meta, err := toMeta(event)
	if err != nil {
		return fmt.Errorf("can't convert meta %w", err)
	}

	if err := d.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return fmt.Errorf("can't exec Cmd %w", err)
	}

	return nil
}

// func toMeta identifying meta, else return error
func toMeta(event events.Event) (Meta, error) {
	meta, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("can't get meta %w", ErrUnknownMetaType)
	}

	return meta, nil
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

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.EventType {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
