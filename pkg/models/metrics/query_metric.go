package metrics

import (
	"context"
	"github.com/expectedsh/expected/pkg/apps/agent/metrics"
	"github.com/expectedsh/expected/pkg/services"
)

func CreateMetric(ctx context.Context, metric metrics.Metric) error {
	_, err := services.Postgres().Client().ExecContext(ctx, `
		INSERT INTO metrics VALUES ($8, $1, $2, $3, $4, $5, $6, $7)
	`, metric.ID, metric.Memory, metric.NetInput, metric.BlockOutput, metric.BlockInput, metric.BlockOutput,
		metric.Cpu, metric.Time)
	return err
}
