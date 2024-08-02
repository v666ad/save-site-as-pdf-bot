package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"save-site-as-pdf-bot/config"
	"save-site-as-pdf-bot/handlers/update"
	"save-site-as-pdf-bot/models"
)

func main() {
	config, err := config.LoadFromFile("config.json")
	if err != nil {
		log.Panic(err)
	}
	log.Println("загружена конфигурация", config.Name)

	bot, err := tgbotapi.NewBotAPI(config.TelegramApiToken) // "7301570305:AAF47O-AtTFC39yu9C5XXxsUqM74xRvXgDE"
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	log.Printf("Авторизован на %s", bot.Self.UserName)

	db, err := gorm.Open(sqlite.Open(config.Database), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	for _, s := range []any{models.User{}, models.Snapshot{}} {
		log.Printf("миграция: %T", s)
		err = db.AutoMigrate(s)
		if err != nil {
			log.Panic(err)
		}
		log.Println("мигрировано")
	}
	err = db.Exec("update users set busy = false").Error
	if err != nil {
		log.Panic(err)
	}

	updHandler := updatehandler.New(config, bot, db)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		go updHandler.Process(&update)
	}
}
