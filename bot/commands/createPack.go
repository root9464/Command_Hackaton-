package commands

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (hb *HomeworkBot) HandleCreatePack(updates <-chan tgbotapi.Update, msgChatID int64) {
	for i := 0; i < 2; i++ {
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
			array, _ := ReadDir("./tmp")
			hb.state.SwapPhoto3 = array
		}

		if hb.state.SwapPhoto1 != nil && hb.state.SwapPhoto3 != nil {
			img1 := *hb.state.SwapPhoto1
			img2 := hb.state.SwapPhoto3
			hb.SwapLocal(img1, img2)
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

func (hb *HomeworkBot) SwapLocal(img1 []tgbotapi.PhotoSize, img2 []string) {
	// Скачиваем первое изображение
	file1, err := hb.bot.GetFile(tgbotapi.FileConfig{FileID: img1[len(img1)-1].FileID})
	if err != nil {
		log.Println("Ошибка: не удалось скачать файл. \n", err)
		return
	}

	// Для каждого изображения из img2
	for index, imgPath := range img2 {
		// Открываем файл в папке ./tmp
		file2, err := os.Open("./tmp/" + imgPath)
		if err != nil {
			log.Println("Ошибка: не удалось открыть файл. \n", err)
			continue
		}
		defer file2.Close()

		// Создаем POST-запрос на сервер
		req, err := http.NewRequest("POST", "http://localhost:5000/swap", nil)
		if err != nil {
			log.Println("Ошибка: не удалось создать запрос. \n", err)
			continue
		}

		// Создаем форму для запроса
		form := &bytes.Buffer{}
		writer := multipart.NewWriter(form)

		// Добавляем первое изображение в форму
		fw1, err := writer.CreateFormFile("image1", "image1.jpg")
		if err != nil {
			log.Println("Ошибка: не удалось создать поле формы. \n", err)
			continue
		}
		// Скачиваем первое изображение с сервера Telegram
		resp, err := http.Get("https://api.telegram.org/file/bot" + hb.bot.Token + "/" + file1.FilePath)
		if err != nil {
			log.Println("Ошибка: не удалось скачать файл. \n", err)
			continue
		}
		defer resp.Body.Close()
		// Копируем содержимое файла в поле формы
		_, err = io.Copy(fw1, resp.Body)
		if err != nil {
			log.Println("Ошибка: не удалось копировать содержимое файлов. \n", err)
			continue
		}
		// Добавляем второе изображение в форму
		fw2, err := writer.CreateFormFile("image2", "image2.jpg")
		if err != nil {
			log.Println("Ошибка: не удалось создать поле формы. \n", err)
			continue
		}
		// Копируем содержимое файла в поле формы
		_, err = io.Copy(fw2, file2)
		if err != nil {
			log.Println("Ошибка: не удалось копировать содержимое файлов. \n", err)
			continue
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
			continue
		}
		if resp.StatusCode == 200 {
			CreateDir(resp, index)
		}

		defer resp.Body.Close()

		log.Print("\033[31mRed\033[0m", index)

	}
}

func CreateDir(resp *http.Response, index int) {
	err := os.MkdirAll("uploads", 0755)
	if err != nil {
		log.Println("Ошибка: не удалось создать папку uploads. \n", err)
		return
	}

	// Сохраняем ответ сервера в файл с уникальным именем
	file, err := os.Create(fmt.Sprintf("uploads/image_%d.png", index))
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
