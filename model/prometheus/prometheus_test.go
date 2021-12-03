package prometheus

import (
	"fmt"
	"testing"
	"time"
)

func TestGetPrometheus(t *testing.T) {
	p := NewPromeDataSource()
	values := p.QueryVector("", "kube_node_info{instance=\"172.16.0.19:8080\"}")
	for _, value := range values {
		fmt.Println(value)
	}
}

var (
	sql = "ceil((1-((sum(increase(node_cpu{mode=\"idle\"}[2m])) by (instance,tags)) / (sum(increase(node_cpu{}[2m])) by (instance,tags)))) * 100)"
)

func TestSendPrometheus(t *testing.T) {
	repo := NewPromeDataSource()
	value := repo.QueryVector("", sql)
	metrics := []*MetricPoint{}
	for _, s := range value {
		metrics = append(metrics, &MetricPoint{
			Metric:  "node_cpu_new",
			TagsMap: s.Metric.Map(),
			Time:    s.T,
			Value:   s.V,
		})
	}
	repo.Push2Queue(metrics)
}

func TestSendOneMetricstoPrometheus(t *testing.T) {
	repo := NewPromeDataSource()
	value := repo.QueryVector("", sql)
	metrics := []*MetricPoint{}
	for _, s := range value {
		metrics = append(metrics, &MetricPoint{
			Metric:  "node_cpu_new",
			TagsMap: s.Metric.Map(),
			Time:    s.T,
			Value:   s.V,
		})
	}
}

func TestConnery(t *testing.T) {
	rateLimit := make(chan struct{}, 10)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			rateLimit <- struct{}{}
			time.Sleep(time.Second)
			fmt.Println(i)
			<-rateLimit
		}(i)
	}
	select {}
}

func TestGeneralCloud(t *testing.T) {
	metrics := []string{
		"es_cpu_load_count",
	}
	instance := []string{
		"http://10.10.4.98:9090",
	}
	result := []string{}
	sql := "insert into captrue_metric(cap_name,cap_sql,status,instance) value(\"%s\",\"%s{}\",1,\"%s\");"
	for _, m := range metrics {
		for _, i := range instance {
			result = append(result, fmt.Sprintf(sql, m, m, i))
		}
	}
	for _, i := range result {
		fmt.Println(i)
	}
}

func TestGeneralPrint1(t *testing.T) {
	metrics := []string{
		"slb_inactive_conn_count",
		"slb_active_conn_count",
		"slb_inactive_conn_count",
		"slb_new_conn_count",
		"slb_down_data_pkg_count",
		"slb_up_data_pkg_count",
		"slb_down_bandwidth_kbps",
		"slb_up_bandwidth_kbps",
		"slb_l7_http_2xx_count",
		"slb_l7_http_3xx_count",
		"slb_l7_http_4xx_count",
		"slb_l7_http_5xx_count",
		"slb_l7_http_other_status_count",
		"slb_l7_http_404_count",
		"slb_l7_http_499_count",
		"slb_l7_http_502_count",
		"slb_l7_rt_ms",
		"slb_l7_qps_count",
	}
	instance := []string{
		"http://10.0.0.1:9090",
	}
	result := []string{}
	sql := "insert into captrue_metric(cap_name,cap_sql,status,instance) value(\"%s\",\"%s{}\",1,\"%s\");"
	for _, m := range metrics {
		for _, i := range instance {
			result = append(result, fmt.Sprintf(sql, m, m, i))
		}
	}
	for _, i := range result {
		fmt.Println(i)
	}
}

func TestGeneralPrint2(t *testing.T) {
	metrics := []string{
		"kube_node_status_condition",
	}
	instance := []string{
		"http://10.0.0.1:9090",
	}
	result := []string{}
	sql := "insert into captrue_metric(cap_name,cap_sql,status,instance) value(\"%s\",\"%s{}\",1,\"%s\");"
	for _, m := range metrics {
		for _, i := range instance {
			result = append(result, fmt.Sprintf(sql, m, m, i))
		}
	}
	for _, i := range result {
		fmt.Println(i)
	}
}
