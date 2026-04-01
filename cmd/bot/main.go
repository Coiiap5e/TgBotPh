package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Coiiap5e/TgBotPh/internal/adapter/kafka"
	"github.com/Coiiap5e/TgBotPh/internal/adapter/notifier"
	"github.com/Coiiap5e/TgBotPh/internal/config"
	"github.com/Coiiap5e/TgBotPh/internal/logs"
	"github.com/Coiiap5e/TgBotPh/internal/service"
)

func main() {
	logger, cleanup := logs.InitLogger()
	defer cleanup()

	slog.SetDefault(logger)

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Default().Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger.Info("Starting Kafka Consumer service...")

	// Используем StdoutNotifier для временного тестирования
	// tgNotifier, err := notifier.NewTelegramNotifier(cfg.Telegram.BotToken, cfg.Telegram.ChannelID)
	// if err != nil {
	// 	logger.Error("Failed to initialize Telegram notifier", "error", err)
	// 	os.Exit(1)
	// }
	// logger.Info("Telegram notifier initialized successfully")

	var notificationService service.Notifier = notifier.NewStdoutNotifier()
	logger.Info("Stdout notifier initialized successfully for temporary Kafka consumer testing.")

	kafkaConsumer := kafka.NewKafkaConsumer(
		kafka.ConsumerConfig{
			BrokerURLs: cfg.Kafka.BrokerURLs,
			Topic:      cfg.Kafka.Topic,
			GroupID:    cfg.Kafka.GroupID,
		},
		notificationService,
		logger,
	)
	logger.Info("Kafka consumer initialized successfully")

	ctx, cancel := context.WithCancel(context.Background())

	go kafkaConsumer.ConsumeMessages(ctx)

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down service...")

	// Отмена контекста для остановки consumer'а
	cancel()

	// Закрытие Kafka reader
	if err := kafkaConsumer.Close(); err != nil {
		logger.Error("Error closing Kafka consumer", "error", err)
	} else {
		logger.Info("Kafka consumer closed successfully")
	}

	logger.Info("Service stopped gracefully")
}
