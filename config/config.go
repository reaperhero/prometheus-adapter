package config

import (
	"go.uber.org/atomic"
	"net/http"
	"net/url"
	"time"
)

var (
	READ_API = GetArrayEnvWithDeafult("READ_API", "")

	Config = PromeSection{
		Batch:                        100,
		MaxRetry:                     5,
		LookbackDeltaMinute:          2,
		MaxConcurrentQuery:           10,
		MaxSamples:                   50000,
		MaxFetchAllSeriesLimitMinute: 5,
		SlowLogRecordSecond:          5,
		RemoteRead:                   getremoteReadsConfig(),
		RemoteWrite: RemoteConfig{
			Url:                 GetEnvWithDeafult("WRITE_API", "http://10.20.23.6:9090") + "/api/v1/write",
			RemoteTimeoutSecond: 40,
		},
	}
)

func getremoteReadsConfig() (reads []RemoteConfig) {
	for _, s := range READ_API {
		reads = append(reads, RemoteConfig{
			Url:                 s + "/api/v1/read",
			RemoteTimeoutSecond: 40,
		})
	}
	return
}

type RemoteConfig struct {
	Url                 string `yaml:"url"`
	RemoteTimeoutSecond int    `yaml:"remoteTimeoutSecond"`
}

type HttpClient struct {
	remoteName string // Used to differentiate clients in metrics.
	Url        *url.URL
	Client     *http.Client
	Timeout    time.Duration
}

type PromeSection struct {
	Batch                        int            `yaml:"batch"`
	MaxRetry                     int            `yaml:"maxRetry"`
	LookbackDeltaMinute          int            `yaml:"lookbackDeltaMinute"`
	MaxConcurrentQuery           int            `yaml:"maxConcurrentQuery"`
	MaxSamples                   int            `yaml:"maxSamples"`
	MaxFetchAllSeriesLimitMinute int64          `yaml:"maxFetchAllSeriesLimitMinute"`
	SlowLogRecordSecond          float64        `yaml:"slowLogRecordSecond"`
	RemoteRead                   []RemoteConfig `yaml:"remoteRead"`
	RemoteWrite                  RemoteConfig   `yaml:"remoteWrite"`
}

type SafePromQLNoStepSubqueryInterval struct {
	value atomic.Int64
}

func (i *SafePromQLNoStepSubqueryInterval) Get(int64) int64 {
	return i.value.Load()
}
