package env

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
	_ "golang.org/x/sync/singleflight"
)

var RedisCli *redis.Client

func InitRedis() (err error) {
	return
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

func ExampleClient() {
	err := RedisCli.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := RedisCli.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := RedisCli.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
