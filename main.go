package main

import (
	"./tools"
	"github.com/tidwall/gjson"
	"time"
	"sync"
	"log"
	"strings"
	"github.com/Narrator69/telegram-bot-api"
)

func main () {
	data := tools.GetFileContent("data.json")

	botsWaitGroup := sync.WaitGroup{}

	gjson.Parse(data).ForEach(func(key, value gjson.Result) bool {
		botsWaitGroup.Add(1)
		go iterateOverUsers(&botsWaitGroup, key.String(), value)
		return true // keep iterating
	})

	botsWaitGroup.Wait()

	log.Println("Done")
}

func iterateOverUsers(botsWaitGroup *sync.WaitGroup, token string, users gjson.Result) {
	defer botsWaitGroup.Done()

	messagesWaitGroup := sync.WaitGroup{}

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Printf("%s: " + "Error of token", token)
		log.Println(err)
		return
	}

	bot.Debug = true

	sentMessagesCounter := 0

	users.ForEach(func(uKey, messages gjson.Result) bool {
		ucid := uKey.Int()

		if ucid == 0 {
			log.Printf("%s: %s: " + "Skipped as empty user", token, uKey.String())
			return true // keep iterating
		}

		messages.ForEach(func(mKey, message gjson.Result) bool {
			if sentMessagesCounter > 0 && sentMessagesCounter % 30 == 0 {
				time.Sleep(time.Second)
				log.Printf("%s: " + "Slept for a second", token)
			}

			contentOfMessage := message.String()

			if strings.Contains(contentOfMessage, "file:") {
				path := strings.Split(contentOfMessage, ":")[1]
				contentOfMessage = tools.GetFileContent(path)
			}

			if contentOfMessage == "" {
				log.Printf("%s: %s: %s: " + "Skipped as empty message", token, uKey.String(), message.String())
				return true // keep iterating
			}

			msg := tgbotapi.NewMessage(ucid, contentOfMessage)

			messagesWaitGroup.Add(1)

			go func() {
				bot.Send(msg)
				messagesWaitGroup.Done()
			}()

			sentMessagesCounter++

			return true // keep iterating
		})

		return true // keep iterating
	})

	messagesWaitGroup.Wait()
}