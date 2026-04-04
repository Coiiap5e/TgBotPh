package metrics

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type FileExporter struct {
	filePath string
	logger   *slog.Logger
	mu       sync.Mutex
	file     *os.File
}

func NewFileExporter(filePath string, logger *slog.Logger) (*FileExporter, error) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory for metrics file: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open metrics file: %w", err)
	}

	return &FileExporter{
		filePath: filePath,
		logger:   logger,
		file:     file,
	}, nil
}

func (e *FileExporter) Temporality(kind sdkmetric.InstrumentKind) metricdata.Temporality {
	return metricdata.CumulativeTemporality
}

func (e *FileExporter) Aggregation(kind sdkmetric.InstrumentKind) sdkmetric.Aggregation {
	return sdkmetric.AggregationSum{}
}

func (e *FileExporter) ForceFlush(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.file != nil {
		return e.file.Sync()
	}
	return nil
}

func (e *FileExporter) Export(ctx context.Context, rm *metricdata.ResourceMetrics) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, scopeMetrics := range rm.ScopeMetrics {
		for _, metricData := range scopeMetrics.Metrics {
			if metricData.Name == "kafka_get_messages_total" {
				if sum, ok := metricData.Data.(metricdata.Sum[int64]); ok {
					for _, dp := range sum.DataPoints {
						line := fmt.Sprintf("[%s] kafka_get_messages_total: %d\n", time.Now().Format(time.RFC3339), dp.Value)
						if _, err := e.file.WriteString(line); err != nil {
							e.logger.Error("failed to write metric to file", "error", err)
							return fmt.Errorf("failed to write metric to file: %w", err)
						}
						e.logger.Info("metric written to file", "metric", metricData.Name, "value", dp.Value, "file", e.filePath)
						return nil
					}
				}
			}
		}
	}
	return nil
}

func (e *FileExporter) Shutdown(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.file != nil {
		if err := e.file.Close(); err != nil {
			e.logger.Error("failed to close metrics file", "error", err)
			return fmt.Errorf("failed to close metrics file: %w", err)
		}
	}
	return nil
}
