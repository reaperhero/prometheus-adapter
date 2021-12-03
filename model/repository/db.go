package repository

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/reaperhero/prometheus-adapter/config"
	"github.com/reaperhero/prometheus-adapter/model"
	"github.com/toolkits/pkg/logger"
	"log"
	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"
)

type XormRepo struct {
	M *xorm.Engine
}

var (
	Xrepo XormRepo
)

func InitDB() {
	// Get/Exist/Find/Iterate/Count/Rows/Sum  Cols
	engine, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8",
		config.GetEnvWithDeafult("MYSQL_USER", "root"),
		config.GetEnvWithDeafult("MYSQL_PASS", "Lzslov123!"),
		config.GetEnvWithDeafult("MYSQL_ADDR", "127.0.0.1"),
		config.GetEnvWithDeafult("MYSQL_DB", "mudu_export"),
	))
	engine.Logger().SetLevel(xlog.LOG_ERR)
	engine.ShowSQL(false)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	Xrepo = XormRepo{M: engine}
	Xrepo.SyncTable()
	go Xrepo.syncFilterLable()
}

func (x *XormRepo) SyncTable() {
	err := x.M.Sync(
		new(model.CaptrueMetric),
		new(model.FilterLabel),
	)
	if err != nil {
		logger.Errorf("[sync table err %s]", err)
	}
}
