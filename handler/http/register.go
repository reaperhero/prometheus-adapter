package http

import (
	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"os"
)

func HttpRun() {
	handler := httphandler{}
	e := echo.New()
	e.HideBanner = true
	e.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Skipper: echoMiddleware.DefaultSkipper,
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Output:           os.Stdout,
	}))
	e.Use(echoMiddleware.Recover())

	handler.setRouter(e)

	if err := e.Start(":80"); err != nil {
		logrus.WithError(err).Fatalln("[echo.Start]")
	}
}

func (h *httphandler) setRouter(e *echo.Echo) {
	home := e.Group("/adapter")
	home.GET("/log/:level", func(c echo.Context) error { //设置本服务的日志等级
		logLevel := c.Param("level")

		switch logLevel {
		case "1":
			logrus.SetLevel(logrus.DebugLevel)
		case "2":
			logrus.SetLevel(logrus.InfoLevel)
		case "3":
			logrus.SetLevel(logrus.WarnLevel)
		case "4":
			logrus.SetLevel(logrus.ErrorLevel)
		}

		log.Debugln("this is a debug log")
		log.Infoln("this is a info log")
		log.Warnln("this is a warn log")
		log.Errorln("this is a error log")

		return c.String(200, "set log level success")
	})
	home.GET("/rules", h.getRules)
	home.POST("/rule", h.createRules)
	home.DELETE("/rule/:id", h.deleteRules)
	home.PUT("/rule/:id", h.updateRules)
	home.PUT("/rulestatus/:id", h.updateRuleStatus)
}
