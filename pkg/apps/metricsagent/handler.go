package metricsagent

import (
	"github.com/expectedsh/expected/pkg/apps/metricsagent/collector"
	"github.com/expectedsh/expected/pkg/protocol"
	"github.com/sirupsen/logrus"
)

func (a *App) GetMetrics(_ *protocol.MetricsRequest, res protocol.Metrics_GetMetricsServer) error {
	packets := collector.GetPackets()
	unsendedPackets := make([][]byte, 0)
	for _, packet := range packets {

		error := res.Send(&protocol.MetricsResponse{
			Metrics: packet,
		})

		if error != nil {
			logrus.WithError(error).Error("can't send stream")
			unsendedPackets = append(unsendedPackets, packet...)
		}
	}
	collector.AddPackets(unsendedPackets)
	return nil
}
