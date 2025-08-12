package main

import (
	"flag"
	"log"
	event_consumer "read_adviser_bot/consumer/event-consumer"
	"read_adviser_bot/storage/files"

	tgClient "read_adviser_bot/clients/telegram"
	"read_adviser_bot/events/telegram"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

// 7045335569:AAGOYKxTf2cycD1T-BCI3LoNNAbQr_6Abs4.

func main() {

	// s := files.New(storagePath)
	// s, err := sqlite.New(sqliteStoragePath)
	// if err != nil {
	// 	log.Fatal("can't connect to storage: ", err)
	// }

	// if err := s.Init(context.TODO()); err != nil {
	// 	log.Fatal("can't init storage: ", err)
	// }

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	// bot -tg-bot-token 'my token'
	token := flag.String("tg-bot-token", "", "token for acces to telegram bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
