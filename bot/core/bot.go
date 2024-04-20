package core

import (
	"log"
	"os"

	"root/bot/commands"

	"github.com/joho/godotenv"
)

func Core() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("ошибка прочтения переменных окружения. \n", err)
	}
	token := os.Getenv("TOKEN")

	bot, err := commands.NewHomeworkBot(token)
	if err != nil {
		log.Fatal("епта не работает")
	}

	bot.Start()
}
