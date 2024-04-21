package commands

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (hb *HomeworkBot) HandleCreatePack(updates <-chan tgbotapi.Update, msgChatID int64) {
	for i := 0; i < 3; i++ {
		update := <-updates
		if i == 0 && hb.state.SwapPhoto1 == nil && hb.state.SwapPhoto2 == nil {
			msg := tgbotapi.NewMessage(msgChatID, "ожидаю фото...")
			if _, err := hb.bot.Send(msg); err != nil {
				log.Panic(err)
			}
		} else if i == 1 && hb.state.SwapPhoto1 == nil {
			msg := tgbotapi.NewMessage(msgChatID, "фото 1 получено")
			if _, err := hb.bot.Send(msg); err != nil {
				log.Panic(err)
			}
			hb.state.SwapPhoto1 = &update.Message.Photo
		} else if i == 2 && hb.state.SwapPhoto1 != nil && hb.state.SwapPhoto2 == nil {
			msg := tgbotapi.NewMessage(msgChatID, "фото 2 получено")
			if _, err := hb.bot.Send(msg); err != nil {
				log.Panic(err)
			}
			array, _ := ReadDir("./tmp")
			log.Fatal(array)
		}

		if hb.state.SwapPhoto1 != nil && hb.state.SwapPhoto2 != nil {
			log.Fatal("все фото получены")
		}
	}
}

func ReadDir(path string) ([]string, error) {
	array := []string{}

	// Открываем текущую директорию
	dir, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer dir.Close()

	// Получаем список файлов и папок
	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println(err)
	}

	// Выводим имена файлов и папок
	for _, file := range files {
		array = append(array, file.Name())
	}
	return array, nil
}
