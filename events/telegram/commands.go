package telegram

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"TGBot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

// func doCmd serves as a wrapper, executing the right function when the client calls it
func (d *Dispatcher) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	if isAddCmd(text) {
		return d.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return d.sendRandom(chatID, username)
	case HelpCmd:
		return d.sendHelp(chatID)
	case StartCmd:
		return d.sendHello(chatID)
	default:
		return d.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (d *Dispatcher) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("can't execute save command %w", err)
		}
	}()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := d.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return d.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := d.storage.Save(page); err != nil {
		return err
	}

	if err := d.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

// func isAddCmd is an util command to check for AddCmd
func isAddCmd(text string) bool {
	return isURL(text)
}

// func isURL checks for page entity
func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
