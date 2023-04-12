package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	client "github.com/RSheremeta/read-adviser-bot/clients/telegram"
	"github.com/RSheremeta/read-adviser-bot/consumer/eventconsumer"
	"github.com/RSheremeta/read-adviser-bot/events/telegram"
	"github.com/RSheremeta/read-adviser-bot/storage"

	"github.com/RSheremeta/read-adviser-bot/storage/files"
	"github.com/RSheremeta/read-adviser-bot/storage/sqlite"
)

const (
	tokenFlagName = "tg_bot_token"
	tgBotHost     = "api.telegram.org"
	batchSize     = 100
)

const (
	filesStoragePath    = "files_storage"
	sqlStoragePath      = "data/sqlite/storage.db"
	storageTypeFlagName = "storage_type"
	sqliteStorageOpt    = "sqlite"
	filesStorageOpt     = "files"
)

func main() {
	var token, storageType string

	flag.StringVar(
		&token,
		tokenFlagName,
		"",
		"token for the Telegram bot accessing",
	)

	flag.StringVar(
		&storageType,
		storageTypeFlagName,
		"",
		fmt.Sprintf("type of storage for the links usage. the options are %q or %q. and %q is by default",
			sqliteStorageOpt, filesStorageOpt, sqliteStorageOpt),
	)

	if storageType == "" {
		storageType = sqliteStorageOpt
	}
	flag.Parse()

	if token == "" {
		log.Fatal("the token value is not specified")
	}

	tgClient := client.New(tgBotHost, token)

	var storage storage.Storage

	if storageType == filesStorageOpt {
		storage = files.New(filesStoragePath)
		log.Println("initializing storage of files type")
	} else if storageType == sqliteStorageOpt {
		storage, err := sqlite.New(sqlStoragePath)
		if err != nil {
			log.Fatal("cannot connect to the sqlite storage: ", err)
		}

		if err = storage.Init(context.TODO()); err != nil {
			log.Fatal("cannot init the sqlite storage: ", err)
		}
		log.Println("initializing storage of sqlite type")
	}

	eventsProcessor := telegram.New(tgClient, storage)

	log.Println("service is started")

	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatalf("service stopped due to the error: %s", err.Error())
	}
}
