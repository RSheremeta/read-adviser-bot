package main

import (
	"context"
	"flag"
	"log"

	tgCl "github.com/RSheremeta/read-adviser-bot/clients/telegram"
	"github.com/RSheremeta/read-adviser-bot/consumer/eventconsumer"
	"github.com/RSheremeta/read-adviser-bot/events/telegram"

	//	"github.com/RSheremeta/read-adviser-bot/storage/files"
	"github.com/RSheremeta/read-adviser-bot/storage/sqlite"
)

const (
	tokenFlagName  = "tg-bot-token"
	tgBotHost      = "api.telegram.org"
	storagePath    = "files_storage"
	storageSqlPath = "data/sqlite/storage.db"
	batchSize      = 100
)

func main() {
	token := mustToken()

	tgClient := tgCl.New(tgBotHost, token)

	// an alternative option to use - a file storage via Gob - if use this, change storage variable passing into tg.new() func
	// storage := files.New(storagePath)
	// or
	// sqlite storage
	storage, err := sqlite.New(storageSqlPath)
	if err != nil {
		log.Fatal("cannot connect to storage: ", err)
	}

	if err = storage.Init(context.TODO()); err != nil {
		log.Fatal("cannot init storage: ", err)
	}

	eventsProcessor := telegram.New(tgClient, storage)

	log.Println("service started")

	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatalf("service stopped due to the error: %s", err.Error())
	}
}

func mustToken() string {
	t := flag.String(
		tokenFlagName,
		"",
		"token for the Telegram bot accessing",
	)
	flag.Parse()

	if *t == "" {
		log.Fatal("the token value is not specified")
	}

	return *t
}
