package prometheus

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/common/promlog"
	pc "github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/storage/remote"
	"github.com/reaperhero/prometheus-adapter/config"
	"github.com/reaperhero/prometheus-adapter/model/repository"
	"github.com/sirupsen/logrus"
	"github.com/toolkits/pkg/container/list"
	"github.com/toolkits/pkg/logger"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type PromeDataSource struct {
	LocalTmpDir  string
	QueryableMap map[string]storage.SampleAndChunkQueryable
	EngineMap    map[string]*promql.Engine // read
	PushQueue    *list.SafeListLimited
	WriteTargets []*config.HttpClient // write
}

func NewPromeDataSource() DataSource {
	p := PromeDataSource{
		PushQueue:    list.NewSafeListLimited(100000),
		QueryableMap: make(map[string]storage.SampleAndChunkQueryable),
		EngineMap:    make(map[string]*promql.Engine),
	}
	for k, remoteConfig := range config.Config.RemoteRead {
		// 模拟创建本地存储目录
		dbDir, err := ioutil.TempDir("", fmt.Sprintf("tsdb-api-ready-%d", k))
		if err != nil {
			logger.Errorf("[error_create_local_tsdb_dir][err: %v]", err)
			return nil
		}
		ur, err := url.Parse(remoteConfig.Url)
		if err != nil {
			logger.Errorf("[prome_ds_init_error][parse_url_error][url:%+v][err:%+v]", remoteConfig.Url, err)
			return nil
		}
		remoteS := remote.NewStorage(promlog.New(&promlog.Config{}), nil, func() (int64, error) {
			return 0, nil
		}, dbDir, 1*time.Minute, nil)
		remoteReadC := &pc.RemoteReadConfig{
			URL:           &config_util.URL{URL: ur},
			RemoteTimeout: model.Duration(time.Duration(remoteConfig.RemoteTimeoutSecond) * time.Second),
			ReadRecent:    true,
		}
		if err = remoteS.ApplyConfig(&pc.Config{RemoteReadConfigs: []*pc.RemoteReadConfig{remoteReadC}}); err != nil {
			logger.Errorf("%s", err)
		}
		pLogger := log.NewNopLogger()

		noStepSubqueryInterval := &config.SafePromQLNoStepSubqueryInterval{}

		queryQueueDir, err := ioutil.TempDir(dbDir, fmt.Sprintf("prom_query_concurrency_%d", k))
		if err != nil {
			logger.Errorf("[error info %s]", err)
		}

		opts := promql.EngineOpts{
			Logger:                   log.With(pLogger, "component", "query engine"),
			Reg:                      nil,
			MaxSamples:               config.Config.MaxSamples,
			Timeout:                  30 * time.Second,
			ActiveQueryTracker:       promql.NewActiveQueryTracker(queryQueueDir, config.Config.MaxConcurrentQuery, log.With(pLogger, "component", "activeQueryTracker")),
			LookbackDelta:            time.Duration(config.Config.LookbackDeltaMinute) * time.Minute,
			NoStepSubqueryIntervalFn: noStepSubqueryInterval.Get,
			EnableAtModifier:         true,
		}

		p.EngineMap[ur.Host] = promql.NewEngine(opts)
		p.QueryableMap[ur.Host] = remoteS
	}

	// write
	wur, _ := url.Parse(config.Config.RemoteWrite.Url)
	p.WriteTargets = []*config.HttpClient{{
		Url:     wur,
		Client:  &http.Client{},
		Timeout: time.Duration(config.Config.RemoteWrite.RemoteTimeoutSecond) * time.Second,
	}}
	//queus consume
	go p.remoteWrite()
	return &p
}

func (p *PromeDataSource) QueryVector(instance string, ql string) promql.Vector {
	t := time.Now()
	q, err := p.EngineMap[instance].NewInstantQuery(p.QueryableMap[instance], ql, t)
	if err != nil {
		logger.Errorf("[prome_query_error][new_insQuery_error][err:%+v][ql:%+v]", err, ql)
		return nil
	}
	res := q.Exec(context.Background())

	if res.Err != nil {
		logger.Errorf("[prome_query_error][insQuery_exec_error][err:%+v][ql:%+v]", err, ql)
		return nil
	}
	defer q.Close()
	switch v := res.Value.(type) {
	case promql.Vector:
		return v
	case promql.Scalar:
		return promql.Vector{promql.Sample{Point: promql.Point(v), Metric: labels.Labels{}}}
	default:
		logger.Errorf("[prome_query_error][insQuery_res_error rule result is not a vector or scalar][err:%+v][ql:%+v]", err, ql)
		return nil
	}
}

func (p *PromeDataSource) buildWriteRequest(samples []prompb.TimeSeries) ([]byte, error) {

	req := &prompb.WriteRequest{
		Timeseries: samples,
		Metadata:   nil,
	}

	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	compressed := snappy.Encode(nil, data)
	return compressed, nil
}

func (p *PromeDataSource) processWrite(payload []byte) {
	for _, c := range p.WriteTargets {
		newC := c
		go func(cc *config.HttpClient, payload []byte) {
			defer func() {
				cc.Client.CloseIdleConnections()
			}()
			for i := 0; i < 5; i++ {
				err := remoteWritePost(cc, payload)
				if err == nil {
					break
				}
				time.Sleep(time.Millisecond * 1000)
			}
		}(newC, payload)
	}
}

func remoteWritePost(c *config.HttpClient, req []byte) error {
	httpReq, err := http.NewRequest("POST", c.Url.String(), bytes.NewReader(req))
	if err != nil {
		return err
	}

	httpReq.Header.Add("Content-Encoding", "snappy")
	httpReq.Header.Set("Content-Type", "application/x-protobuf")
	//httpReq.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	httpReq = httpReq.WithContext(ctx)

	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		var ht *nethttp.Tracer
		httpReq, ht = nethttp.TraceRequest(
			parentSpan.Tracer(),
			httpReq,
			nethttp.OperationName("Remote Store"),
			nethttp.ClientTrace(false),
		)
		defer ht.Finish()
	}

	httpResp, err := c.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(ioutil.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	if httpResp.StatusCode/100 != 2 {
		scanner := bufio.NewScanner(io.LimitReader(httpResp.Body, 512))
		line := ""
		if scanner.Scan() {
			line = scanner.Text()
		}

		if httpResp.StatusCode == 400 {
			//400的错误是客户端的问题，不返回给上层，输出到debug日志中
			logrus.Errorf("server returned HTTP status %s: %s req:%v", httpResp.Status, line, getSamples(req))
		} else {
			err = errors.Errorf("server returned HTTP status %s: %s", httpResp.Status, line)
		}
	}

	return err
}

func getSamples(compressed []byte) []prompb.TimeSeries {
	var samples []prompb.TimeSeries
	req := &prompb.WriteRequest{
		Timeseries: samples,
		Metadata:   nil,
	}

	d, _ := snappy.Decode(nil, compressed)
	proto.Unmarshal(d, req)
	if len(req.Timeseries) == 0 {
		return nil
	}
	return req.Timeseries[:1]
}

func TransformLabel(ls labels.Labels, name string) map[string]string {
	var (
		found bool
	)
	newSlice := labels.Labels{}
	for _, l := range ls {
		if l.Name == "__name__" {
			found = true
			l.Value = name
		}
		if repository.FilterLabels.MatchValueWithValue(l.Name) {
			continue
		}
		newSlice = append(newSlice, l)
	}
	if !found {
		newSlice = append(newSlice, labels.Label{
			Name:  "__name__",
			Value: name,
		})
	}
	return newSlice.Map()
}
