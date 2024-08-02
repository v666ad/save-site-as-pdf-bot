package updatehandler

import (
	"log"
	"save-site-as-pdf-bot/config"
	"save-site-as-pdf-bot/handlers/message"
	"save-site-as-pdf-bot/handlers/command"
	"save-site-as-pdf-bot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type UpdateHandler struct {
	Bot            *tgbotapi.BotAPI
	DB             *gorm.DB
	CommandHandler *commandhandler.CommandHandler
	MessageHandler *messagehandler.MessageHandler
	Config         *config.Config
}

func New(config *config.Config, bot *tgbotapi.BotAPI, db *gorm.DB) *UpdateHandler {
	h := &UpdateHandler{Bot: bot, DB: db}

	cmdHandler := &commandhandler.CommandHandler{
		Bot: bot,
		DB: db,
		Config: config,
	}
	h.CommandHandler = cmdHandler

	msgHandler := &messagehandler.MessageHandler{
		Bot:    bot,
		DB:     db,
		Config: config,
	}
	h.MessageHandler = msgHandler

	return h

}

func (h *UpdateHandler) Process(update *tgbotapi.Update) {
	sentFrom := update.SentFrom()
	if sentFrom == nil {
		return
	}

	user, valid, err := h.validateUser(sentFrom.ID)
	if err != nil {
		log.Panic(err)
	}
	if !valid {
		return
	}

	if update.Message != nil {
		if update.Message.IsCommand() {
			h.CommandHandler.Process(update, user)
		} else {
			h.MessageHandler.Process(update, user)
		}
	}
}

func (h *UpdateHandler) validateUser(tgID int64) (*models.User, bool, error) {
	user := new(models.User)
	err := h.DB.First(user, "tg_id = ?", tgID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user.TgID = tgID
			h.DB.Create(user)
			return user, true, nil
		} else {
			return nil, false, err
		}
	}
	return user, true, nil
}
