package mysql_models

import (
	"github.com/astaxie/beego/orm"
)

type CAppinfo struct {
	Appid           string `orm:"column(appid);pk"`
	Secret          string `orm:"column(secret);size(300)"`
	LoginDuration   int    `orm:"column(login_duration);null"`
	SessionDuration int    `orm:"column(session_duration);null"`
	QcloudAppid     string `orm:"column(qcloud_appid);size(300);null"`
	Ip              string `orm:"column(ip);size(50);null"`
}

func (t *CAppinfo) TableName() string {
	return "cAppinfo"
}

func init() {
	orm.RegisterModel(new(CAppinfo))
}
