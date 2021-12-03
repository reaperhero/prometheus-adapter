package http

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/reaperhero/prometheus-adapter/model"
	"github.com/reaperhero/prometheus-adapter/model/repository"
	"strconv"
)

type httphandler struct {
}

func (h *httphandler) getRules(c echo.Context) error {
	instance := c.QueryParam("instance")
	var (
		count int64
		list  []model.CaptrueMetric
	)
	if instance != "" {
		count, list = repository.Xrepo.GetCaptrueRuleWithInstance(instance)
	} else {
		count, list = repository.Xrepo.GetCaptrueRule()
	}

	type rule struct {
		Id     int64  `json:"id"`
		Status bool   `json:"status"`
		Sql    string `json:"sql"`
	}
	type respParam struct {
		Count int64             `json:"count"`
		Rules map[string][]rule `json:"rules"`
	}
	response := respParam{
		Count: count,
		Rules: make(map[string][]rule),
	}
	for _, metric := range list {
		r := response.Rules[metric.Instance]
		r = append(r, rule{
			Id:     metric.Id,
			Status: metric.Status,
			Sql:    metric.CapSql,
		})
		response.Rules[metric.Instance] = r
	}
	return c.JSON(200, response)
}

func (h *httphandler) createRules(c echo.Context) error {
	param := model.CaptrueMetric{}
	if err := c.Bind(param); err != nil {
		return c.JSON(200, fmt.Errorf("param is err,%s", err))
	}
	response := struct {
		Msg  string      `json:"msg"`
		Data interface{} `json:"data,omitempty"`
	}{}
	list := repository.Xrepo.GetRulesWithInstance(param.Instance, param.CapSql)
	if len(list) != 0 {
		response.Msg = "rule is exist"
		response.Data = list
		return c.JSON(200, response)
	}
	repository.Xrepo.CreateRuleWithInstance(param.Instance, param.CapName, param.CapSql)
	response.Msg = "ok"
	return c.JSON(200, response)
}

func (h *httphandler) deleteRules(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id == 0 {
		return c.JSON(200, "param is err")
	}

	err := repository.Xrepo.DeleteRulesWithIds(id)

	return c.JSON(200, model.GetErrorMap(err))
}

func (h *httphandler) updateRules(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	request := model.CaptrueMetric{}
	if err := c.Bind(request); err != nil {
		return c.JSON(200, model.GetErrorMap(model.ErrInvalidParam))
	}
	if id == 0 {
		return c.JSON(200, "param is err")
	}
	request.Id = id
	err := repository.Xrepo.UpdateRulesWithIds(request)

	return c.JSON(200, model.GetErrorMap(err))
}

func (h *httphandler) updateRuleStatus(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	status, err := strconv.ParseBool(c.QueryParam("status"))
	if id < 0 || err != nil {
		return c.JSON(200, model.GetErrorMap(model.ErrInvalidParam))
	}
	err = repository.Xrepo.UpdateRulesStatusWithId(id, status)

	return c.JSON(200, model.GetErrorMap(err))
}
