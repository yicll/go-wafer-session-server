package redis_models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"wafer-session-server/common"
	"wafer-session-server/redis_helper"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

// redis存储app信息结构
type RAppinfo struct {
	Appid           string `redis:"appid"`
	Secret          string `redis:"secret"`
	LoginDuration   int    `redis:"ld"`
	SessionDuration int    `redis:"sd"`
	Ip              string `redis:"ip"`
}

// redis存储session信息结构
type RSessionInfo struct {
	Appid         string `redis:"-"`
	Uuid          string `redis:"uuid"`
	Skey          string `redis:"skey"`
	CreateTime    int64  `redis:"ct"`
	LastVisitTime int64  `redis:"lvt"`
	OpenId        string `redis:"openid"`
	SessionKey    string `redis:"sesskey"`
	UserInfo      string `redis:"userinfo"`
}

// 获取appinfo存在在redis中的key
func GetAppinfoKey(appid string) string {
	return fmt.Sprintf("app_%s", appid)
}

func GetSessionKey(appid string, uuid string) string {
	return fmt.Sprintf("sess_%s_%s", appid, uuid)
}

// 获取appinfo信息
func GetAppinfById(appid string) (r RAppinfo, err error) {
	appKey := GetAppinfoKey(appid)

	conn := redis_helper.RedisPool.Get()
	defer conn.Close()

	appExist, err := redis.Bool(conn.Do("EXISTS", appKey))
	if err != nil {
		beego.Error("get info from redis error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "DB_ERR_REDIS_EXISTS"}
		return r, err
	}

	if appExist {
		v, err := redis.Values(conn.Do("HGETALL", appKey))
		if err != nil {
			beego.Error("get appinfo from redis error, message: ", err.Error())
			err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "DB_ERR_REDIS_APPINFO"}
			return r, err
		}

		if err := redis.ScanStruct(v, &r); err != nil {
			beego.Error("scan redis struct to appinfo struct error, message: ", err.Error())
			err = common.ServerError{common.RETURN_CODE_MA_SERVER_ERR, "SERVER_ERR_SCAN_STRUCT"}
			return r, err
		}

		beego.Informational(fmt.Sprintf("result app info: %+v", r))

		return r, nil
	} else if appid == "" {
		// 兼容wafer sdk不传appid的情况
		keys, err := redis.Strings(conn.Do("KEYS", "app_*"))
		if err != nil {
			beego.Error("get keys from redis error, message: ", err.Error())
			err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "DB_ERR_REDIS_KEYS"}
			return r, err
		}

		if len(keys) == 0 || keys == nil {
			beego.Error("redis app info is empty")
			err = common.ServerError{common.RETURN_CODE_MA_NO_APPID, "REDIS_NOT_EXISTS_APPID"}
			return r, err
		} else {
			appKey = keys[0]
		}

		v, err := redis.Values(conn.Do("HGETALL", appKey))
		if err != nil {
			beego.Error("get appinfo from redis error, message: ", err.Error())
			err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "DB_ERR_REDIS_APPINFO"}
			return r, err
		}

		if err := redis.ScanStruct(v, &r); err != nil {
			beego.Error("scan redis struct to appinfo struct error, message: ", err.Error())
			err = common.ServerError{common.RETURN_CODE_MA_SERVER_ERR, "SERVER_ERR_SCAN_STRUCT"}
			return r, err
		}

		beego.Informational(fmt.Sprintf("result app info: %+v", r))

		return r, nil
	} else {
		beego.Informational("redis not exist appid: ", appid)
		err = common.ServerError{common.RETURN_CODE_MA_NO_APPID, "REDIS_NOT_EXISTS_APPID"}
		return r, err
	}

	return
}

// 保存session
func SaveSessionInfo(info *RSessionInfo, expire int) (err error) {

	beego.Debug(fmt.Sprintf("redis session info: %+v", info))

	sessKey := GetSessionKey(info.Appid, info.Uuid)

	conn := redis_helper.RedisPool.Get()
	defer conn.Close()

	sessExist, err := redis.Bool(conn.Do("EXISTS", sessKey))
	if err != nil {
		beego.Error("check key from redis error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "DB_ERR_REDIS_EXISTS"}
		return err
	}

	// 存在key，则删除重新插入
	if sessExist {
		_, err := conn.Do("DEL", sessKey)
		if err != nil {
			beego.Error("delete key from redis error, message: ", err.Error())
			err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "DB_ERR_DEL_SESS"}
			return err
		}
	}

	if _, err := conn.Do("HMSET", redis.Args{}.Add(sessKey).AddFlat(info)...); err != nil {
		beego.Error("set key to redis error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "DB_ERR_SET_SESS"}
		return err
	}

	if expire == 0 {
		expire = 30 * 86400
	}

	if _, err := conn.Do("EXPIRE", sessKey, expire); err != nil {
		beego.Error(fmt.Sprintf("set key[%s] expire[%d] to redis error, message: %s", sessKey, expire, err.Error()))
		err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "DB_ERR_SET_SESS_EXPIRE"}
		return err
	}

	return nil
}

// 检查登录信息
func CheckSession(appid string, uuid string, skey string, lduration int, sduration int) (r common.ResponseData, err error) {

	sessKey := GetSessionKey(appid, uuid)

	conn := redis_helper.RedisPool.Get()
	defer conn.Close()

	sessExist, err := redis.Bool(conn.Do("EXISTS", sessKey))
	if err != nil {
		beego.Error("check key exist from redis error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "DB_ERR_REDIS_EXISTS"}
		return r, err
	}

	// 登录信息不存在
	if !sessExist {
		err = common.ServerError{common.RETURN_CODE_MA_AUTH_ERR, "AUTH_ERR_NO_EXIST"}
		return r, err
	}

	var sess RSessionInfo
	v, err := redis.Values(conn.Do("HGETALL", sessKey))
	if err != nil {
		beego.Error("redis command hgetall error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "DB_ERR_REDIS_SESSINFO"}
		return r, err
	}

	if err := redis.ScanStruct(v, &sess); err != nil {
		beego.Error("scan redis struct to sessinfo struct error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_SERVER_ERR, "SERVER_ERR_SCAN_STRUCT"}
		return r, err
	}
	beego.Informational(fmt.Sprintf("result session info: %+v", sess))

	now := time.Now().Unix()
	// 登录过期
	if (now-sess.CreateTime)/86400 > int64(lduration) {
		err = common.ServerError{common.RETURN_CODE_MA_AUTH_ERR, "AUTH_ERR_LOGIN_EXPIRED"}
		return r, err
	}

	// session过期
	if now-sess.LastVisitTime > int64(sduration) {
		err = common.ServerError{common.RETURN_CODE_MA_AUTH_ERR, "AUTH_ERR_SESSION_EXPIRED"}
		return r, err
	}

	// 解析user info信息
	userInfoStr, err := base64.StdEncoding.DecodeString(sess.UserInfo)
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
	if _, err := conn.Do("hset", sessKey, "lvt", time.Now().Unix()); err != nil {
		beego.Error("redis command hset error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "REDIS_ERR_UPDATE_LASTVISITTIME"}
		return r, err
	}

	r.Id = uuid
	r.Skey = skey
	r.Duration = int(now - sess.LastVisitTime)
	r.UserInfo = userInfo

	return
}

func InitApp(appid string, secret string) (r common.ResponseData, err error) {

	if appid == "" {
		err = common.ServerError{common.RETURN_CODE_MA_NO_PARA, "PARA_ERR_APPID_EMPTY"}
		return r, err
	}

	if secret == "" {
		err = common.ServerError{common.RETURN_CODE_MA_NO_PARA, "PARA_ERR_SECRET_EMPTY"}
		return r, err
	}

	appKey := GetAppinfoKey(appid)

	conn := redis_helper.RedisPool.Get()
	defer conn.Close()

	appExist, err := redis.Bool(conn.Do("EXISTS", appKey))
	if err != nil {
		beego.Error("get info from redis error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "DB_ERR_REDIS_EXISTS"}
		return r, err
	}

	if appExist {
		err = common.ServerError{common.RETURN_CODE_MA_INIT_APPINFO_ERR, "APPID_ALREADY_EXISTS"}
		return r, err
	}

	var appinfo RAppinfo
	appinfo.Appid = appid
	appinfo.Secret = secret
	appinfo.LoginDuration = 30
	appinfo.SessionDuration = 30 * 86400
	appinfo.Ip = "0.0.0.0"

	if _, err := conn.Do("HMSET", redis.Args{}.Add(appKey).AddFlat(&appinfo)...); err != nil {
		beego.Error("add appinfo to redis error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_DB_ERR, "DB_ERR_ADD_APPINFO"}
		return r, err
	}

	return
}
