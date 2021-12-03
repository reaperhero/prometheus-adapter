package repository

import (
	"fmt"
	"github.com/reaperhero/prometheus-adapter/model"
	"github.com/robfig/cron"
	"github.com/toolkits/pkg/logger"
	"testing"
	"time"
	xlog "xorm.io/xorm/log"
)

func init() {
	InitDB()
	Xrepo.SyncTable()
}

func TestAddCapmetrics(t *testing.T) {
	c := model.CaptrueMetric{
		CapName:  "kafka_partitions_count",
		CapSql:   "kafka_partitions_count{}",
		Status:   true,
		Instance: "http://10.10.4.36:9090",
	}
	if _, err := Xrepo.M.InsertOne(&c); err != nil {
		logger.Info(err)
	}
}

func TestUpdateCapmetrics(t *testing.T) {
	c := model.CaptrueMetric{}
	_, err := Xrepo.M.Where("id=?", 0).Get(&c)
	if err != nil {
		logger.Error(err)
	}
	c.Instance = "http://10.10.4.36:9090"
	code, err := Xrepo.M.Update(c)
	if err != nil {
		return
	}
	logger.Info(code)
}

func TestGetCapmetrics(t *testing.T) {
	c := model.CaptrueMetric{}
	has, err := Xrepo.M.Where("id=?", 0).Get(&c)
	if err != nil {
		logger.Error(err)
	}
	logger.Info(c, has)
}

func TestDeleteCapmetrics(t *testing.T) {
	ok, err := Xrepo.M.Where("id=?", 0).Delete(&model.CaptrueMetric{})
	if err != nil {
		logger.Error(err)
	}
	logger.Info(ok)
}

func TestPrintSecond(t *testing.T) {
	fmt.Println(60 - time.Now().Second())
}

func TestCronSecond(t *testing.T) {
	cron := cron.New()
	cron.AddFunc("*/5 * * * * ?", func() {
		fmt.Println(1)
	})
	cron.Start()
	time.Sleep(time.Minute)
}

func print(i int) {
	fmt.Println(i)
}

func TestPrint(t *testing.T) {
	for i := 0; i < 100000; i++ {
		go print(i)
	}
	time.Sleep(time.Second * 3)
}

func TestPrintLables(t *testing.T) {
	go Xrepo.syncFilterLable()
	time.Sleep(time.Second * 1)
	fmt.Println(Xrepo.GetFilterLable())
}

func TestSlivePrint(t *testing.T) {
	a := []int{}
	for i := 0; i < 20000; i++ {
		a = append(a, i)
	}
	for _, i2 := range a {
		go func() {
			print(i2)
		}()
	}
	time.Sleep(time.Second * 3)
}

func TestIn(t *testing.T) {
	a := []model.CaptrueMetric{}
	Xrepo.M.Logger().SetLevel(xlog.LOG_DEBUG)
	Xrepo.M.ShowSQL(true)
	if err := Xrepo.M.Where("instance=?", "http://10.10.4.36:9090").In("cap_sql", "net_bytes_recv{}").Find(&a); err != nil {
		fmt.Println(err)
	}

}
