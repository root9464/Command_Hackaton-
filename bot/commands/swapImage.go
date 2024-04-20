package commands

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (hb *HomeworkBot) SwapFace(update tgbotapi.Update) {

	if update.Message.Text == "/swap" {
		hb.state.Image1 = nil
		hb.state.Image2 = nil

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите свои фото:")
		if _, err := hb.bot.Send(msg); err != nil {
			log.Println("Ошибка: не удалось отправить сообщение. \n", err)
		}

		hb.state.Image1 = update.Message.Photo
	} else if hb.state.Image1 != nil && hb.state.Image2 == nil {
		hb.state.Image2 = update.Message.Photo

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите фото для слияния:")
		if _, err := hb.bot.Send(msg); err != nil {
			log.Println("Ошибка: не удалось отправить сообщение. \n", err)
		}

	} else if hb.state.Image1 != nil && hb.state.Image2 != nil {
		//отправить hb.state.Image1  и hb.state.Image2 на эндпоинт
	}
}
