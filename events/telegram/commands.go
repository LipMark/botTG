package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"TGBot/storage"
)

// TODO :closures5

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

// func savePage is a client cmd to save a given link-article
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

	isExists, err := d.storage.IsExists(context.Background(), page)
	if err != nil {
		return err
	}
	if isExists {
		return d.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := d.storage.Save(context.Background(), page); err != nil {
		return err
	}

	if err := d.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

// func sendRandom is a client cmd to take a random article from storage and remove it afterwards
func (d *Dispatcher) sendRandom(chatID int, username string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("can't execute sendRandom command %w", err)
		}
	}()

	page, err := d.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return d.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := d.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return d.storage.Remove(context.Background(), page)
}

// func sendHelp is used by user to get a list of bot commands
func (d *Dispatcher) sendHelp(chatID int) error {
	return d.tg.SendMessage(chatID, msgHelp)
}

// func sendHello is used to make a hello-message for user
func (d *Dispatcher) sendHello(chatID int) error {
	return d.tg.SendMessage(chatID, msgHello)
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
