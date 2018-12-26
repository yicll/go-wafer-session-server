package mysql_models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	"wafer-session-server/common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 通过appid获取appinfo相关信息
func GetAppinfoById(appid string) (r CAppinfo, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CAppinfo))
	if appid == "" {
		qs = qs.Limit(1)
	} else {
		qs = qs.Filter("Appid", appid)
	}

	var result CAppinfo
	err = qs.One(&result)
	if err == orm.ErrNoRows {
		beego.Error(fmt.Sprintf("query appinfo from db error, appid: %s, message: %s", appid, err.Error()))
		err = common.ServerError{common.RETURN_CODE_MA_NO_APPID, err.Error()}
		return r, err
	}
	if err != nil {
		beego.Error(fmt.Sprintf("query appinfo from db error, appid: %s, message: %s", appid, err.Error()))
		err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, err.Error()}
		return r, err
	}

	return result, nil
}

// 更新db里的session信息
func SaveSessionInfo(info *CSessionInfo) (err error) {

	beego.Debug(fmt.Sprintf("session info: %+v", info))

	o := orm.NewOrm()
	qs := o.QueryTable(new(CSessionInfo))
	qs = qs.Filter("Appid", info.Appid).Filter("OpenId", info.OpenId)
	qs = qs.Limit(1).OrderBy("-Id")

	var oldInfo CSessionInfo
	err = qs.One(&oldInfo)

	if err != nil && err != orm.ErrNoRows {
		beego.Error("query session from db error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "QUERY_SESSION_DB_ERR"}
		return err
	}

	// 存在记录，则删除掉老记录再重新插入新记录
	if err != orm.ErrNoRows {
		if _, err := o.Delete(&oldInfo); err != nil {
			beego.Error("delete old session info error, message: ", err.Error())
			err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "DELETE_SESSION_DB_ERR"}
			return err
		}
	}
	beego.Informational("delete old session info success, id: ", oldInfo.Id)

	id, err := o.Insert(info)
	if err != nil {
		beego.Error("insert new session info to db error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "ADD_NEW_SESSION_DB_ERR"}
		return err
	}
	beego.Informational("add new session success, id: ", id)

	return nil
}

// 检查登录状态
func CheckSession(appid string, uuid string, skey string, lduration int, sduration int) (r common.ResponseData, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CSessionInfo))
	qs = qs.Filter("Appid", appid).Filter("Uuid", uuid).Filter("Skey", skey)
	qs = qs.Limit(1).OrderBy("-Id")

	var sInfo CSessionInfo
	err = qs.One(&sInfo)

	if err == orm.ErrNoRows {
		err = common.ServerError{common.RETURN_CODE_MA_AUTH_ERR, "AUTH_ERR_NO_EXIST"}
		return r, err
	}

	if err != nil {
		beego.Error("query session info from db error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_AUTH_ERR, "AUTH_ERR_QUERY_ERR"}
		return r, err
	}

	now := time.Now().Unix()
	// 登录过期
	if (now-sInfo.CreateTime.Unix())/86400 > int64(lduration) {
		err = common.ServerError{common.RETURN_CODE_MA_AUTH_ERR, "AUTH_ERR_LOGIN_EXPIRED"}
		return r, err
	}

	// session过期
	if now-sInfo.LastVisitTime.Unix() > int64(sduration) {
		err = common.ServerError{common.RETURN_CODE_MA_AUTH_ERR, "AUTH_ERR_SESSION_EXPIRED"}
		return r, err
	}

	// 解析userinfo 信息
	userInfoStr, err := base64.StdEncoding.DecodeString(sInfo.UserInfo)
	if err != nil {
		beego.Error("base64 decode userinfo string error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_AUTH_ERR, "AUTH_ERR_DECRYPT_USER_ERR"}
		return r, err
	}

	var userInfo common.UserInfo
	err = json.Unmarshal(userInfoStr, &userInfo)
	if err != nil {
		beego.Error("json decode userinfo error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_AUTH_ERR, "AUTH_ERR_JSON_DECODE_ERR"}
		return r, err
	}

	// 更新最后登录时间
	sInfo.LastVisitTime = time.Now()
	if _, err := o.Update(&sInfo); err != nil {
		beego.Error("update last visit time to db error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "DB_ERR_UPDATE_LASTVISITTIME"}
		return r, err
	}

	r.Id = uuid
	r.Skey = skey
	r.Duration = int(now - sInfo.LastVisitTime.Unix())
	r.UserInfo = userInfo

	return
}

// 初始化app信息
func InitApp(appid string, secret string) (r common.ResponseData, err error) {

	if appid == "" {
		err = common.ServerError{common.RETURN_CODE_MA_NO_PARA, "PARA_ERR_APPID_EMPTY"}
		return r, err
	}

	if secret == "" {
		err = common.ServerError{common.RETURN_CODE_MA_NO_PARA, "PARA_ERR_SECRET_EMPTY"}
		return r, err
	}

	o := orm.NewOrm()
	qs := o.QueryTable(new(CAppinfo))
	appExist := qs.Filter("Appid", appid).Exist()
	if appExist {
		err = common.ServerError{common.RETURN_CODE_MA_INIT_APPINFO_ERR, "APPID_ALREADY_EXISTS"}
		return r, err
	}

	var appinfo CAppinfo
	appinfo.Appid = appid
	appinfo.Secret = secret
	appinfo.LoginDuration = 30
	appinfo.SessionDuration = 30 * 86400
	appinfo.QcloudAppid = "appid_qcloud"
	appinfo.Ip = "0.0.0.0"

	id, err := o.Insert(&appinfo)
	if err != nil {
		beego.Error("insert new appinfo to mysql error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "DB_ERR_INSERT_APPINFO"}
		return r, err
	}
	beego.Debug("add new appinfo to db, id: ", id)

	return
}
