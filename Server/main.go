package main

import (
	"GZHU-Pi/env"
	"GZHU-Pi/models"
	"GZHU-Pi/routers"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/mux"
	"github.com/prest/adapters/postgres"
	"github.com/prest/cmd"
	"github.com/prest/config"
	"github.com/prest/config/router"
	"github.com/prest/middlewares"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {

	err := env.InitConfigure()
	if err != nil {
		log.Fatal(err)
	}

	// init db without handling errors
	go models.InitDb()

	r := mux.NewRouter()

	//获取阿里云函数计算容器内部端口
	port := os.Getenv("FC_SERVER_PORT")
	if port == "" {
		port = "9000"
		logs.Info("自主部署 Server on port: " + port)
		_ = env.InitLogger("./log/")
		r = r.PathPrefix("/api/v1").Subrouter()
	} else {
		logs.Info("阿里云云函数部署 Server on port: " + port)
		_ = env.InitLogger("/tmp/log/")
		r = r.PathPrefix("/2016-08-15/proxy/GZHU-API/go/api/v1").Subrouter()
	}

	if env.Conf.App.PRest == true {

		logs.Info("启用pRest接口服务")
		runWithPRest(r)

	} else {

		customRouter(r)
		if err := http.ListenAndServe(":"+port, r); err != nil {
			log.Fatal(err)
		}
	}
}

func runWithPRest(r *mux.Router) {

	//同步应用配置 覆盖pRest部分配置
	viper.Set("pg.host", viper.GetString("db.host"))
	viper.Set("pg.user", viper.GetString("db.user"))
	viper.Set("pg.pass", viper.GetString("db.password"))
	viper.Set("pg.port", viper.GetInt("db.port"))
	viper.Set("pg.database", viper.GetString("db.dbname"))
	viper.Set("ssl.mode", viper.GetString("db.sslmode"))

	// load config for pREST
	config.Load()

	// Load Postgres Adapter
	postgres.Load()

	// Get pREST app
	n := middlewares.GetApp()

	// Register custom middleware
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		//该中间件用于消除路由前缀对pRest内部路由的影响
		logs.Info(r.URL.Path)
		if !strings.Contains(r.URL.Path, "/api/v") {
			_, _ = w.Write([]byte("path must contains /api/v\\d"))
			return
		}
		reg := regexp.MustCompile(`[\s\S]*/api/v\d`)
		match := reg.FindStringSubmatch(r.URL.Path)
		if len(match) > 0 {
			r.URL.Path = strings.ReplaceAll(r.URL.Path, match[0], "")
		}
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		next(w, r)
	})

	n.UseFunc(routers.TableAccessHandle)

	// Get pREST router
	r = router.Get()

	// just for test
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Pong!"))
	})

	// Register custom routes
	customRouter(r)

	// Call pREST cmd
	cmd.Execute()
}

func customRouter(r *mux.Router) *mux.Router {

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello!"))
	})

	r.HandleFunc("/auth", routers.PanicMV(routers.Auth)).Methods("POST")

	//微信公众号接口
	//r.HandleFunc("/wx/check", routers.PanicMV(routers.WeChatCheck))
	//r.HandleFunc("/wx/check", routers.PanicMV(routers.Hello))

	//教务系统
	r.HandleFunc("/jwxt/course", routers.PanicMV(routers.JWMiddleWare(routers.Course))).Methods("POST")
	r.HandleFunc("/jwxt/exam", routers.PanicMV(routers.JWMiddleWare(routers.Exam))).Methods("POST")
	r.HandleFunc("/jwxt/grade", routers.PanicMV(routers.JWMiddleWare(routers.Grade))).Methods("POST")
	r.HandleFunc("/jwxt/classroom", routers.PanicMV(routers.JWMiddleWare(routers.EmptyRoom))).Methods("POST")
	r.HandleFunc("/jwxt/achieve", routers.PanicMV(routers.JWMiddleWare(routers.Achieve))).Methods("POST")

	//图书馆
	r.HandleFunc("/library/search", routers.PanicMV(routers.BookSearch)).Methods("GET")
	r.HandleFunc("/library/holdings", routers.PanicMV(routers.BookHoldings)).Methods("GET")

	//第二课堂学分系统
	r.HandleFunc("/second/my", routers.PanicMV(routers.SecondMiddleWare(routers.MySecond))).Methods("POST")
	r.HandleFunc("/second/search", routers.PanicMV(routers.SecondMiddleWare(routers.SecondSearch))).Methods("POST")
	r.HandleFunc("/second/image", routers.PanicMV(routers.SecondMiddleWare(routers.SecondImage))).Methods("POST")

	//物理实验平台
	//r.HandleFunc("/exp", test).Methods("POST")

	//四六级、普通话考试查询
	//r.HandleFunc("/exam/cet", test).Methods("POST")
	//r.HandleFunc("/exam/chinese", test).Methods("POST")
	return r
}
