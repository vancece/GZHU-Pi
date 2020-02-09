package routers

import (
	"GZHU-Pi/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

//后端响应数据通信协议
type ResponseProto struct {
	Status     int         `json:"status"`      //接口状态码
	Msg        string      `json:"msg"`         //状态信息
	Data       interface{} `json:"data"`        //响应数据
	Api        string      `json:"api"`         //api接口
	Method     string      `json:"method"`      //post,put,get,delete
	Count      int         `json:"count"`       //Data若是数组，算其长度
	Time       int64       `json:"time"`        //请求响应时间，毫秒
	UpdateTime string      `json:"update_time"` //响应处理时间
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
	defer func() {
		logs.Info("============== Responded ==============")
	}()
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
	//if statusCode != 0 && statusCode != 200 {
	//	statusCode = -1
	//} else {
	//	statusCode = 0
	//}
	resp := ResponseProto{}
	resp.Api = r.URL.Path
	resp.Status = statusCode
	resp.Msg = msg
	resp.Data = data
	resp.Method = r.Method
	resp.Time = last.Nanoseconds() / 1000000
	resp.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

	//保存请求记录
	u, _ := ReadRequestArg(r, "username")
	username, _ := u.(string)
	models.SaveApiRecord(&models.TApiRecord{
		Username: username,
		Uri:      r.RequestURI,
		Duration: resp.Time,
	})

	//统计数组/切片长度
	if data != nil {
		k := reflect.TypeOf(data)
		if k.Kind() == reflect.Array {
			resp.Count = k.Len()
		} else if k.Kind() == reflect.Slice {
			resp.Count = int(k.Size())
		}
	}
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

func PanicMV(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				err := fmt.Errorf("recover a panic: %+v", err)
				logs.Error(err)
				Response(w, r, nil, http.StatusInternalServerError, err.Error())
			}
		}()
		w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusOK)
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
		defer r.Body.Close()
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
	return nil, fmt.Errorf("unsupported method: %s with content type: %v", r.Method, r.Header["Content-Type"])
}

//传入用户id，生成并返回Token
func GenerateToken(userID int64) (string, error) {
	//TODO config
	SecretKey := []byte("pi") //设置密钥

	token := jwt.New(jwt.SigningMethodHS256) //指定签名方式，创建token对象
	claims := token.Claims.(jwt.MapClaims)   //Claims (Payload):声明 token 有关的重要信息

	claims["authorized"] = true
	claims["iss"] = userID                                     //指明Token用户
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix() //过期时间

	tokenString, err := token.SignedString(SecretKey)

	if err != nil {
		logs.Error(err)
		return "", err
	}
	return tokenString, nil
}

//传入token字符串，解析Token，返回用户id
func ParseToken(tokenStr string) (userID int64, err error) {
	SecretKey := []byte("pi") //设置密钥

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//methodt, ok := token.Method.(*jwt.SigningMethodHMAC)	//查看加密方式
		return SecretKey, nil
	})
	if err != nil {
		logs.Error(err)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logs.Error("claims not a jwt.MapClaims")
		return
	}
	iss, ok := claims["iss"].(float64)
	if !ok {
		err = fmt.Errorf("iss not float64")
		logs.Error(err)
		return
	}
	if token.Valid {
		return int64(iss), nil
	} else {
		return int64(iss), errors.New("token无效")
	}
}

//从token读取用户id
func GetUserID(r *http.Request) (userID int64, err error) {
	var token string
	if len(r.Cookies()) == 0 {
		err = fmt.Errorf("cookie not set, please authenticate first")
		return
	}
	cookie, err := r.Cookie("token")
	if err != nil {
		logs.Error(err)
		return
	}
	token = cookie.Value
	userID, err = ParseToken(token)
	if err != nil {
		logs.Error(err)
		return
	}
	return
}

func NewCookie(userID int64) (newCookie string, err error) {
	newToken, err := GenerateToken(userID)
	if err != nil {
		logs.Error(err)
		return
	}
	cookie := http.Cookie{
		Name:     "token",
		Value:    newToken,
		HttpOnly: true,
	}
	newCookie = cookie.String()
	return
}
