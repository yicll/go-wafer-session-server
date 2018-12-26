package redis_helper

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

var RedisPool *redis.Pool

// 初始化连接池
// Params:
// host: connection host
// port: connection port
// db: select redis db
// mi: MaxIdle 最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
// ma: MaxActive 最大的激活连接数，表示同时最多有N个连接 ，为0事表示没有限制
// timeout: IdleTimeout 最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
func NewReidsPool(host string, port int, pwd string, db int, mi int, ma int, timeout time.Duration) (err error) {
	connstr := fmt.Sprintf("%s:%d", host, port)

	RedisPool = &redis.Pool{
		MaxIdle:     mi,
		MaxActive:   ma,
		IdleTimeout: timeout,
		Dial: func() (c redis.Conn, err error) {
			if pwd == "" {
				c, err = redis.Dial("tcp", connstr, redis.DialDatabase(db))
			} else {
				c, err = redis.Dial("tcp", connstr,
					redis.DialDatabase(db), redis.DialPassword(pwd))
			}

			if err != nil {
				c.Close()
				beego.Error("connect to redis fail, message: ", err.Error())
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}
