package notifier

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Coiiap5e/TgBotPh/internal/errors"
	"github.com/Coiiap5e/TgBotPh/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

type TelegramNotifier struct {
	bot       *tgbotapi.BotAPI
	channelID int64
}

func NewTelegramNotifier(token string, channelID string) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeTelegramBotInit, "failed to initialize telegram bot")
	}
	id, err := strconv.ParseInt(channelID, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInvalidInput, "invalid telegram channel id")
	}

	return &TelegramNotifier{
		bot:       bot,
		channelID: id,
	}, nil
}

func (n *TelegramNotifier) Notify(shoot model.Shoot) error {
	var sb strings.Builder

	sb.WriteString("<b>Upcoming Shoot!</b>")
	sb.WriteString(fmt.Sprintf("<b>ID:</b> %d", shoot.Id))
	sb.WriteString(fmt.Sprintf("<b>Date:</b> %s", shoot.ShootDate.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("<b>Time:</b> %s", shoot.StartTime.Format("15:04")))
	sb.WriteString(fmt.Sprintf("<b>Type:</b> %s", shoot.ShootType))
	sb.WriteString(fmt.Sprintf("<b>Location:</b> %s", shoot.ShootLocation))
	sb.WriteString(fmt.Sprintf("<b>Price:</b> %d RUB (%.2f USD)", shoot.ShootPrice, shoot.PriceUSD))

	mainClient, ok := lo.Find(shoot.Clients, func(c model.ShootClientInfo) bool {
		return c.IsMainClient
	})

	if ok {
		sb.WriteString("<b>Main Client:</b>")
		sb.WriteString(fmt.Sprintf("  <b>Name:</b> %s %s", mainClient.FirstName, mainClient.LastName))
		sb.WriteString(fmt.Sprintf("  <b>Phone:</b> %s", mainClient.Phone))
	} else {
		sb.WriteString("<b>No main client assigned.</b>")
	}

	msg := tgbotapi.NewMessage(n.channelID, sb.String())
	msg.ParseMode = tgbotapi.ModeHTML
	if _, err := n.bot.Send(msg); err != nil {
		return errors.Wrap(err, errors.ErrCodeTelegramBotSend, "failed to send message to telegram")
	}
	return nil
}

func (n *TelegramNotifier) NotifyMessage(message string) error {
	msg := tgbotapi.NewMessage(n.channelID, message)
	msg.ParseMode = tgbotapi.ModeHTML
	if _, err := n.bot.Send(msg); err != nil {
		return errors.Wrap(err, errors.ErrCodeTelegramBotSend, "failed to send message to telegram")
	}
	return nil
}
