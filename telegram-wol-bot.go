package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Token string
	ChatID float64
	Computers []Computers
}

type Computers struct {
	Name string
	Mac string
}

func fatalError(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	var opts struct {
		ConfigPath string `short:"c" long:"config" description:"Path of the config file" required:"true"`
	}

	_, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	configPath, err := filepath.Abs(opts.ConfigPath)
	fatalError(err)

	jsonFile, err := os.Open(configPath)
	if err != nil {
		fmt.Printf("Error opening JSON file: %s\n", err.Error())
		os.Exit(1)
	}
	defer jsonFile.Close()

	content, err := ioutil.ReadAll(jsonFile)
	fatalError(err)

	var config Config
	err = json.Unmarshal(content, &config)
	fatalError(err)

	runBot(config)
}

func runBot(config Config) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Chat.ID != int64(config.ChatID) {
			continue
		}

		if update.Message.Text[0:5] == "/help" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, `<code>/boot &lt;machine&gt;</code> - Boot the computer user requested
<code>/list</code> - Show a list of computers`)
			msg.ParseMode = "HTML"
			_, _ = bot.Send(msg)
		}

		if update.Message.Text[0:5] == "/list" {
			multi := ""
			if len(config.Computers) == 1 {
				multi = "s"
			}
			var msgText strings.Builder
			msgText.WriteString(fmt.Sprintf("<b>%d computer%s may be waked:</b>\n\n", len(config.Computers), multi))
			for _, e := range config.Computers {
				msgText.WriteString(fmt.Sprintf("<code>%s</code>\n", e.Name))
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText.String())
			msg.ParseMode = "HTML"
			_, _ = bot.Send(msg)
		}

		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID
		// _, _ = bot.Send(msg)
	}
}
