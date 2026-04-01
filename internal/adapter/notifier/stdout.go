package notifier

import (
	"fmt"

	"github.com/Coiiap5e/TgBotPh/internal/model"
)

type StdoutNotifier struct{}

func NewStdoutNotifier() *StdoutNotifier {
	return &StdoutNotifier{}
}

func (n *StdoutNotifier) Notify(shoot model.Shoot) error {
	fmt.Printf("--- Kafka Message (Shoot) ---\n")
	fmt.Printf("ID: %d\n", shoot.Id)
	fmt.Printf("Date: %s\n", shoot.ShootDate.Format("2006-01-02"))
	fmt.Printf("Time: %s\n", shoot.StartTime.Format("15:04"))
	fmt.Printf("Type: %s\n", shoot.ShootType)
	fmt.Printf("Location: %s\n", shoot.ShootLocation)
	fmt.Printf("Price: %d RUB (%.2f USD)\n", shoot.ShootPrice, shoot.PriceUSD)
	if len(shoot.Clients) > 0 {
		for _, client := range shoot.Clients {
			if client.IsMainClient {
				fmt.Printf("Main Client: %s %s, Phone: %s\n", client.FirstName, client.LastName, client.Phone)
			}
		}
	} else {
		fmt.Printf("No main client assigned.\n")
	}
	fmt.Printf("-----------------------------\n")
	return nil
}

func (n *StdoutNotifier) NotifyMessage(message string) error {
	fmt.Printf("--- Kafka Message (Raw) ---\n")
	fmt.Printf("%s\n", message)
	fmt.Printf("---------------------------\n")
	return nil
}
