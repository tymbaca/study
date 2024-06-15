package main

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
	"gopkg.in/telebot.v3"
	tele "gopkg.in/telebot.v3"
)

func main() {
	log.SetLevel(log.DebugLevel)

	token := os.Getenv("TOKEN")
	log.Debugf("token: %s", token)
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// menu := &telebot.ReplyMarkup{}
	// btnBuy := menu.Text("🚀 buy")
	// menu.Reply(menu.Row(btnBuy))

	// b.Handle("/start", func(c telebot.Context) error {
	// 	return c.Send("Hello!", menu)
	// })

	// b.Handle(&btnBuy, func(c telebot.Context) error {
	// 	c.Send("kill yourself")
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return nil
	// })

	var (
		// Universal markup builders.
		menu     = &tele.ReplyMarkup{ResizeKeyboard: true}
		selector = &tele.ReplyMarkup{}

		// Reply buttons.
		btnHelp     = menu.Text("ℹ Help")
		btnSettings = menu.Text("⚙ Settings")

		// Inline buttons.
		//
		// Pressing it will cause the client to
		// send the bot a callback.
		//
		// Make sure Unique stays unique as per button kind
		// since it's required for callback routing to work.
		//
		btnPrev = selector.Data("⬅", "prev")
		btnNext = selector.Data("➡", "next")
	)

	menu.Reply(
		menu.Row(btnHelp),
		menu.Row(btnSettings),
	)
	selector.Inline(
		selector.Row(btnPrev, btnNext),
	)

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Hello!", menu)
	})

	// On reply button pressed (message)
	b.Handle(&btnHelp, func(c tele.Context) error {
		return c.Edit("Here is some help: ...")
	})

	// On inline button pressed (callback)
	b.Handle(&btnPrev, func(c tele.Context) error {
		return c.Respond()
	})

	log.Info("waiting for connections")
	b.Start()
}
