package commandhandler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"save-site-as-pdf-bot/config"
	"save-site-as-pdf-bot/models"
)

type CommandHandler struct {
	Bot    *tgbotapi.BotAPI
	DB     *gorm.DB
	Config *config.Config
}

func (h *CommandHandler) Process(update *tgbotapi.Update, user *models.User) {
	sentFrom := update.SentFrom()

	var initiator *models.User
	h.DB.First(&initiator, "tg_id = ?", sentFrom.ID)

	command := update.Message.Command()
	if command == "start" || command == "help" {
		msg := tgbotapi.NewMessage(sentFrom.ID, "Отправь ссылку на ресурс")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("сменить метод"),
			),
		)
		h.Bot.Send(msg)
	}
}
