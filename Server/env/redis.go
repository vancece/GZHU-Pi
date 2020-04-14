package env

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
	_ "golang.org/x/sync/singleflight"
)

var RedisCli *redis.Client

func InitRedis() (err error) {
	host := Conf.Redis.Host
	port := Conf.Redis.Port
	password := Conf.Redis.Password

	RedisCli = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       0,
	})

	logs.Info("Ping Redis: %s:%d", host, port)
	_, err = RedisCli.Ping().Result()
	if err != nil {
		logs.Error("init redis failed: ", err)
		return
	}
	logs.Info(fmt.Sprintf("connect to redis %s:%d", host, port))

	return
}

