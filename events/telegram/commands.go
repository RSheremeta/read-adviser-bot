package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/RSheremeta/read-adviser-bot/clients/telegram"
	"github.com/RSheremeta/read-adviser-bot/lib/e"
	"github.com/RSheremeta/read-adviser-bot/storage"
)

const (
	StartCmd = "/start"
	HelpCmd  = "/help"
	RndCmd   = "/rnd"
)

func (p *Processor) doCmd(chatID int, text, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command %q from %q", text, username)

	// add page link

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case StartCmd:
		return p.sendHello(chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	case RndCmd:
		return p.sendRandom(chatID, username)
	default:
		return p.tg.SendMsg(chatID, msgUnknownCommand)
	}

	// get random page: /rnd
	// /help
	// /start: hi + help
}

func (p *Processor) savePage(chatID int, pageURL, username string) (err error) {
	defer func() { err = e.WrapIfErr("cannot do command 'save page'", err) }()

	// sendMsg := msgSender(chatID, p.tg)

	page := &storage.Page{
		URL:      pageURL,
		Username: username,
	}

	isExist, err := p.storage.IsExist(context.TODO(), page)
	if err != nil {
		return err
	}
	if isExist {
		// sendMsg(msgAlreadyExists)
		return p.tg.SendMsg(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(context.TODO(), page); err != nil {
		return err
	}

	if err := p.tg.SendMsg(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("cannot do command 'send random'", err) }()

	page, err := p.storage.PickRandom(context.TODO(), username)
	if err != nil && errors.Is(err, storage.ErrNoUserSavingsHistory) {
		return p.tg.SendMsg(chatID, msgNoSavingsHistory)
	}
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMsg(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMsg(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(context.TODO(), page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMsg(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMsg(chatID, msgHello)
}

func msgSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMsg(chatID, msg)
	}
}

func isAddCmd(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
