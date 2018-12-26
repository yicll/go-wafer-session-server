package common

import (
	"fmt"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	Error error
}

// 全局返回结果结构体
type RenderResult struct {
	Code          RETURN_CODE `json:"returnCode"`
	Message       string      `json:"returnMessage"`
	Data          interface{} `json:"returnData"`
	Version       int         `json:"version"`
	ComponentName string      `json:"componentName"`
}

// 接口请求接口信息
type RequestInterface struct {
	InterfaceName string            `json:"interfaceName"`
	Appid         string            `json:"appid"`
	Para          map[string]string `json:"para"`
}

// 接口请求数据
type RequestData struct {
	Version       int              `json:"version"`
	ComponentName string           `json:"componentName"`
	Interface     RequestInterface `json:"interface"`
}

// 解析成功之后的用户信息结构
type UserInfo struct {
	OpenId    string `json:"openId"`
	UnionId   string `json:"unionId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	Language  string `json:"language"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
}

// 接口返回数据结构体
type ResponseData struct {
	Id       string   `json:"id"`
	Skey     string   `json:"skey"`
	UserInfo UserInfo `json:"user_info"`
	Duration int      `json:"duration"`
}

// 校验参数是否存在
func (this *BaseController) CheckParams(keys []string, data map[string]string) (err error) {
	for _, key := range keys {
		if ret, ok := data[key]; !ok || ret == "" {
			beego.Warning(fmt.Sprintf("request param %s is empty", key))
			err = ServerError{RETURN_CODE_MA_PARA_ERR, "PARA_ERR"}
			return err
		}
	}

	return
}

// 渲染结果
func (this *BaseController) RenderJSON(data interface{}) {

	var (
		code    RETURN_CODE
		message string
	)

	if this.Error == nil {
		code = RETURN_CODE_MA_OK
		message = "SUCCESS"
	} else {
		code = this.Error.(ServerError).GetCode()
		message = this.Error.(ServerError).GetMessage()
	}

	res := RenderResult{code, message, data, 1, "MA"}
	this.Data["json"] = res
	this.ServeJSON()
}
