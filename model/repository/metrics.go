package repository

import (
	"github.com/reaperhero/prometheus-adapter/model"
	"github.com/sirupsen/logrus"
	"github.com/toolkits/pkg/logger"
)

type Instance struct {
	Ip       string `json:"ip"`
	Interval int64  `json:"interval"`
}

func (x *XormRepo) GetCaptrueRuleWithTrue(instance string) (cs []model.CaptrueMetric) {
	if err := x.M.Where("status=? and instance=?", true, instance).Find(&cs); err != nil {
		logger.Errorf("[get metrics err %s]", err)
		return nil
	}
	return
}

func (x *XormRepo) GetCaptrueRule() (count int64, cs []model.CaptrueMetric) {
	count, err := x.M.FindAndCount(&cs)
	if err != nil {
		logrus.Errorf("[XormRepo.GetCaptrueRule] %s", err)
	}
	return
}

func (x *XormRepo) GetCaptrueRuleWithInstance(instance string) (count int64, cs []model.CaptrueMetric) {
	count, err := x.M.Where("instance=?", instance).FindAndCount(&cs)
	if err != nil {
		logrus.Errorf("[XormRepo.GetCaptrueRule] %s", err)
	}
	return
}
