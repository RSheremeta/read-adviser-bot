package telegram

import (
	"errors"

	"github.com/RSheremeta/read-adviser-bot/clients/telegram"
	"github.com/RSheremeta/read-adviser-bot/events"
	"github.com/RSheremeta/read-adviser-bot/lib/e"
	"github.com/RSheremeta/read-adviser-bot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	// offset isn't set here since it starts from 0 by default
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("cannot get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("cannot process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("cannot process message", err)
	}

	if err := p.doCmd(meta.ChatID, event.Text, meta.Username); err != nil {
		return e.Wrap("cannot process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta) // type assertion
	if !ok {
		return Meta{}, e.Wrap("cannot get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(u telegram.Update) events.Event {
	updType := fetchType(u)

	res := events.Event{
		Type: updType,
		Text: fetchText(u),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   u.Message.Chat.ID,
			Username: u.Message.From.Username,
		}
	}

	return res
}

func fetchType(u telegram.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func fetchText(u telegram.Update) string {
	if u.Message == nil {
		return ""
	}

	return u.Message.Text
}
