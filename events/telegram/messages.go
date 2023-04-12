package telegram

const msgHelp = `I can save and keep your pages. Also, I can offer you them to read.

In order to save the page, just send me a link to it.

In order to get a random page from your list, just send me command "/rnd"
Caution! After sending, the random page will be removed from your list!`

const msgHello = "Hi there! 👋 \n\n" + msgHelp

const (
	msgUnknownCommand   = "Unknown command 🤔 "
	msgNoSavedPages     = "You have no saved pages 🛑 "
	msgNoSavingsHistory = "You haven't save anything yet. Send me the first link to save 🙌 "
	msgSaved            = "Saved! 👌 "
	msgAlreadyExists    = "You already have this page in your list 😇 "
)
