package commands

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HomeworkBot struct {
	bot *tgbotapi.BotAPI

	state struct {
		SwapPhoto1 *[]tgbotapi.PhotoSize
		SwapPhoto2 *[]tgbotapi.PhotoSize
		SwapPhoto3 []string
		Iteration  bool
	}
}

func NewHomeworkBot(token string) (*HomeworkBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	return &HomeworkBot{
		bot: bot,
	}, nil
}

func (hb *HomeworkBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := hb.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		msgChatID := update.Message.Chat.ID

		var message1, message2 string
		switch update.Message.Text {
		case "/test":
			for i := 0; i < 2; i++ {
				update := <-updates

				if i == 0 {
					msg := tgbotapi.NewMessage(msgChatID, "1")
					if _, err := hb.bot.Send(msg); err != nil {
						log.Panic(err)
					}
					message1 = update.Message.Text
				} else if i == 1 {
					msg := tgbotapi.NewMessage(msgChatID, "2")
					if _, err := hb.bot.Send(msg); err != nil {
						log.Panic(err)
					}
					message2 = update.Message.Text
				}

				if message1 != "" && message2 != "" {
					log.Fatal("ok")
				}
			}

		case "/swap":
			hb.HandleSwap(updates, msgChatID)
		case "/hello":
			msg := tgbotapi.NewMessage(msgChatID, "Привет")
			_, err := hb.bot.Send(msg)
			if err != nil {
				log.Panic(err)
			}

		case "/check":
			chatID := update.Message.Chat.ID // ID канала
			userID := update.Message.From.ID // ID юзера
			hb.checkSubscription(hb.bot, chatID, userID)

		case "/pack":
			hb.HandleCreatePack(updates, msgChatID)

		case "/help":
			msg := tgbotapi.NewMessage(msgChatID, "Вставьте ссылку в браузер: http://localhost:3000/chat")
			if _, err := hb.bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}

	}

}
