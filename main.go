package main

import (
	"github.com/reaperhero/prometheus-adapter/config"
	"github.com/reaperhero/prometheus-adapter/handler/http"
	"github.com/reaperhero/prometheus-adapter/model"
	"github.com/reaperhero/prometheus-adapter/model/prometheus"
	"github.com/reaperhero/prometheus-adapter/model/repository"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strings"
	"time"
)

func initLogLevel(level string) {
	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}
func init() {
	initLogLevel(os.Getenv("LOG_LEVEL"))
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05",
	})
	rand.Seed(time.Now().Unix())
}

func main() {
	repository.InitDB()
	usecase := prometheus.NewPromeDataSource()
	Run(usecase)
	logrus.Println("prometheus adapter is running....")
	http.HttpRun()
}

func Run(source prometheus.DataSource) {
	cron := cron.New()
	cron.AddFunc("20 * * * * ?", func() {
		for _, instance := range config.READ_API {
			caps := repository.Xrepo.GetCaptrueRuleWithTrue(instance)
			instance = strings.TrimPrefix(instance, "http://")
			go runQueueMessage(caps, instance, source)
		}
	})
	cron.Start()
}

func runQueueMessage(caps []model.CaptrueMetric, instance string, source prometheus.DataSource) {
	rateLimit := make(chan struct{}, 20)
	for _, c := range caps {
		go func(c model.CaptrueMetric) {
			rateLimit <- struct{}{}
			if value := source.QueryVector(instance, c.CapSql); value != nil {
				metrics := make([]*prometheus.MetricPoint, len(value))
				for k, s := range value {
					ns := s
					metrics[k] = &prometheus.MetricPoint{
						Metric:  c.CapName,
						TagsMap: prometheus.TransformLabel(ns.Metric, c.CapName),
						Time:    s.T,
						Value:   s.V,
					}
				}
				source.Push2Queue(metrics)
			}
			<-rateLimit
		}(c)
	}
}
