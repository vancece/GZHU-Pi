package routers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

//后端响应数据通信协议
type ResponseProto struct {
	Status      int         `json:"status"`       //状态 0正常，小于0出错，大于0可能有问题
	Msg         string      `json:"msg"`          //状态信息
	Data        interface{} `json:"data"`         //响应数据
	Api         string      `json:"api"`          //api接口
	Method      string      `json:"method"`       //post,put,get,delete
	Count       int         `json:"count"`        //Data若是数组，算其长度
	Time        int64       `json:"time"`         //请求响应时间，毫秒
	UpdatedTime string      `json:"updated_time"` //响应处理时间
}

//前端请求数据通讯协议
//type RequestProto struct {
//	Action   string      `json:"action"` //请求类型GET/POST/PUT/DELETE
//	Data     interface{} `json:"data"`   //请求数据
//	Sets     []string    `json:"sets"`
//	OrderBy  string      `json:"orderBy"`  //排序要求
//	Filter   string      `json:"filter"`   //筛选条件
//	Page     int         `json:"page"`     //分页
//	PageSize int         `json:"pageSize"` //分页大小
//}

//统一响应处理函数
func Response(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int, msg string) {
	if w == nil || r == nil {
		http.Error(w, "Unknown Error", http.StatusInternalServerError)
		return
	}
	//计算响应时长
	startTime := r.Context().Value("startTime")
	last := time.Duration(0)
	if startTime != nil {
		_, ok := startTime.(time.Time)
		if ok {
			last = time.Since(startTime.(time.Time))
		}
	}
	resp := ResponseProto{}
	resp.Api = r.URL.Path
	resp.Status = statusCode
	resp.Msg = msg
	resp.Data = data
	resp.Method = r.Method
	resp.Time = last.Nanoseconds() / 1000000
	resp.UpdatedTime = time.Now().Format("2006-01-02 15:04:05")

	//统计数组/切片长度
	//if data != nil {
	//	k := reflect.TypeOf(data)
	//	//if k.Kind() == reflect.Slice || k.Kind() == reflect.Array {
	//	if  k.Kind() == reflect.Array {
	//		resp.Count = k.Len()
	//	}
	//}
	response, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ErrorHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				err := fmt.Errorf("recover a panic: %+v", err)
				logs.Error(err)
				Response(w, r, nil, http.StatusInternalServerError, err.Error())
			}
		}()
		w.Header().Set("Content-Type", "application/json")
		//请求开始时间
		startTime := time.Now()
		ctx := context.WithValue(r.Context(), "startTime", startTime)
		// 创建新的请求
		r = r.WithContext(ctx)
		h(w, r)
	}
}

//识别基础请求类型读取参数
func ReadRequestArg(r *http.Request, key string) (value interface{}, err error) {
	if key == "" {
		return nil, fmt.Errorf("invalid key")
	}
	//application/json post请求解析body
	if r.Method == "POST" && reflect.DeepEqual(r.Header["Content-Type"], []string{"application/json"}) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logs.Error(err)
			return "", err
		}
		_ = r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		if len(body) == 0 {
			return "", fmt.Errorf("non body")
		}
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			logs.Error(err)
			return "", err
		}
		return data[key], nil
	}
	//get请求参数或者application/x-www-form-urlencoded
	if r.Method == "GET" || reflect.DeepEqual(r.Header["Content-Type"], []string{"application/x-www-form-urlencoded"}) {
		_ = r.ParseForm()
		value = r.Form.Get(key)
		return value, nil
	}
	return nil, fmt.Errorf("unsupported method: %s or content type: %v", r.Method, r.Header["Content-Type"])
}
