package telegram

const msgHelp = `I can save and keep your pages.

In order to save the page, just send me link to it.

In order to get a random page from your list, send me command /rnd.
Caution! After that rnd command, this page will be removed from your list!`

const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command"
	msgNoSavedPages   = "You have no saved pages"
	msgSaved          = "Saved!"
	msgAlreadyExists  = "You have already have this page in your list"
)
