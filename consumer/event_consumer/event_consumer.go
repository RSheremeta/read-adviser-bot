package eventconsumer

import (
	"log"
	"time"

	"github.com/RSheremeta/read-adviser-bot/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s\n", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Println(err)
			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, e := range events {
		log.Printf("got new event: %s\n", e.Text)

		if err := c.processor.Process(e); err != nil {
			log.Printf("cannot handle event %s\n", err.Error())
			continue
		}
	}

	return nil
}
