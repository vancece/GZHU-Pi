package main

import (
	"GZHU-Pi/env"
	"GZHU-Pi/routers"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type App struct {
	Client *http.Client
	Config []string
	Router []http.Handler
}

func (app *App) InitApp() error {

	return nil
}

func main() {

	_ = env.InitLogger()

	r := mux.NewRouter()
	r = r.PathPrefix("/api/v1").Subrouter()

	//微信公众号接口
	//r.HandleFunc("/wx/check", routers.PanicMV(routers.WeChatCheck))
	//r.HandleFunc("/wx/check", routers.PanicMV(routers.Hello))

	//教务系统
	r.HandleFunc("/jwxt/course", routers.PanicMV(routers.JWMiddleWare(routers.Course))).Methods("POST")
	r.HandleFunc("/jwxt/exam", routers.PanicMV(routers.JWMiddleWare(routers.Exam))).Methods("POST")
	r.HandleFunc("/jwxt/grade", routers.PanicMV(routers.JWMiddleWare(routers.Grade))).Methods("POST")
	r.HandleFunc("/jwxt/room", routers.PanicMV(routers.JWMiddleWare(routers.EmptyRoom))).Methods("POST")
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

	http.Handle("/", r)
	//获取阿里云函数计算容器端口
	port := os.Getenv("FC_SERVER_PORT")
	if port == "" {
		port = "9000"
	}
	fmt.Println("Server on port: " + port)
	err1 := http.ListenAndServe(":"+port, r)
	if err1 != nil {
		log.Print(err1)
		return
	}
}
