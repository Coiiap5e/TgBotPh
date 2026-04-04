package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/Coiiap5e/TgBotPh/internal/errors"
	"github.com/Coiiap5e/TgBotPh/internal/model"
	"github.com/Coiiap5e/TgBotPh/internal/service" // Добавлен импорт service
	kafkago "github.com/segmentio/kafka-go"
)

type ConsumerConfig struct {
	BrokerURLs []string
	Topic      string
	GroupID    string
}

type KafkaConsumer struct {
	reader                  *kafkago.Reader
	logger                  *slog.Logger
	notifier                service.Notifier // Заглушка
	meter                   metric.Meter
	kafkaGetMessagesCounter metric.Int64Counter
}

func NewKafkaConsumer(cfg ConsumerConfig, notifier service.Notifier, logger *slog.Logger, meter metric.Meter) (*KafkaConsumer, error) {
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

	consumer := &KafkaConsumer{
		reader:   r,
		logger:   logger,
		notifier: notifier,
		meter:    meter,
	}

	kafkaGetMessagesCounter, err := meter.Int64Counter(
		"kafka_get_messages_total",
		metric.WithDescription("Total number of get messeges from Kafka"),
		metric.WithUnit("1"),
	)
	if err != nil {
		logger.Error("failed to create Kafka get messages counter", "error", err)
		return nil, errors.Wrap(err, errors.ErrCodeKafkaConsume, "failed to create Kafka get messages counter")
	}

	consumer.kafkaGetMessagesCounter = kafkaGetMessagesCounter

	return consumer, nil
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
					c.logger.Error("Error sending notification for shoot", "error", err, "shoot_id", shoot.Id)
				} else {
					c.logger.Info("Notification sent successfully for shoot", "shoot_id", shoot.Id)
				}

			case "general_message":
				message := string(m.Value)
				if err := c.notifier.NotifyMessage(message); err != nil {
					c.logger.Error("Error sending message notification", "error", err, "message", message)
				} else {
					c.logger.Info("Message notification sent successfully", "message", message)
				}

			default:
				c.logger.Warn("Received Kafka message with unknown key", "key", string(m.Key), "value", string(m.Value))
			}

			// Подтверждение обработки сообщения
			if err := c.reader.CommitMessages(ctx, m); err != nil {
				c.logger.Error("Error committing message", "error", err)
			}
			c.kafkaGetMessagesCounter.Add(ctx, 1, metric.WithAttributes(
				attribute.String("topic", m.Topic),
				attribute.String("message_type", string(m.Key)),
			))

		}
	}
}

func (c *KafkaConsumer) Close() error {
	c.logger.Info("Closing Kafka reader")
	return c.reader.Close()
}
