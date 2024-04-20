package commands

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"

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

		img1 := hb.state.Image1
		img2 := hb.state.Image2

		file1, err := hb.bot.GetFile(tgbotapi.FileConfig{FileID: img1[len(img1)-1].FileID})
		if err != nil {
			log.Println("Ошибка: не удалось скачать файл. \n", err)
			return
		}
		file2, err := hb.bot.GetFile(tgbotapi.FileConfig{FileID: img2[len(img2)-1].FileID})
		if err != nil {
			log.Println("Ошибка: не удалось скачать файл. \n", err)
			return
		}

		// Create a new HTTP request
		req, err := http.NewRequest("POST", "http://localhost:5000/swap", nil)
		if err != nil {
			log.Println("Ошибка: не удалось создать запрос. \n", err)
			return
		}

		// Create a new multipart form
		form := &bytes.Buffer{}
		writer := multipart.NewWriter(form)

		// Add img1 to the form
		fw1, err := writer.CreateFormFile("image1", "image1.jpg")
		if err != nil {
			log.Println("Ошибка: не удалось создать поле формы. \n", err)
			return
		}
		resp, err := http.Get("https://api.telegram.org/file/bot" + hb.bot.Token + "/" + file1.FilePath)
		if err != nil {
			log.Println("Ошибка: не удалось скачать файл. \n", err)
			return
		}
		defer resp.Body.Close()
		io.Copy(fw1, resp.Body)

		// Add img2 to the form
		fw2, err := writer.CreateFormFile("image2", "image2.jpg")
		if err != nil {
			log.Println("Ошибка: не удалось создать поле формы. \n", err)
			return
		}
		resp, err = http.Get("https://api.telegram.org/file/bot" + hb.bot.Token + "/" + file2.FilePath)
		if err != nil {
			log.Println("Ошибка: не удалось скачать файл. \n", err)
			return
		}
		defer resp.Body.Close()
		io.Copy(fw2, resp.Body)

		// Close the writer
		writer.Close()

		// Set the request body and headers
		req.Body = io.NopCloser(form)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Send the request
		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			log.Println("Ошибка: не удалось отправить запрос. \n", err)
			return
		}
		defer resp.Body.Close()

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			log.Println("Ошибка: сервер вернул код состояния", resp.StatusCode)
			return
		}

		// If everything is successful, send a success message
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Фото успешно отправлены на сервер!")
		if _, err := hb.bot.Send(msg); err != nil {
			log.Println("Ошибка: не удалось отправить сообщение. \n", err)
		}

	}
}
