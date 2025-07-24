package main

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("7948439516:AAGHu5ITqmKNBnR_crFBXSyju_MczwTOGsQ")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	// Установка webhook
	wh, err := tgbotapi.NewWebhook("https://realtor-bot.mooo.com:8443/webhook")
	if err != nil {
		log.Fatal(err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}

	// HTTP сервер для обработки webhook
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		update, err := bot.HandleUpdate(r)
		if err != nil {
			log.Println(err)
			return
		}

		// Обработка сообщения
		if update.Message != nil {
			handleMessage(bot, update.Message)
		}
	})

	go http.ListenAndServe(":8080", nil)
	log.Println("Server started on :8080")

	// Бесконечный цикл для поддержания работы программы
	select {}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	switch message.Command() {
	case "start":
		msg.Text = "Привет! Я бот с webhook."
	case "help":
		msg.Text = "Доступные команды: /start, /help"
	default:
		msg.Text = "Неизвестная команда"
	}

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}
