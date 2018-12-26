package common

type RETURN_CODE int

const (
	RETURN_CODE_MA_OK                       RETURN_CODE = 0     // 成功返回码
	RETURN_CODE_MA_DB_ERR                   RETURN_CODE = 1001  // DB错误等
	RETURN_CODE_MA_NO_INTERFACE             RETURN_CODE = 1002  // 接口参数不存在
	RETURN_CODE_MA_PARA_ERR                 RETURN_CODE = 1003  // 参数错误
	RETURN_CODE_MA_WEIXIN_NET_ERR           RETURN_CODE = 1005  // 连接微信服务器失败
	RETURN_CODE_MA_CHANGE_SESSION_ERR       RETURN_CODE = 1006  // 新增修改SESSION失败
	RETURN_CODE_MA_WEIXIN_RETURN_ERR        RETURN_CODE = 1007  // 微信返回值错误
	RETURN_CODE_MA_UPDATE_LASTVISITTIME_ERR RETURN_CODE = 1008  // 更新最近访问时间失败
	RETURN_CODE_MA_REQUEST_ERR              RETURN_CODE = 1009  // 请求包不是json
	RETURN_CODE_MA_INTERFACE_ERR            RETURN_CODE = 1010  // 接口名称错误
	RETURN_CODE_MA_NO_PARA                  RETURN_CODE = 1011  // 不存在参数
	RETURN_CODE_MA_NO_APPID                 RETURN_CODE = 1012  // 不能获取AppID
	RETURN_CODE_MA_INIT_APPINFO_ERR         RETURN_CODE = 1013  // 初始化AppID失败
	RETURN_CODE_MA_SERVER_ERR               RETURN_CODE = 2000  // 服务器处理错误
	RETURN_CODE_MA_WEIXIN_CODE_ERR          RETURN_CODE = 40029 // CODE无效
	RETURN_CODE_MA_AUTH_ERR                 RETURN_CODE = 60012 // 鉴权失败
	RETURN_CODE_MA_DECRYPT_ERR              RETURN_CODE = 60021 // 解密失败
)
