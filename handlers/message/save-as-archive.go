package messagehandler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"os/exec"
	"save-site-as-pdf-bot/models"
	"strconv"
	"time"
)

func (h *MessageHandler) saveAsArchive(update *tgbotapi.Update, initiator *models.User) {
	startTime := time.Now()
	sentFrom := update.SentFrom()
	url := update.Message.Text
	log.Println(sentFrom, url)
	outputDir := fmt.Sprintf("%d", sentFrom.ID)
	outputArchive := fmt.Sprintf("%d.tar.gz", sentFrom.ID)

	// Создание временной директории для сохранения сайта
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка создания временной директории"))
		return
	}
	defer os.RemoveAll(outputDir)

	// Вызов wget для сохранения сайта
	h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Сохраняю сайт в архив..."))
	cmd := exec.Command("wget", "--limit-rate=10m", "--no-clobber", "--convert-links", "--wait", "0.1", "-r", "-p", "-E", "-e", "robots=off", "--level", "1", "-U", "mozilla", url, "-P", outputDir)
	err = cmd.Run()
	if err != nil {
		h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка сохранения сайта"))
		return
	}

	// Упаковка сохраненного сайта в архив
	cmd = exec.Command("tar", "-czf", outputArchive, outputDir)
	err = cmd.Run()
	if err != nil {
		h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка упаковки сайта в архив"))
		return
	}

	// Открытие архива
	file, err := os.Open(outputArchive)
	if err != nil {
		h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка открытия архива"))
		return
	}
	defer file.Close()

	// Отправка архива пользователю
	doc := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FilePath(outputArchive))
	msg, err := h.Bot.Send(doc)
	if err != nil {
		h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка выгрузки архива"))
		return
	}

	// Сохраняем снапшот в базу
	h.DB.Create(&models.Snapshot{
		Initiator:    initiator.ID,
		Site:         url,
		ResultFileID: strconv.Itoa(msg.MessageID),
	})

	// Удаление временного архива
	// err = os.Remove(outputArchive)
	// if err != nil {
	// h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка удаления временного архива"))
	// }

	h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Готово за %.2f секунд", time.Since(startTime).Seconds())))
}
