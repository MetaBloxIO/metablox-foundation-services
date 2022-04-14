package main

import (
	"github.com/metabloxDID/contract"
	"github.com/metabloxDID/controllers"
	"github.com/metabloxDID/daily"
	"github.com/metabloxDID/dao"
	"github.com/metabloxDID/routers"
	"github.com/metabloxDID/settings"
	logger "github.com/sirupsen/logrus"
)

func main() {
	err := settings.Init()
	if err != nil {
		logger.Error(err)
		return
	}

	err = dao.InitSql()
	if err != nil {
		logger.Error(err)
		return
	}

	err = contract.Init()
	if err != nil {
		logger.Error(err)
		return
	}

	controllers.Init()
	err = daily.StartDailyTimer()
	if err != nil {
		logger.Error(err)
		return
	}

	routers.Setup()
}
