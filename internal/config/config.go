package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Coiiap5e/TgBotPh/internal/errors"
	"github.com/joho/godotenv"
)

const (
	EnvTelegramBotToken  = "TELEGRAM_BOT_TOKEN"
	EnvTelegramChannelID = "TELEGRAM_CHANNEL_ID"
	EnvKafkaBrokerURLs   = "KAFKA_BROKER_URLS"
	EnvKafkaTopic        = "KAFKA_TOPIC"
	EnvKafkaGroupID      = "KAFKA_GROUP_ID"
)

type Config struct {
	Telegram TelegramConfig
	Kafka    KafkaConfig
	Metrics  MetricsConfig
}

type TelegramConfig struct {
	BotToken  string
	ChannelID string
}

type MetricsConfig struct {
	FilePath       string
	ExportInterval time.Duration
}

type KafkaConfig struct {
	BrokerURLs []string
	Topic      string
	GroupID    string
}

// Load - загружает конфигурацию приложения из переменных окружения.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeConfig, "error loading .env file")
	}

	//telegramConfig, err := loadTelegramConfig()
	//if err != nil {
	//	return nil, errors.Wrap(err, errors.ErrCodeConfig, "failed to load Telegram config")
	//}

	kafkaConfig, err := loadKafkaConfig()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeConfig, "failed to load Kafka config")
	}

	metricsConfig, err := loadMetricsConfig()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeConfig, "failed to load metrics config")
	}

	return &Config{
		//Telegram: *telegramConfig,
		Kafka:   *kafkaConfig,
		Metrics: *metricsConfig,
	}, nil
}

func loadMetricsConfig() (*MetricsConfig, error) {
	filePath := getEnv("METRICS_FILE_PATH", "metrics.log")
	exportIntervalStr := getEnv("METRICS_EXPORT_INTERVAL_SECONDS", "10")

	exportIntervalSeconds, err := strconv.Atoi(exportIntervalStr)
	if err != nil {
		return nil, errors.New(errors.ErrCodeConfig, "invalid METRICS_EXPORT_INTERVAL_SECONDS")
	}

	return &MetricsConfig{
		FilePath:       filePath,
		ExportInterval: time.Duration(exportIntervalSeconds) * time.Second,
	}, nil
}

func loadTelegramConfig() (*TelegramConfig, error) {
	telegramBotToken := os.Getenv(EnvTelegramBotToken)
	if telegramBotToken == "" {
		return nil, errors.New(errors.ErrCodeConfig, "Telegram bot token is not set")
	}
	telegramChannelID := os.Getenv(EnvTelegramChannelID)
	if telegramChannelID == "" {
		return nil, errors.New(errors.ErrCodeConfig, "Telegram channel ID is not set")
	}

	return &TelegramConfig{
		BotToken:  telegramBotToken,
		ChannelID: telegramChannelID,
	}, nil
}

func loadKafkaConfig() (*KafkaConfig, error) {
	kafkaBrokerURLsStr := os.Getenv(EnvKafkaBrokerURLs)
	if kafkaBrokerURLsStr == "" {
		return nil, errors.New(errors.ErrCodeConfig, "Kafka broker URLs are not set")
	}
	kafkaBrokerURLs := strings.Split(kafkaBrokerURLsStr, ",")
	if len(kafkaBrokerURLs) == 0 {
		return nil, errors.New(errors.ErrCodeConfig, "No Kafka broker URLs provided")
	}

	kafkaTopic := os.Getenv(EnvKafkaTopic)
	if kafkaTopic == "" {
		return nil, errors.New(errors.ErrCodeConfig, "Kafka topic is not set")
	}

	kafkaGroupID := getEnv(EnvKafkaGroupID, "telegram-bot-consumer-group") // Group ID по умолчанию

	return &KafkaConfig{
		BrokerURLs: kafkaBrokerURLs,
		Topic:      kafkaTopic,
		GroupID:    kafkaGroupID,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
