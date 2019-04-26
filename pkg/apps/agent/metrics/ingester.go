package metrics

import (
	"github.com/expectedsh/expected/pkg/services"
	"github.com/expectedsh/expected/pkg/services/stan"
	"github.com/sirupsen/logrus"
	"sync"
)

func ingest(metricList []Metric) {
	group := sync.WaitGroup{}
	group.Add(len(metricList))

	metricsSended := 0

	for _, packet := range metricList {
		go func(metric Metric) {
			data, err := metric.MarshalBinary()
			if err != nil {
				logrus.WithError(err).Error("can't marshal into binary metric")
				group.Done()
				return
			}

			err = services.Stan().Client().Publish(stan.SubjectMetric, data)
			if err != nil {
				logrus.WithError(err).Error("can't ingest metric")
			} else {
				metricsSended++
			}
			group.Done()
		}(packet)
	}
	group.Wait()
	logrus.Infof("%d/%d metrics ingested", metricsSended, len(metricList))
}
