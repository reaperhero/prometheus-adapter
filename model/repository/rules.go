package repository

import (
	"github.com/reaperhero/prometheus-adapter/model"
	"github.com/sirupsen/logrus"
	"strings"
)

func (x *XormRepo) GetRulesWithInstance(host string, rules string) (caps []model.CaptrueMetric) {
	if err := x.M.Where("instance=? and cap_sql=?", host, rules).Find(&caps); err != nil {
		logrus.Errorf("[XormRepo.GetRulesWithInstance] %s", err)
	}
	return
}

func (x *XormRepo) CreateRuleWithInstance(host string, capname string, rule string) {
	if capname == "" {
		caps := strings.Split(rule, "{")
		capname = caps[0]
	}
	cap := model.CaptrueMetric{
		CapName:  capname,
		CapSql:   rule,
		Status:   true,
		Instance: host,
	}
	x.M.Insert(&cap)

	return
}

func (x *XormRepo) DeleteRulesWithIds(id int64) (err error) {
	if _, err := x.M.Where("id=?", id).Delete(new(model.CaptrueMetric)); err != nil {
		logrus.Errorf("[XormRepo.DeleteRulesWithIds] %s", err)
		return model.ErrDbOperation
	}
	return nil
}

func (x *XormRepo) UpdateRulesWithIds(rule model.CaptrueMetric) (err error) {
	idData := model.CaptrueMetric{}
	if err := x.M.Where("id=?", rule.Id).Find(&idData); err != nil {
		logrus.Infof("[XormRepo.UpdateRulesWithIds] %s", err)
		return model.ErrDbOperation
	}
	if rule.CapName != "" {
		idData.CapName = rule.CapName
	}
	if rule.CapSql != "" {
		idData.CapSql = rule.CapSql
	}
	if rule.Instance != "" {
		idData.Instance = rule.Instance
	}
	if err := x.M.Sync2(idData); err != nil {
		return model.ErrDbOperation
	}
	return nil
}

func (x *XormRepo) UpdateRulesStatusWithId(id int64, status bool) (err error) {
	idData := model.CaptrueMetric{}
	if err := x.M.Where("id=?", id).Find(&idData); err != nil {
		logrus.Infof("[XormRepo.UpdateRulesWithIds] %s", err)
		return model.ErrDbOperation
	}
	idData.Status = status
	if err := x.M.Sync2(idData); err != nil {
		return model.ErrDbOperation
	}
	return nil
}
