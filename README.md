Golang Wafer 会话服务器
===============

本项目是 [Wafer](https://github.com/tencentyun/wafer) 组成部分中的wafer-session-server的升级版本

使用golang开发，支持`redis`和`mysql`，提供会话服务供 SDK 或独立使用。

会话服务的实现细请参考 [Wiki](https://github.com/tencentyun/wafer/wiki/%E4%BC%9A%E8%AF%9D%E6%9C%8D%E5%8A%A1)。


## 接口协议

### 请求

会话服务器提供 HTTP 接口来实现会话管理，下面是协议说明。

* 协议类型：`HTTP`
* 传输方式：`POST`
* 编码类型：`UTF-8`
* 编码格式：`JSON`

请求示例：

```http
POST / HTTP/1.1
Content-Type: application/json;charset=utf-8

{
    "version": 1,
    "componentName": "MA",
    "interface": {
        "interfaceName": "qcloud.cam.id_skey",
		"appid": your appid", //该字段是在原生的wafer session Server api上新增的，不传的话代表只支持单server
        "para": { "code": "...", "encrypt_data": "..." }
    }
}
```

### 响应

HTTP 输出为响应内容，下面是响应内容说明：

* 内容编码：`UTF-8`
* 内容格式：`JSON`

响应示例：

```json
{
    "returnCode": 0,
    "returnMessage": "OK",
    "returnData": {
    	"id": "session info uuid", // 本次登录的session标识uuid
    	"skey": "session info skey", // 本次登录的session标识skey，建议和id一起组合成session登录唯一标识
		"user_info": {
			"openId": "oFRD80IWuaryqCHQmcHV5jKpCHYo",
			"uionId": "o3L7v1GhRiwsgvcaFGSVu5cbtOWo",
			"nickName": "yicll",
			"gender": 1,
			"language": "zh_CN",
			"city": "Fangshan",
			"province": "Beijing",
			"country": "China",
			"avatarUrl": "https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTL1F9UiaAlLTeJRAIXfmyUxWj485QI08dWR9eU16qKvfVediaYZNDOBUiars3icjJSW7fxJDlO8ymWYBA/132",
		}
    }
}
```

* `returnCode` - 返回码，如果成功则取值为 `0`，如果失败则取值为具体错误码；
* `returnMessage` - 如果返回码非零，内容为出错信息；
* `returnData` - 返回的数据

### qcloud.cam.id_skey

`qcloud.cam.id_skey` 处理用户登录请求。

使用示例：

```sh
curl -i -d'{"version":1,"componentName":"MA","interface":{"interfaceName":"qcloud.cam.id_skey","appid":"wxe325db015fc632af","para":{"code":"001EWYiD1CVtKg0jXGjD1e6WiD1EWYiC","encrypt_data":"DNlJKYA0mJ3+RDXD/syznaLVLlaF4drGzeZvJFmjnEKtOAi37kAzC/1tCBr7KqGX8EpiLuWl8qt/kcH9a4LxDC5LQvlRLJlDogTEIwtlT/2jBWBuWwBC3vWFhm7Uuq5AOLZV+xG9UmWPKECDZX9UZpWcPRGQpiY8OOUNBAywVniJv6rC2eADFimdRR2qPiebdC3cry7QAvgvttt1Wk56Nb/1TmIbtJRTay5wb+6AY1H7AT1xPoB6XAXW3RqODXtRR0hZT1s/o5y209Vcc6EBal5QdsbJroXa020ZSD62EnlrOwgYnXy5c8SO+bzNAfRw59SVbI4wUNYz6kJb4NDn+y9dlASRjlt8Rau4xTQS+fZSi8HHUwkwE6RRak3qo8YZ7FWWbN2uwUKgQNlc/MfAfLRcfQw4XUqIdn9lxtRblaY="}}}' http://127.0.0.1:8888
```

响应数据：

* `id` - 会话 id
* `skey` - 会话 skey
* `userInfo` - 用户信息

### qcloud.cam.auth

使用 `qcloud.cam.auth` 接口检查用户登录态。

响应数据：

* 同前一个接口，根据returnCode判断是否登录成功，0表示成功，将返回用户信息，否则根据提示引导用户重新授权登录

### 错误码
<table>
  <tbody>
  <tr>
    <th>错误码</th>
    <th>解释</th>
  </tr>
  <tr>
    <td>0</td>
    <td>成功</td>
  </tr>
  <tr>
    <td>1001</td>
    <td>数据库错误</td>
  </tr>
   <tr>
    <td>1002</td>
    <td>接口不存在</td>
  </tr>
  <tr>
    <td>1003</td>
    <td>参数错误</td>
  </tr>
  <tr>
    <td>1005</td>
    <td>连接微信服务器失败</td>
  </tr>
   <tr>
    <td>1006</td>
    <td>新增或修改 SESSION 失败</td>
  </tr>
  <tr>
    <td>1007</td>
    <td>微信返回值错误</td>
  </tr>
  <tr>
    <td>1008</td>
    <td>更新最近访问时间失败</td>
  </tr>
  <tr>
    <td>1009</td>
    <td>请求包不是json</td>
  </tr>
  <tr>
    <td>1010</td>
    <td>接口名称错误</td>
  </tr>
  <tr>
    <td>1011</td>
    <td>参数不存在</td>
  </tr>
   <tr>
    <td>1012</td>
    <td>不能获取 AppID</td>
  </tr>
  <tr>
    <td>1013</td>
    <td>初始化 AppID 失败</td>
  </tr>
  <tr>
    <td>2000</td>
    <td>server处理出现错误</td>
  </tr>
  <tr>
    <td>40029</td>
    <td>CODE 无效</td>
  </tr>
  <tr>
    <td>60021</td>
    <td>解密失败</td>
  </tr>
  <tr>
    <td>60012</td>
    <td>鉴权失败</td>
  </tr>
  </tbody>
</table>

## 项目配置文件解析

配置文件位置在src/wafer-session-server/conf/app.conf

具体配置信息及其含义如下

```
appname = wafer-session-server 		// app名称，建议不要修改
httpport = 8888						// 服务监听端口，当dev、prod节点里没有配置该项的时候，取该端口为默认值
runmode = dev						// 当前运行环境，dev：开发环境、prod：生产环境
autorender = false					// 自动渲染模板，api项目，该项必须为false
copyrequestbody = true				// 默认为true，不允许更改
EnableDocs = true					// 是否开始自动化文档功能

[dev]								// 当rundmode=dev的时候，读取该节点的配置信息
httpport = 8888						// 服务监听端口
dbdriver = redis					// 当前服务db引擎，mysql、redis
									// 当dbdriver=mysql时，下面mysql.*配置必须配置正确
mysql.host = 127.0.0.1				// mysql数据库地址
mysql.port = 3306					// mysql数据库端口
mysql.user = root					// mysql用户名
mysql.password = 123456				// mysql密码
mysql.db = cAuth					// mysql数据库名称
mysql.debug = true					// 是否开启debug模式，开启的话，将会在stdout中输出sql信息
									// 当dbdriver=redis时，下面的redis.*配置必须配置正确
redis.host = 127.0.0.1				// redis连接地址
redis.port = 6379					// redis连接端口
redis.password = 					// redis连接密码，不需要密码的话，该配置为空
redis.db = 1						// redis连接db
redis.maxidle = 10					// redis连接池允许最大空闲连接数
redis.maxactive = 1000				// redis连接池允许最大连接数
redis.timeout = 30					// 空闲连接最大允许留存超时时间


[prod]								// 当rundmode=prod的时候，读取该节点的配置信息
httpport = 8889						// 服务监听端口
dbdriver = redis					// 当前服务db引擎，mysql、redis
									// 当dbdriver=mysql时，下面mysql.*配置必须配置正确
mysql.host = 127.0.0.1				// mysql数据库地址
mysql.port = 3306					// mysql数据库端口
mysql.user = root					// mysql用户名
mysql.password = 123456				// mysql密码
mysql.db = cAuth					// mysql数据库名称
mysql.debug = true					// 是否开启debug模式，开启的话，将会在stdout中输出sql信息
									// 当dbdriver=redis时，下面的redis.*配置必须配置正确
redis.host = 127.0.0.1				// redis连接地址
redis.port = 6379					// redis连接端口
redis.password = 					// redis连接密码，不需要密码的话，该配置为空
redis.db = 1						// redis连接db
redis.maxidle = 10					// redis连接池允许最大空闲连接数
redis.maxactive = 1000				// redis连接池允许最大连接数
redis.timeout = 30					// 空闲连接最大允许留存超时时间
```


## 搭建会话管理服务器

### 环境准备

```
golang version: 1.9.2
```

golang环境请自行搭建，最小需要版本为1.9.2，GOPATH、GOROOT等环境变量请自己设置好

### 数据库

根据项目选择，redis、mysql数据库请自行部署，部署成功之后修改配置文件里的相关配置

如果选择mysql的话，需要使用sql语句进行建表，详见[db.sql](http://git.domob-inc.cn/mp-lib/go-wafer-session-server/blob/master/src/wafer-session-server/sql/db.sql)

### 项目编译打包

将项目下载到本地

```sh
git clone http://git.domob-inc.cn/mp-lib/go-wafer-session-server wafer-session-server
cd wafer-session-server
./build.sh
#target目录下的wafer-session-server.tar.gz就是用于上线部署的压缩包
```

### 服务启动

将打包成功的tar包拷贝到需要上线的服务器，例如你将tar包拷贝到了服务器的/tmp目录下

```sh
mkdir wafer-session-server
cd wafer-session-server
cp /tmp/wafer-session-server.tar.gz .
tar zxvf wafer-session-server.tar.gz
./control.sh start
```
   
## 初始化 appId 和 appSecret

### API初始化

直接请求服务API进行appid和secret初始化

```sh
curl -i -d'{"version":1,"componentName":"MA","interface":{"interfaceName":"qcloud.cam.initapp","appid":"[替换成你的APPID]","para":{"secret":"[替换成你的APP SECRET]"}}}' htt://127.0.0.1:8888
```

### mysql

登录到 MySql 后，手动插入配置到 `cAuth` 表中。

```
use cAuth;
insert into cAppinfo set appid='Your appid',secret='Your secret';
```

### redis

登录到redis后，手动插入app key

```sh
redis-cli -h 127.0.0.1 -p 6379
> select 1
> hmset app_[your appid] appid [your appid] secret [your sceret] ld 30 sd 2592000 ip 0.0.0.0
```
    
### 测试服务可用性

```sh
sh
curl -i -d'{"version":1,"componentName":"MA","interface":{"interfaceName":"qcloud.cam.id_skey","appid":"wxe325db015fc632af","para":{"code":"001btntB1UQJ7d0VVasB1ZKGtB1btnth","encrypt_data":"qs8afGiRlAsjIcNuG9CqxMbMgr6tpaTqOrpa9szUSrYfObQR54ThGhmAadEhkuW/6Flyqa+r+p/4BuKnCLx81TzwqM+7gP3pdOG4rLvlvWCtDes2blsGZm2wNFOqqwj+xfVQqj25JznX75lNbObY5Ic67ZTiaszMzJym0QDy7vaBQMCwdGLfTiVPc35cpfq9ZZzGDVVewHoNGauhPrkOxdu+ec/M6/Fp39J32yEyfi/7lkUwauobdDl7ovazjoFGvfeBOjdXlmBGuF0+W5KKjdsXINLHWL1m4gZD5twLQxICC4A6W6YvXoLAHr41eslvfFvGptIJFOW4GXnEZyhzc7tgubiSvMy9cMA0NcB6o8qIh7GrZ1sp6FSdrCDaDj3zXNlHzgbXvNfX/Q7PkQ18AaofjapSnoEOUxfiHwR/yNpK05yqviCgY7UdoNUSKd3GtMXg+KJTG5yfvOfN23JiaQnJ4P30wJ15IJb07pQsEMk0C6QthDfPvRnxpU07ERgGL7FKmP1f3Z2HlrzET/Z2Jw==","iv":"nU6TJmoVfrz8Vt8FJbrZYA=="}}}' http://127.0.0.1:8888
```
    
    
## mysql数据库设计

当项目选择使用mysql为存储，则下面是mysql的数据库详细设计，如果选择redis为存储，请忽略

全局信息表 `cAppInfo` 保存会话服务所需要的配置项。

<table>
  <tbody>
  <tr>
    <th>Field</th>
    <th>Type</th>
    <th>Null</th>
    <th>key</th>
    <th>Extra</th>
  </tr>
  <tr>
    <td>appid</td>
    <td>varchar(200)</td>
    <td>NO</td>
    <td>PRI</td>
    <td>申请微信小程序开发者时，微信分配的 appId</td>
  </tr>
  <tr>
    <td>secret</td>
    <td>varchar(300)</td>
    <td>NO</td>
    <td></td>
    <td>申请微信小程序开发者时，微信分配的 appSecret</td>
  </tr>
  <tr>
    <td>login_duration</td>
    <td>int(11)</td>
    <td>NO</td>
    <td></td>
    <td>登录过期时间，单位为天，默认 30 天</td>
  </tr>
  <tr>
    <td>session_duration</td>
    <td>int(11)</td>
    <td>NO</td>
    <td></td>
    <td>会话过期时间，单位为秒，默认为 2592000 秒(即30天)</td>
  </tr>
  
  </tbody>
</table>
    

会话记录 `cSessionInfo` 保存每个会话的数据。

<table>
  <tbody>
  <tr>
    <th>Field</th>
    <th>Type</th>
    <th>Null</th>
    <th>key</th>
    <th>Extra</th>
  </tr>
  <tr>
    <td>id</td>
    <td>int(11)</td>
    <td>NO</td>
    <td>MUL</td>
    <td>会话 ID（自增长）</td>
  </tr>
  <tr>
    <td>appid</td>
    <td>varchar(100)</td>
    <td>NO</td>
    <td></td>
    <td>session归属的appid</td>
  </tr>
   <tr>
    <td>uuid</td>
    <td>varchar(100)</td>
    <td>NO</td>
    <td></td>
    <td>会话 uuid</td>
  </tr>
  <tr>
    <td>skey</td>
    <td>varchar(100)</td>
    <td>NO</td>
    <td></td>
    <td>会话 Skey</td>
  </tr>
  <tr>
    <td>create_time</td>
    <td>datetime</td>
    <td>NO</td>
    <td></td>
    <td>会话创建时间，用于判断会话对应的 open_id 和 session_key 是否过期（是否超过 `cAppInfo` 表中字段 `login_duration` 配置的天数）</td>
  </tr>
  <tr>
    <td>last_visit_time</td>
    <td>datetime</td>
    <td>NO</td>
    <td></td>
    <td>最近访问时间，用于判断会话是否过期（是否超过 `cAppInfo` 表中字段 `session_duration` 的配置的秒数）</td>
  </tr>
  <tr>
    <td>open_id</td>
    <td>varchar(100)</td>
    <td>NO</td>
    <td>MUL</td>
    <td>微信服务端返回的 `open_id` 值 </td>
  </tr>
  <tr>
    <td>session_key</td>
    <td>varchar(100)</td>
    <td>NO</td>
    <td></td>
    <td>微信服务端返回的 `session_key` 值 </td>
  </tr>
  <tr>
    <td>user_info</td>
    <td>varchar(2048)</td>
    <td>YES</td>
    <td></td>
    <td>已解密的用户数据</td>
  </tr>
  </tbody>
</table>

建数据库的详细 SQL 脚本请参考 [db.sql](http://git.domob-inc.cn/mp-lib/go-wafer-session-server/blob/master/src/wafer-session-server/sql/db.sql)

## 老版本wafer-session-server数据迁移

如果选择mysql为存储，可以直接复用原来的db结构，需要更改cSessionInfo表，修改如下

```
alter table cSessionInfo add appid varchar(200) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL default '';
update cSessionInfo set appid=[Your appid];
```