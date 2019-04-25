package metricsagent

import (
	"github.com/expectedsh/expected/pkg/apps/metricsagent/collector"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/sirupsen/logrus"
)

// GetMetrics send a stream to the caller.
func (a *App) GetMetrics(_ *protocol.MetricsRequest, res protocol.Metrics_GetMetricsServer) error {
	packets := collector.GetPackets()
	unsendedPackets := make([][]byte, 0)

	for _, packet := range packets {
		err := res.Send(&protocol.MetricsResponse{
			Metrics: packet,
		})

		if err != nil {
			logrus.WithError(err).Error("can't send stream")
			unsendedPackets = append(unsendedPackets, packet...)
		}
	}
	collector.AddPackets(unsendedPackets)
	return nil
}