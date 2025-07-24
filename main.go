package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Println(err.Error())
	}
	certFile := "/etc/letsencrypt/live/realtor-bot.mooo.com/fullchain.pem"
	certData, err := os.ReadFile(certFile)
	if err != nil {
		log.Println(err.Error())
	}

	// 2. Создаем RequestFileData для сертификата
	certReader := tgbotapi.FileBytes{
		Name:  "cert.pem",
		Bytes: certData,
	}

	// Настройка вебхука
	webhookURL := "https://realtor-bot.mooo.com:8443/" + bot.Token
	wh, err := tgbotapi.NewWebhookWithCert(webhookURL, certReader)
	if err != nil {
		log.Println(err.Error())
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Println(err.Error())
	}

	// Получаем обновления через вебхук
	updates := bot.ListenForWebhook("/" + bot.Token)
	go func() {
		err = http.ListenAndServeTLS(":8443",
			"/etc/letsencrypt/live/realtor-bot.mooo.com/fullchain.pem",
			"/etc/letsencrypt/live/realtor-bot.mooo.com/privkey.pem",
			nil)
		if err != nil {
			log.Println(err.Error())
		}
	}()

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Echo: "+update.Message.Text)
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
