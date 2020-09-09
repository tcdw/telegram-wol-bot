package main

import (
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/mdlayher/wol"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Token string `json:"token"`
	ChatID float64 `json:"chatID"`
	Computers []Computer `json:"computers"`
}

type Computer struct {
	Name string `json:"name"`
	Mac string `json:"mac"`
	IP *string `json:"broadcast"`
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
	b, err := tb.NewBot(tb.Settings{
		Token: config.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(m *tb.Message) {
		if !verify(config.ChatID, m) {
			return
		}
		_, _ = b.Send(m.Sender, "This bot is started!")
	})

	b.Handle("/help", func(m *tb.Message) {
		if !verify(config.ChatID, m) {
			return
		}
		_, _ = b.Send(m.Sender, `<code>/boot &lt;machine&gt;</code> - Boot the computer user requested
<code>/list</code> - Show a list of computers`, &tb.SendOptions{ParseMode: "HTML"})
	})

	b.Handle("/list", func(m *tb.Message) {
		if !verify(config.ChatID, m) {
			return
		}
		multi := "s"
		if len(config.Computers) == 1 {
			multi = ""
		}
		var msgText strings.Builder
		msgText.WriteString(fmt.Sprintf("<b>%d computer%s may be waked:</b>\n\n", len(config.Computers), multi))
		for _, e := range config.Computers {
			msgText.WriteString(fmt.Sprintf("<code>%s</code>\n", e.Name))
		}
		_, _ = b.Send(m.Sender, msgText.String(), &tb.SendOptions{ParseMode: "HTML"})
	})

	b.Handle("/boot", func(m *tb.Message) {
		if !verify(config.ChatID, m) {
			return
		}
		target := m.Payload
		var item *Computer
		for _, e := range config.Computers {
			if e.Name == target {
				item = &e
				break
			}
		}
		if item == nil {
			_, _ = b.Send(m.Sender, "Unknown machine. Use <code>/list</code> to see a list of computers.",
				&tb.SendOptions{ParseMode: "HTML"})
			return
		}
		var ip string = "255.255.255.255:9"
		if item.IP != nil {
			ip = *item.IP
		}
		err = wake(ip, item.Mac, []byte(""))
		result := "Boot command sent successfully."
		if err != nil {
			result = fmt.Sprintf("Unable to send boot command: %s", err.Error())
		}
		_, _ = b.Send(m.Sender, result)
	})

	log.Println("Starting bot service")
	b.Start()
}

func wake(addr string, tar string, password []byte) error {
	target, err := net.ParseMAC(tar)
	c, err := wol.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()
	return c.WakePassword(addr, target, password)
}

func verify(chatID float64, m *tb.Message) bool {
	if int64(chatID) != m.Chat.ID {
		log.Printf("Chat ID %d is not authorized, skipping\n", m.Chat.ID)
		return false
	}
	log.Printf("User %d issued command: %s\n", m.Sender.ID, m.Text)
	return true
}
