package telegram

import (
	"read-adviser-bot/clients/telegram"
	"read-adviser-bot/events"
	"read-adviser-bot/lib/errwrap"
	"read-adviser-bot/storage"
)

type Processor struct {
	tgClient *telegram.Client
	offset   int
	storage  storage.Storage
}

func NewProcessor(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tgClient: client,
		storage:  storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	update, err := p.tgClient.Updates(p.offset, limit)
	if err != nil {
		return nil, errwrap.Wrap("can't get events", err)
	}

	res := make([]events.Event, 0, len(update))

	for _, upd := range update {
		res = append(res, event(upd))
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	res := events.Event{
		Type: fetchType(upd),
		Text: fetchText(upd),
	}
}

func fetchType(upd telegram.Update) events.Type {

}

func fetchText(upd telegram.Update) telegram.Message {

}
