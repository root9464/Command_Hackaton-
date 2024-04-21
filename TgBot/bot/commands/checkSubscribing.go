package commands

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (hb *HomeworkBot) checkSubscription(bot *tgbotapi.BotAPI, chatID int64, channelID int64) {
	req := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{ChatID: channelID, UserID: chatID},
	}

	resp, err := bot.GetChatMember(req)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.Status == "administrator" || resp.Status == "creator" || resp.Status == "member" {
		log.Println("ok")
	} else {
		log.Println("no")
	}
}
