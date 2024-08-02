package messagehandler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"save-site-as-pdf-bot/config"
	"save-site-as-pdf-bot/models"
	"time"
)

const UserAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0"

type MessageHandler struct {
	Bot    *tgbotapi.BotAPI
	DB     *gorm.DB
	Config *config.Config
}

func (h *MessageHandler) Process(update *tgbotapi.Update, user *models.User) {
	sentFrom := update.SentFrom()

	var initiator *models.User
	h.DB.First(&initiator, "tg_id = ?", sentFrom.ID)

	if update.Message.Text == "сменить метод" {
		if initiator.SaveMode == "pdf" {
			initiator.SaveMode = "archive"
		} else {
			initiator.SaveMode = "pdf"
		}
		h.DB.Model(initiator).Update("save_mode", initiator.SaveMode)
		h.Bot.Send(tgbotapi.NewMessage(initiator.TgID, "Метод изменён на "+initiator.SaveMode))
		return
	}

	if initiator.Busy {
		h.Bot.Send(tgbotapi.NewMessage(initiator.TgID, "Ты уже что-то скачиваешь. Ожидай."))
		return
	}

	timeSinceLastSnapshot := time.Since(initiator.LastSnapshotTime)
	if timeSinceLastSnapshot < h.Config.DelayBetweenSnapshots {
		h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Скачивать страницу можно раз в %.2f секунд(ы). Осталось секунд %.2f",
			h.Config.DelayBetweenSnapshots.Seconds(),
			(h.Config.DelayBetweenSnapshots-timeSinceLastSnapshot).Seconds())))
		return
	}

	h.DB.Model(initiator).Update("last_snapshot_time", time.Now()).Update("busy", true)
	defer func(initiator *models.User) {
		h.DB.Model(initiator).Update("busy", false)
	}(initiator)

	if initiator.SaveMode == "pdf" {
		h.saveAsPDF(update, initiator)
	} else {
		h.saveAsArchive(update, initiator)
	}
}
