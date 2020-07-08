package env

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
	_ "golang.org/x/sync/singleflight"
	"time"
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

type CacheOptions struct {
	Key      string
	Duration time.Duration
	Fun      func() (interface{}, error) //用于获取结果数据的函数，错误向上传递
	Receiver interface{}                 //初始化的 结构体指针或者map，用来存放结果
}

//获取缓存，不存在则调用fun获取数据加入缓存，fun的错误向上传递
func GetSetCache(c *CacheOptions) (using bool, err error) {

	if c == nil || c.Receiver == nil || c.Key == "" {
		err = fmt.Errorf("illegal arguments")
		logs.Error(err)
		return
	}

	//====查询缓存
	val, err := RedisCli.Get(c.Key).Result()
	if err != nil && err != redis.Nil {
		logs.Error(err)
		return
	}

	if err == redis.Nil {
		//调用函数获取数据
		c.Receiver, err = c.Fun()
		if err != nil {
			return
		}
		if fmt.Sprint(c.Receiver) == "<nil>" {
			return false, nil
		}

		//加入缓存
		logs.Debug("Set cache %s", c.Key)
		var buf []byte
		buf, err = json.Marshal(&c.Receiver)
		if err != nil {
			logs.Error(err)
			return
		}
		err = RedisCli.Set(c.Key, string(buf), c.Duration).Err()
		if err != nil {
			logs.Error(err)
			return
		}

	} else {
		//解析缓存
		using = true
		logs.Debug("Hit cache %s", c.Key)
		err = json.Unmarshal([]byte(val), &c.Receiver)
		if err != nil {
			logs.Error(err)
			return
		}
	}
	return
}
