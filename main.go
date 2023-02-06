package main

import (
	"github.com/MetaBloxIO/metablox-foundation-services/comm/log"
	"github.com/MetaBloxIO/metablox-foundation-services/contract"
	"github.com/MetaBloxIO/metablox-foundation-services/controllers"
	"github.com/MetaBloxIO/metablox-foundation-services/dao"
	"github.com/MetaBloxIO/metablox-foundation-services/did"
	"github.com/MetaBloxIO/metablox-foundation-services/routers"
	"github.com/MetaBloxIO/metablox-foundation-services/service"
	"github.com/MetaBloxIO/metablox-foundation-services/settings"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	pflag.Parse()

	err := settings.Init()
	if err != nil {
		logger.Error(err)
		return
	}

	loggerConf := &log.Config{}
	viper.UnmarshalKey("logger", loggerConf)
	err = log.Init(loggerConf)
	if err != nil {
		logger.Error(err)
		return
	}

	err = dao.InitSql()
	if err != nil {
		logger.Error(err)
		return
	}

	sqlConf := &dao.Config{}
	viper.UnmarshalKey("wifiDB", sqlConf)
	err = dao.InitWifiDB(sqlConf)
	if err != nil {
		logger.Error(err)
		return
	}

	err = contract.Init()
	if err != nil {
		logger.Error(err)
		return
	}

	controllers.InitializeValues()

	err = did.Init(&did.Config{
		Passphrase: viper.GetString("metablox.wallet.passphrase"),
		Keystore:   viper.GetString("metablox.wallet.keystore"),
	})
	if err != nil {
		logger.Error(err)
		return
	}

	if err = service.InitEvent(); err != nil {
		logger.Error(err)
		return
	}

	routers.Setup()
}
