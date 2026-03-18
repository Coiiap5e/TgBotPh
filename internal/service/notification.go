package service

import "github.com/Coiiap5e/TgBotPh/internal/model"

type Notifier interface {
	Notify(shoot model.Shoot) error
	NotifyMessage(message string) error
}
