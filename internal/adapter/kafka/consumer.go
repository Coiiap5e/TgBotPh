package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Coiiap5e/TgBotPh/internal/adapter/notifier"
	"github.com/Coiiap5e/TgBotPh/internal/model"
	kafkago "github.com/segmentio/kafka-go"
)

type ConsumerConfig struct {
	BrokerURLs []string
	Topic      string
	GroupID    string
}

type KafkaConsumer struct {
	reader   *kafkago.Reader
	logger   *slog.Logger
	notifier *notifier.TelegramNotifier
}

func NewKafkaConsumer(cfg ConsumerConfig, notifier *notifier.TelegramNotifier, logger *slog.Logger) *KafkaConsumer {
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:        cfg.BrokerURLs,
		Topic:          cfg.Topic,
		GroupID:        cfg.GroupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		Logger:         kafkago.LoggerFunc(logger.Debug),
		ErrorLogger:    kafkago.LoggerFunc(logger.Error),
	})

	return &KafkaConsumer{
		reader:   r,
		logger:   logger,
		notifier: notifier,
	}
}

func (c *KafkaConsumer) ConsumeMessages(ctx context.Context) {
	c.logger.Info("Starting Kafka consumer", "topic", c.reader.Config().Topic, "groupID", c.reader.Config().GroupID)
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Kafka consumer stopped due to context cancellation")
			return
		default:
			m, err := c.reader.FetchMessage(ctx)
			if err != nil {
				c.logger.Error("Error fetching message from Kafka", "error", err)
				if ctx.Err() != nil { // Check if the error is due to context cancellation
					return
				}
				time.Sleep(time.Second)
				continue
			}

			c.logger.Info("Received message from Kafka", "topic", m.Topic, "partition", m.Partition, "offset", m.Offset, "key", string(m.Key))

			// Обработка сообщения в зависимости от ключа
			switch string(m.Key) {
			case "shoot_notification":
				var shoot model.Shoot
				if err := json.Unmarshal(m.Value, &shoot); err != nil {
					c.logger.Error("Error unmarshalling Kafka message to Shoot", "error", err, "value", string(m.Value))
					if commitErr := c.reader.CommitMessages(ctx, m); commitErr != nil {
						c.logger.Error("Error committing message after unmarshalling failure", "error", commitErr)
					}
					continue
				}

				if err := c.notifier.Notify(shoot); err != nil {
					c.logger.Error("Error sending Telegram notification for shoot", "error", err, "shoot_id", shoot.Id)
				} else {
					c.logger.Info("Telegram notification sent successfully for shoot", "shoot_id", shoot.Id)
				}

			case "general_message":
				message := string(m.Value)
				if err := c.notifier.NotifyMessage(message); err != nil {
					c.logger.Error("Error sending Telegram message notification", "error", err, "message", message)
				} else {
					c.logger.Info("Telegram message notification sent successfully", "message", message)
				}

			default:
				c.logger.Warn("Received Kafka message with unknown key", "key", string(m.Key), "value", string(m.Value))
			}

			// Подтверждение обработки сообщения
			if err := c.reader.CommitMessages(ctx, m); err != nil {
				c.logger.Error("Error committing message", "error", err)
			}
		}
	}
}

func (c *KafkaConsumer) Close() error {
	c.logger.Info("Closing Kafka reader")
	return c.reader.Close()
}
