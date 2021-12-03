package prometheus

import (
	"github.com/prometheus/prometheus/promql"
)

type DataSource interface {
	QueryVector(instance string, ql string) promql.Vector
	Push2Queue(points []*MetricPoint)
}
