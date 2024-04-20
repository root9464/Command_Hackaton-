package commands

import (
	"log"
	"root/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type HomeworkBot struct {
	bot *tgbotapi.BotAPI
	db  *gorm.DB

	state struct {
		Iter1  bool
		Iter2  bool
		Iter3  bool
		Image1 []tgbotapi.PhotoSize
		Image2 []tgbotapi.PhotoSize
	}
}

func NewHomeworkBot(token string) (*HomeworkBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	methodDB := database.MethodDB()
	db, err := methodDB.Connect()
	if err != nil {
		return nil, err
	}
	return &HomeworkBot{
		bot: bot,
		db:  db,
	}, nil
}

func (hb *HomeworkBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := hb.bot.GetUpdatesChan(u)
	commandMessageSent := false

	for update := range updates {
		if update.Message == nil || update.Message.From.ID == hb.bot.Self.ID {
			continue
		}

		if !commandMessageSent {

			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Выберите шаблон"),
				),
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите команду:")
			msg.ReplyMarkup = keyboard

			if _, err := hb.bot.Send(msg); err != nil {
				log.Println("ошибка: не удалось отправить сообщение. \n", err)
			}

			commandMessageSent = true
		}

		hb.processesBotUpdate(update)

	}
}

func (hb *HomeworkBot) processesBotUpdate(update tgbotapi.Update) {
	hb.SwapFace(update)
}
