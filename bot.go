package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var botUsers = os.Getenv("BOT_USERS")

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}
	users := strings.Split(botUsers, ",")

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			message := update.Message
			atStr := strings.TrimPrefix(update.Message.CommandWithAt(), update.Message.Command())
			if !message.IsCommand() || (atStr != "" && atStr != "@"+bot.Self.UserName) || !in(strconv.Itoa(int(message.From.ID)), users) {
				continue
			}
			switch message.Command() {
			case "status":
				var text string
				text += GetHostInfo() + "\n"
				text += GetCPUPercents() + "\n"
				text += GetMemStats() + "\n"
				text += GetRootUsage() + "\n"
				text += GetNetworkAllStats()
				sendMsg(bot, message, text)
			case "cpu_status":
				sendMsg(bot, message, GetCPUCoresPercents())
			case "parts_status":
				sendMsg(bot, message, GetPartsStats())
			case "network_status":
				sendMsg(bot, message, GetNetworkStats())
			}
		}
	}
}

func sendMsg(bot *tgbotapi.BotAPI, message *tgbotapi.Message, text string) {
	text = "```\n" + text + "\n```"
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyToMessageID = message.MessageID

	_, err := bot.Send(msg)

	if err != nil {
		log.Println(err)
	}
}
