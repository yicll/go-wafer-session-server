package models

import (
	"time"

	"wafer-session-server/common"
	"wafer-session-server/common/utils"
	"wafer-session-server/models/mysql_models"
	"wafer-session-server/models/redis_models"

	"github.com/astaxie/beego"
)

// 微信鉴权接口返回数据格式
type WxResponseData struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	Errcode    int    `json:"errcode"`
	ErrMsg     string `json:"errMsg"`
}

// appinfo基本信息
type Appinfo struct {
	Appid           string
	Secret          string
	LoginDuration   int
	SessionDuration int
	QcloudAppid     string
	Ip              string
}

// session info
type SessionInfo struct {
	Appid         string
	Uuid          string
	Skey          string
	CreateTime    time.Time
	LastVisitTime time.Time
	OpenId        string
	SessionKey    string
	UserInfo      string
}

// 通过appid获取appinfo相关信息
func GetAppinfoById(appid string) (r Appinfo, err error) {
	if beego.AppConfig.String("dbdriver") == "mysql" {
		appinfo, err := mysql_models.GetAppinfoById(appid)
		if err != nil {
			return r, err
		}

		var result Appinfo
		utils.Assign(&appinfo, &result)

		return result, err
	} else {
		appinfo, err := redis_models.GetAppinfById(appid)
		if err != nil {
			return r, err
		}

		var result Appinfo
		utils.Assign(&appinfo, &result)

		return result, err
	}

	return r, err
}

// 更新登录态信息
func SaveSessionInfo(appinfo Appinfo, info *SessionInfo) error {

	if beego.AppConfig.String("dbdriver") == "mysql" {

		var sInfo mysql_models.CSessionInfo
		utils.Assign(info, &sInfo)

		return mysql_models.SaveSessionInfo(&sInfo)
	} else {

		var sess redis_models.RSessionInfo
		utils.Assign(info, &sess)
		sess.CreateTime = info.CreateTime.Unix()
		sess.LastVisitTime = info.LastVisitTime.Unix()
		var expire int
		if appinfo.LoginDuration*86400 < appinfo.SessionDuration {
			expire = appinfo.LoginDuration * 86400
		} else {
			expire = appinfo.SessionDuration
		}

		return redis_models.SaveSessionInfo(&sess, expire)
	}

	return nil
}

func CheckSession(uuid string, skey string, appinfo Appinfo) (r common.ResponseData, err error) {
	if beego.AppConfig.String("dbdriver") == "mysql" {
		return mysql_models.CheckSession(appinfo.Appid,
			uuid, skey, appinfo.LoginDuration, appinfo.SessionDuration)
	} else {
		return redis_models.CheckSession(appinfo.Appid, uuid,
			skey, appinfo.LoginDuration, appinfo.SessionDuration)
	}

	return
}

func InitApp(appid string, secret string) (r common.ResponseData, err error) {
	if beego.AppConfig.String("dbdriver") == "mysql" {
		return mysql_models.InitApp(appid, secret)
	} else {
		return redis_models.InitApp(appid, secret)
	}
}
