package ingester

import (
	"github.com/expectedsh/expected/pkg/apps/metricsagent/metrics"
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/stan"
	"github.com/sirupsen/logrus"
	"sync"
)

func Ingest(metricList []metrics.Metric) {
	group := sync.WaitGroup{}
	group.Add(len(metricList))

	metricsSended := 0

	for _, packet := range metricList {
		go func(metric metrics.Metric) {
			data, err := metric.MarshalBinary()
			if err != nil {
				logrus.WithError(err).Error("can't marshal into binary metric")
				group.Done()
				return
			}

			// todo maybe change the ID with the namespace ID ?
			err = services.Stan().Client().Publish(stan.SubjectMetricNamespaceID(packet.ID.String()), data)
			if err != nil {
				logrus.WithError(err).Error("can't ingest metric")
			} else {
				metricsSended++
			}
			group.Done()
		}(packet)
	}
	group.Wait()
	logrus.Infof("%d/%d metrics ingested.", metricsSended, len(metricList))
}
