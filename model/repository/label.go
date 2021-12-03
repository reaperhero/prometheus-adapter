package repository

import (
	"github.com/reaperhero/prometheus-adapter/model"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	filterLabelsWg = sync.Mutex{}
	FilterLabels   = filterLabels{}
)

type filterLabels []model.FilterLabel

func (f filterLabels) MatchValueWithValue(name string) bool {
	for _, label := range f {
		if label.LableName == name {
			return true
		}
	}
	return false
}

func (x *XormRepo) GetFilterLable() (int, []model.FilterLabel) {
	filterLabelsWg.Lock()
	defer filterLabelsWg.Unlock()
	return len(FilterLabels), FilterLabels
}

func (x *XormRepo) syncFilterLable() {
	for {
		var list []model.FilterLabel
		filterLabelsWg.Lock()
		err := x.M.Find(&list)
		FilterLabels = list
		filterLabelsWg.Unlock()
		if err != nil {
			logrus.Errorf("[XormRepo.GetFilterLable] %s", err)
		}
		time.Sleep(time.Minute)
	}
}
