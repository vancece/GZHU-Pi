package env

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

func InitConfig() {

	viper.SetConfigName("config") //指定配置文件的文件名称(不需要制定配置文件的扩展名)
	viper.AddConfigPath(".")      // 设置配置文件和可执行二进制文件在用一个目录
	viper.AutomaticEnv()          //自动从环境变量读取匹配的参数

	//绑定放进变量，会优先读取环境变量的值
	_ = viper.BindEnv("db.host", "GZHUPI_DB_HOST")
	_ = viper.BindEnv("db.port", "GZHUPI_DB_PORT")
	_ = viper.BindEnv("db.user", "GZHUPI_DB_USER")
	_ = viper.BindEnv("db.password", "GZHUPI_DB_PASSWORD")
	_ = viper.BindEnv("db.dbname", "GZHUPI_DB_PASSWORD")

	_ = viper.BindEnv("redis.host", "GZHUPI_REDIS_USER")
	_ = viper.BindEnv("redis.port", "GZHUPI_REDIS_PORT")
	_ = viper.BindEnv("redis.password", "GZHUPI_REDIS_PASSWORD")

	_ = viper.BindEnv("weixin.path", "GZHUPI_WEIXIN_PATH")
	_ = viper.BindEnv("weixin.app_id", "GZHUPI_WEIXIN_APP_ID")
	_ = viper.BindEnv("weixin.secret", "GZHUPI_WEIXIN_PATH")
	_ = viper.BindEnv("weixin.token", "GZHUPI_WEIXIN_TOKEN")

	//读取-c输入的路径参数，初始化配置文件，如： ./main -c config.yaml
	if len(os.Args) >= 3 {
		if os.Args[1] == "-c" {
			cfgFile := os.Args[2]
			viper.SetConfigFile(cfgFile)
		}
	}
	// 根据以上配置读取加载配置文件
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	fmt.Println("获取配置文件的string", viper.GetString(`app.name`))
	fmt.Println("获取配置文件的string", viper.GetInt(`app.foo`))
	fmt.Println("获取配置文件的string", viper.GetBool(`app.bar`))
	fmt.Println("获取配置文件的map[string]string", viper.GetStringMapString(`app`))

}
