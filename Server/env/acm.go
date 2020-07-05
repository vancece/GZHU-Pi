package env

import (
	"GZHU-Pi/services/acm"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/viper"
)

var Acm *acm.Acm

func InitAcm() (a *acm.Acm, err error) {
	enable := viper.GetBool("acm.enabled")
	if !enable {
		logs.Info("disable acm")
		Acm = acm.DefaultDisableAcm()
		a = Acm
		return
	}

	accessKey := viper.GetString("acm.access_key")
	secretKey := viper.GetString("acm.secret_key")
	namespaceId := viper.GetString("acm.namespace_id")
	group := viper.GetString("acm.group")

	Acm, err = acm.DefaultAcm(accessKey, secretKey, namespaceId, group)
	if err != nil {
		return
	}
	a = Acm
	logs.Info("acm init group:%s", a.DefaultGroup)
	return
}
