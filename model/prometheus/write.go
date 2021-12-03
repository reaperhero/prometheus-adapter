package prometheus

import (
	"errors"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/prompb"
	"github.com/sirupsen/logrus"
	"github.com/toolkits/pkg/logger"
	"regexp"
	"time"
)

type MetricPoint struct {
	Metric  string            `json:"metric"` // 监控指标名称
	TagsMap map[string]string `json:"tags"`   // 监控数据标签
	Time    int64             `json:"time"`   // 时间戳，单位是秒
	Value   float64           `json:"-"`      // 内部字段，最终转换之后的float64数值
}

type sample struct {
	labels labels.Labels
	t      int64
	v      float64
}

var MetricNameRE = regexp.MustCompile(`^[a-zA-Z_:][a-zA-Z0-9_:]*$`)

func (p *PromeDataSource) convertOne(item *MetricPoint) (prompb.TimeSeries, error) {
	pt := prompb.TimeSeries{}
	pt.Samples = []prompb.Sample{{}}
	s := sample{}
	s.t = item.Time
	s.v = item.Value
	// name
	if !MetricNameRE.MatchString(item.Metric) {
		return pt, errors.New("invalid metrics name")
	}

	for k, v := range item.TagsMap {
		if model.LabelNameRE.MatchString(k) {
			ls := labels.Label{
				Name:  k,
				Value: v,
			}
			s.labels = append(s.labels, ls)
		}
	}

	pt.Labels = labelsToLabelsProto(s.labels, pt.Labels)
	pt.Samples[0].Timestamp = s.t
	pt.Samples[0].Value = s.v
	return pt, nil
}

func labelsToLabelsProto(labels labels.Labels, buf []prompb.Label) []prompb.Label {
	result := buf[:0]
	if cap(buf) < len(labels) {
		result = make([]prompb.Label, 0, len(labels))
	}
	for _, l := range labels {
		result = append(result, prompb.Label{
			Name:  l.Name,
			Value: l.Value,
		})
	}
	return result
}

func (p *PromeDataSource) Push2Queue(points []*MetricPoint) {
	for _, point := range points {
		pt, err := p.convertOne(point)
		if err != nil {
			logger.Errorf("[prome_convertOne_error][point: %+v][err:%s]", point, err)
			continue
		}
		ok := p.PushQueue.PushFront(pt)
		if !ok {
			logger.Errorf("[prome_push_queue_error][point: %+v] ", point)
		}
	}
}

func (p *PromeDataSource) remoteWrite() {
	for {
		items := p.PushQueue.PopBackBy(1000)
		count := len(items)
		if count == 0 {
			time.Sleep(time.Second * 5)
			continue
		}
		logrus.Debug(items)
		pbItems := make([]prompb.TimeSeries, count)
		for i := 0; i < count; i++ {
			pbItems[i] = items[i].(prompb.TimeSeries)
		}
		payload, err := p.buildWriteRequest(pbItems)
		if err != nil {
			logger.Errorf("[prome_remote_write_error][pb_marshal_error][items: %+v][pb.err: %v]: ", items, err)
			continue
		}
		p.processWrite(payload)
	}
}
