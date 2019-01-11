package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"wafer-session-server/common"
	"wafer-session-server/common/utils"
	"wafer-session-server/models"

	"github.com/astaxie/beego"
	"github.com/xlstudio/wxbizdatacrypt"
)

type MainController struct {
	common.BaseController
}

func (this *MainController) Get() {
	data, err := models.GetAppinfoById("")

	if err != nil {
		this.Error = err
		this.RenderJSON(nil)
	}

	this.RenderJSON(data)
}

func (this *MainController) Post() {
	t1 := time.Now().UnixNano()

	data, err := this.parseRequest(this.Ctx.Input.RequestBody)

	t2 := time.Now().UnixNano()
	beego.Informational(fmt.Sprintf("server exec time: %d ms", (t2-t1)/1e6))

	if err != nil {
		beego.Warning("response has error, info: ", err.Error())
		this.Error = err
		this.RenderJSON(nil)
	} else {
		beego.Informational(fmt.Sprintf("response data: %+v", data))
		this.RenderJSON(data)
	}
}

func (this *MainController) getIdSkey(appid string, code string, encryptData string, iv string) (r common.ResponseData, err error) {
	appinfo, err := models.GetAppinfoById(appid)
	if err != nil {
		return
	}
	appid = appinfo.Appid

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appinfo.Appid, appinfo.Secret, code)
	req, err := http.Get(url)
	if err != nil {
		beego.Error("request url: ", url, "error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_WEIXIN_NET_ERR, "MA_WEIXIN_NET_ERR"}
		return
	}

	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		beego.Error("read reponse body error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_WEIXIN_NET_ERR, "MA_WEIXIN_NET_ERR"}
		return
	}

	var resp models.WxResponseData
	//err = req.ToJSON(&resp)
	err = json.Unmarshal(body, &resp)
	if err != nil {
		beego.Error("unmarshal json to struct error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_WEIXIN_NET_ERR, "MA_WEIXIN_NET_ERR"}
		return
	}

	beego.Informational(fmt.Sprintf("request url: %s, response result: %+v", url, resp))

	// code无效
	if resp.Errcode == 40029 {
		err = common.ServerError{common.RETURN_CODE_MA_WEIXIN_CODE_ERR, "MA_WEIXIN_CODE_ERR"}
		return
	}

	// 其他错误
	if resp.Errcode != 0 {
		err = common.ServerError{common.RETURN_CODE_MA_AUTH_ERR, "WEIXIN_AUTH_ERR"}
		return
	}

	// 返回信息必须有openid和session_key
	if resp.Openid == "" || resp.SessionKey == "" {
		err = common.ServerError{common.RETURN_CODE_MA_WEIXIN_RETURN_ERR, "WEIXIN_RETURN_ERR"}
		return
	}

	sessionKey := resp.SessionKey

	// 解析加密数据
	pc := wxbizdatacrypt.WxBizDataCrypt{AppID: appinfo.Appid, SessionKey: sessionKey}
	decryptData, err := pc.Decrypt(encryptData, iv, true)
	if err != nil {
		beego.Error("decrypt data error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_DECRYPT_ERR, "DECRYPT_FAIL"}
		return
	}

	beego.Debug("decrpyt data: ", decryptData.(string))
	var userInfo common.UserInfo
	err = json.Unmarshal([]byte(decryptData.(string)), &userInfo)
	if err != nil {
		beego.Error("json decode decrypt data error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_SERVER_ERR, "DECRYPT_JSON_FAIL"}
		return
	}

	beego.Debug(fmt.Sprintf("decrypt data success, data: %+v", userInfo))

	userInfoJson, err := json.Marshal(userInfo)
	if err != nil {
		beego.Error("json encode userInfo error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_SERVER_ERR, "JSON_ENCODE_FAIL"}
		return
	}

	sessionInfo := models.SessionInfo{
		Appid:         appinfo.Appid,
		Uuid:          utils.GenerateUuid(),
		Skey:          utils.GenerateSkey(),
		CreateTime:    time.Now(),
		LastVisitTime: time.Now(),
		OpenId:        userInfo.OpenId,
		SessionKey:    sessionKey,
		UserInfo:      base64.StdEncoding.EncodeToString(userInfoJson)}

	err = models.SaveSessionInfo(appinfo, &sessionInfo)
	if err != nil {
		beego.Error("save session info error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_SERVER_ERR, "SAVE_SESSION_INFO_ERR"}
		return
	}

	var data common.ResponseData
	data.Id = sessionInfo.Uuid
	data.Skey = sessionInfo.Skey
	data.UserInfo = userInfo
	data.Duration = appinfo.SessionDuration

	return data, nil
}

// 校验登录态
func (this *MainController) checkAuth(appid string, uuid string, skey string) (r common.ResponseData, err error) {
	appinfo, err := models.GetAppinfoById(appid)
	if err != nil {
		return
	}

	return models.CheckSession(uuid, skey, appinfo)
}

func (this *MainController) initApp(appid string, secret string) (r common.ResponseData, err error) {
	beego.Debug(fmt.Sprintf("appid: %s, secret: %s", appid, secret))
	return models.InitApp(appid, secret)
}

// 解析请求数据并格式化
func (this *MainController) parseRequest(body []byte) (r interface{}, err error) {

	beego.Informational("request body: ", string(body))

	var data common.RequestData
	if err = json.Unmarshal(body, &data); err != nil {
		beego.Error("json decode reqeust body error, message: ", err.Error())
		err = common.ServerError{common.RETURN_CODE_MA_REQUEST_ERR, "REQUEST_IS_NOT_JSON"}
		return
	}

	beego.Informational("interface name: ", data.Interface.InterfaceName)
	if data.Interface.InterfaceName == "" {
		err = common.ServerError{common.RETURN_CODE_MA_NO_INTERFACE, "NO_INTERFACENAME_PARA"}
		return
	}

	if data.Interface.Para == nil {
		err = common.ServerError{common.RETURN_CODE_MA_PARA_ERR, "NO_PARA"}
		return
	}

	switch data.Interface.InterfaceName {
	case "qcloud.cam.id_skey":
		ks := []string{"code", "encrypt_data", "iv"}
		if err = this.CheckParams(ks, data.Interface.Para); err != nil {
			return
		}
		return this.getIdSkey(data.Interface.Appid,
			data.Interface.Para["code"],
			data.Interface.Para["encrypt_data"],
			data.Interface.Para["iv"])
	case "qcloud.cam.auth":
		ks := []string{"id", "skey"}
		if err = this.CheckParams(ks, data.Interface.Para); err != nil {
			return
		}
		return this.checkAuth(data.Interface.Appid,
			data.Interface.Para["id"],
			data.Interface.Para["skey"])
	case "qcloud.cam.decrypt":
		ks := []string{"id", "skey", "encrypt_data"}
		if err = this.CheckParams(ks, data.Interface.Para); err != nil {
			return
		}
	case "qcloud.cam.initapp":
		ks := []string{"secret"}
		if err = this.CheckParams(ks, data.Interface.Para); err != nil {
			return
		}
		return this.initApp(data.Interface.Appid, data.Interface.Para["secret"])
	default:
		err = common.ServerError{common.RETURN_CODE_MA_INTERFACE_ERR, "INTERFACENAME_PARA_ERR"}
		return
	}

	return data, nil
}
