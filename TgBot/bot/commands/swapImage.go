package commands

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (hb *HomeworkBot) HandleSwap(updates <-chan tgbotapi.Update, msgChatID int64) {
	for i := 0; i < 3; i++ {
		update := <-updates

		if i == 0 && hb.state.SwapPhoto1 == nil && hb.state.SwapPhoto2 == nil {
			msg := tgbotapi.NewMessage(msgChatID, "between photos...")
			if _, err := hb.bot.Send(msg); err != nil {
				log.Panic(err)
			}
		} else if i == 1 && hb.state.SwapPhoto1 == nil {
			msg := tgbotapi.NewMessage(msgChatID, "photo 1 received")
			if _, err := hb.bot.Send(msg); err != nil {
				log.Panic(err)
			}
			hb.state.SwapPhoto1 = &update.Message.Photo
		} else if i == 2 && hb.state.SwapPhoto1 != nil && hb.state.SwapPhoto2 == nil {
			msg := tgbotapi.NewMessage(msgChatID, "photo 2 received")
			if _, err := hb.bot.Send(msg); err != nil {
				log.Panic(err)
			}
			hb.state.SwapPhoto2 = &update.Message.Photo
		}

		if hb.state.SwapPhoto1 != nil && hb.state.SwapPhoto2 != nil {
			img1 := *hb.state.SwapPhoto1
			img2 := *hb.state.SwapPhoto2
			hb.Swap(img1, img2)
		}
	}
}

// Swap выполняет обмен изображениями между двумя пользователями.
// Оно скачивает изображения из Telegram и отправляет их на сервер для обмена.
// Затем оно сохраняет обмененное изображение в папку "uploads".
func (hb *HomeworkBot) Swap(img1, img2 []tgbotapi.PhotoSize) {
	// Скачиваем первый изображение
	file1, err := hb.bot.GetFile(tgbotapi.FileConfig{FileID: img1[len(img1)-1].FileID})
	if err != nil {
		log.Println("Ошибка: не удалось скачать файл. \n", err)
		return
	}
	// Скачиваем второй изображение
	file2, err := hb.bot.GetFile(tgbotapi.FileConfig{FileID: img2[len(img2)-1].FileID})
	if err != nil {
		log.Println("Ошибка: не удалось скачать файл. \n", err)
		return
	}

	// Создаем POST-запрос на сервер
	req, err := http.NewRequest("POST", "http://localhost:5000/swap", nil)
	if err != nil {
		log.Println("Ошибка: не удалось создать запрос. \n", err)
		return
	}

	// Создаем форму для запроса
	form := &bytes.Buffer{}
	writer := multipart.NewWriter(form)

	// Добавляем первый изображение в форму
	fw1, err := writer.CreateFormFile("image1", "image1.jpg")
	if err != nil {
		log.Println("Ошибка: не удалось создать поле формы. \n", err)
		return
	}
	// Скачиваем первый изображение с сервера Telegram
	resp, err := http.Get("https://api.telegram.org/file/bot" + hb.bot.Token + "/" + file1.FilePath)
	if err != nil {
		log.Println("Ошибка: не удалось скачать файл. \n", err)
		return
	}
	defer resp.Body.Close()
	// Копируем содержимое файла в поле формы
	_, err = io.Copy(fw1, resp.Body)
	if err != nil {
		log.Println("Ошибка: не удалось копировать содержимое файлов. \n", err)
		return
	}
	// Добавляем второй изображение в форму
	fw2, err := writer.CreateFormFile("image2", "image2.jpg")
	if err != nil {
		log.Println("Ошибка: не удалось создать поле формы. \n", err)
		return
	}
	// Скачиваем второй изображение с сервера Telegram
	resp, err = http.Get("https://api.telegram.org/file/bot" + hb.bot.Token + "/" + file2.FilePath)
	if err != nil {
		log.Println("Ошибка: не удалось скачать файл. \n", err)
		return
	}
	defer resp.Body.Close()
	// Копируем содержимое файла в поле формы
	_, err = io.Copy(fw2, resp.Body)
	if err != nil {
		log.Println("Ошибка: не удалось копировать содержимое файлов. \n", err)
		return
	}
	// Закрываем запись в форму
	writer.Close()

	// Устанавливаем заголовок запроса
	req.Body = io.NopCloser(form)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Отправляем запрос на сервер
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Println("Ошибка: не удалось отправить запрос. \n", err)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус код ответа
	if resp.StatusCode == http.StatusOK {
		// Создаем папку uploads, если она еще не существует
		err := os.MkdirAll("uploads", 0755)
		if err != nil {
			log.Println("Ошибка: не удалось создать папку uploads. \n", err)
			return
		}

		// Сохраняем ответ сервера в файл
		file, err := os.Create("uploads/image.png")
		if err != nil {
			log.Println("Ошибка: не удалось создать файл. \n", err)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			log.Println("Ошибка: не удалось сохранить файл. \n", err)
			return
		}

		log.Println("Файл успешно сохранен в папку uploads.")
	}
}
