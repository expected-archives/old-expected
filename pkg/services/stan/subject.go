package stan

import "fmt"

const (
	SubjectImageDelete      = "image:delete"
	SubjectImageDeleteLayer = "image:delete:layer"
)

func SubjectMetricNamespaceID(namespaceId string) string {
	return fmt.Sprintf("metric:%s", namespaceId)
}
