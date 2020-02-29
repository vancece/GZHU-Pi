package env

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/viper"
	"log"
	"os"
	"reflect"
	"strings"
)

func InitViper() {

	viper.SetConfigName("config") //指定配置文件的文件名称(不需要制定配置文件的扩展名)
	viper.AddConfigPath(".")      // 设置配置文件和可执行二进制文件在用一个目录
	viper.AutomaticEnv()          //自动从环境变量读取匹配的参数

	//绑定放进变量，会优先读取环境变量的值
	_ = viper.BindEnv("secret.jwt", "GZHUPI_SECRET_JWT")

	_ = viper.BindEnv("db.host", "GZHUPI_DB_HOST")
	_ = viper.BindEnv("db.port", "GZHUPI_DB_PORT")
	_ = viper.BindEnv("db.user", "GZHUPI_DB_USER")
	_ = viper.BindEnv("db.password", "GZHUPI_DB_PASSWORD")
	_ = viper.BindEnv("db.dbname", "GZHUPI_DB_DBNAME")
	_ = viper.BindEnv("db.sslmode", "GZHUPI_DB_SSLMODE")

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
	InitViper()
	fmt.Println("获取配置文件的string", viper.GetString(`db.user`))
	fmt.Println("获取配置文件的string", viper.GetInt(`db.user`))
	fmt.Println("获取配置文件的string", viper.GetBool(`app.bar`))
	fmt.Println("获取配置文件的string", viper.Sub(`app`).GetString("logfile"))
	fmt.Println("获取配置文件的map[string]string", viper.GetStringMapString(`db`))
	_ = InitConfigure()
	logs.Info(Conf)
}

//配置文件结构体，配置文件上的内容需要一一对应，可多不可少
type Configure struct {
	App struct {
		Name    string `json:"name" remark:"应用名称"`
		Version string `json:"version" remark:"软件发布版本，对应仓库tag版本"`
		Mode    string `json:"mode" remark:"开发模式develop/test/product"`
		PRest   bool   `json:"prest" remark:"" remark:"是否开启pRest接口服务"`
	}
	Secret struct {
		JWT string `json:"jwt" remark:"jwt密钥"`
	}
	Rpc struct {
		Addr string `json:"addr" remark:"rpc主机地址"`
	}
	Db struct {
		Type     string `json:"type" remark:"数据库类型"`
		Host     string `json:"host" remark:"数据库主机"`
		Port     int64  `json:"port" remark:"数据库端口"`
		User     string `json:"user" remark:"数据库用户"`
		Password string `json:"password" remark:"数据库密码"`
		Dbname   string `json:"dbname" remark:"数据库名"`
		SslMode  string `json:"sslmode" remark:"ssl模式"`
	}
	Redis struct {
		Host     string `json:"host" remark:"redis主机"`
		Port     string `json:"port" remark:"redis端口"`
		Password string `json:"password" remark:"redis密码"`
	}
}

var Conf = &Configure{}

//初始化配置信息，测试需要修改配置文件
func InitConfigure() (err error) {
	InitViper()

	confValue := reflect.ValueOf(Conf).Elem()
	confType := reflect.TypeOf(*Conf)

	for i := 0; i < confType.NumField(); i++ {
		section := confType.Field(i)
		sectionValue := confValue.Field(i)

		//读取节类型信息
		for j := 0; j < section.Type.NumField(); j++ {
			key := section.Type.Field(j)
			keyValue := sectionValue.Field(j)

			sec := strings.ToLower(section.Name) //配置文件节名
			remark := key.Tag.Get("remark")      //配置备注
			tag := key.Tag.Get("json")           //配置键节名
			if tag == "" {
				err = fmt.Errorf("can not found a tag name `json` in struct of [%s].%s", sec, tag)
				logs.Error(err)
				return err
			}

			if key.Type.Kind() == reflect.String {

			}
			//根据类型识别配置字段
			switch key.Type.Kind() {
			case reflect.String:
				value := viper.GetString(sec + "." + tag)
				if value == "" {
					err = fmt.Errorf("get a blank value of [%s].%s %s", sec, tag, remark)
					logs.Error(err)
					return err
				}
				keyValue.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				value := viper.GetInt64(sec + "." + tag)
				keyValue.SetInt(value)
			case reflect.Bool:
				value := viper.GetBool(sec + "." + tag)
				keyValue.SetBool(value)

			default:
				logs.Warn("unsupported config struct key type")
			}

			//读取配置文件初始化结构体

		}
	}
	return
}
