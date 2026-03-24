package main

import (
	"fmt"
	"log"

	swissknife "github.com/Sagleft/swiss-knife"
	"gopkg.in/telebot.v4"
)

type Config struct {
	Telegram TelegramConfig
}

type TelegramConfig struct {
	BotToken string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true" json:"token"`
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}

	log.Println("app finished")
}

func run() error {
	var cfg Config
	if err := swissknife.ProcessConfig(&cfg); err != nil {
		return fmt.Errorf("process config: %w", err)
	}

	b, err := telebot.NewBot(telebot.Settings{
		Token: cfg.Telegram.BotToken,
	})
	if err != nil {
		return fmt.Errorf("init bot: %w", err)
	}

	menu := b.NewMarkup()

	row := menu.Row(
		menu.Chat("Выбрать группу",
			&telebot.ReplyRecipient{
				ID:           1,                  // на каждую кнопку необходим свой ID
				Channel:      false,              // запрашиваем не каналы
				RequestTitle: telebot.Flag(true), // просить отображать название чата
				BotRights:    nil,                // можно оставить nil, если не нужен filtro по правам бота
			},
		),
	)
	menu.Reply(row)

	b.Handle("/start", func(ctx telebot.Context) error {
		return ctx.Send(
			fmt.Sprintf("Привет! Твой ID: `%d`", ctx.Sender().ID),
			menu, telebot.ModeMarkdown,
		)
	})

	b.Handle(telebot.OnChatShared, func(ctx telebot.Context) error {
		if ctx.Message() == nil || ctx.Message().ChatShared == nil {
			return ctx.Reply("нет данных")
		}

		return ctx.Send(
			fmt.Sprintf("`%d`", ctx.Message().ChatShared.ChatID),
			menu, telebot.ModeMarkdown,
		)
	})

	log.Println("app started")

	b.Start()
	return nil
}
