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

func (h *MessageHandler) saveAsPDF(update *tgbotapi.Update, initiator *models.User) {
	startTime := time.Now()
	sentFrom := update.SentFrom()
	url := update.Message.Text
	log.Println(sentFrom, url)
	output := fmt.Sprintf("%d.pdf", sentFrom.ID)

	// Вызов wkhtmltopdf для конвертации URL в PDF
	h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Конвертирую URL в PDF..."))
	cmd := exec.Command("wkhtmltopdf", "--custom-header-propagation", "--custom-header", "User-Agent", UserAgent, url, output)
	err := cmd.Run()
	if err != nil {
		h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка ковертации URL в PDF"))
		return
	}

	// Открытие PDF файла
	file, err := os.Open(output)
	if err != nil {
		h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка открытия PDF файла"))
		return
	}
	defer file.Close()

	// Отправка PDF файла пользователю
	doc := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FilePath(output))
	msg, err := h.Bot.Send(doc)
	if err != nil {
		h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка выгрузки PDF файла"))
		return
	}

	// Сохраняем снапшот в базу
	h.DB.Create(&models.Snapshot{
		Initiator:    initiator.ID,
		Site:         url,
		ResultFileID: strconv.Itoa(msg.MessageID),
	})

	// Удаление временного PDF файла
	err = os.Remove(output)
	if err != nil {
		h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка удаления временного файла"))
	}

	h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Готово за %.2f секунд", time.Since(startTime).Seconds())))
}
