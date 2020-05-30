package routers

import (
	"GZHU-Pi/env"
	"GZHU-Pi/pkg"
	"GZHU-Pi/pkg/gzhu_jw"
	"GZHU-Pi/services/kafka"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
	"gopkg.in/guregu/null.v3"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TODO Note
//var Jwxt = make(map[string]pkg.Jwxt)
var Jwxt sync.Map

func newJWClient(r *http.Request, username, password string) (client pkg.Jwxt, err error) {

	school := r.URL.Query().Get("school")
	if school == "" {
		school = "gzhu"
	}

	if school == "gzhu" {
		return gzhu_jw.BasicAuthClient(username, password)
	}

	if school == "demo" {
		client = &pkg.Demo{Username: username, Password: password}
	}
	return
}

func getCacheKey(r *http.Request, username string) string {
	s := r.URL.Query().Get("school")
	if r.URL.Query().Get("school") == "" {
		s = "gzhu"
	}
	return s + username
}

//教务系统统一中间件，做一些准备客户端的公共操作
func JWMiddleWare(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if strings.ToUpper(r.Method) == "GET" {
			username := r.URL.Query().Get("username")
			if username == "" {
				Response(w, r, nil, http.StatusUnauthorized, "Unauthorized")
				return
			}
			//client, ok := Jwxt[getCacheKey(r, username)]
			c, _ := Jwxt.Load(getCacheKey(r, username))
			client, ok := c.(pkg.Jwxt)
			if !ok {
				Response(w, r, nil, http.StatusUnauthorized, "Unauthorized")
				return
			}
			ctx := context.WithValue(r.Context(), "client", client)
			// 创建新的请求
			r = r.WithContext(ctx)
			next(w, r)
			return
		}

		u, err := ReadRequestArg(r, "username")
		p, err0 := ReadRequestArg(r, "password")
		if err != nil || err0 != nil {
			logs.Error(err, err0)
			Response(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		username, _ := u.(string)
		password, _ := p.(string)
		if username == "" || password == "" {
			Response(w, r, nil, http.StatusUnauthorized, "Unauthorized")
			return
		}
		logs.Info("用户：%s 接口：%s IP: %s", username, r.URL.Path, r.RemoteAddr)

		//从缓存中获取客户端，不存在或者过期则创建
		//client, ok := Jwxt[getCacheKey(r, username)]
		c, _ := Jwxt.Load(getCacheKey(r, username))
		client, ok := c.(pkg.Jwxt)
		if !ok || client == nil || time.Now().After(client.GetExpiresAt()) {
			if client != nil && time.Now().After(client.GetExpiresAt()) {
				logs.Debug("客户端在 %ds 前过期了", time.Now().Unix()-client.GetExpiresAt().Unix())
			}
			client, err = newJWClient(r, username, password)
			if err != nil {
				logs.Error(err)
				Response(w, r, nil, http.StatusUnauthorized, err.Error())
				return
			}
			//将客户端存入缓存
			//Jwxt[getCacheKey(r, username)] = client
			Jwxt.Store(getCacheKey(r, username), client)
		}
		if client != nil && !time.Now().After(client.GetExpiresAt()) {
			logs.Debug("客户端正常 %ds 后过期", client.GetExpiresAt().Unix()-time.Now().Unix())
		}
		//如果客户端不发生错误而被删除，则更新过期时间
		defer func() {
			c, _ := Jwxt.Load(getCacheKey(r, username))
			client, ok := c.(pkg.Jwxt)
			if ok || client != nil {
				client.SetExpiresAt(time.Now().Add(20 * time.Minute))
			}
		}()
		//把客户端通过context传递给下一级
		ctx := context.WithValue(r.Context(), "client", client)
		// 创建新的请求
		r = r.WithContext(ctx)
		next(w, r)

	}
}

//获取课表信息，参数为 year semester
func Course(w http.ResponseWriter, r *http.Request) {
	//从context提取客户端
	c := r.Context().Value("client")
	if c == nil {
		Response(w, r, nil, http.StatusInternalServerError, "get nil client from context")
		return
	}
	client, ok := c.(pkg.Jwxt)
	if !ok {
		Response(w, r, nil, http.StatusInternalServerError, "get a wrong client from context")
		return
	}

	year, sem := gzhu_jw.Year, gzhu_jw.SemCode[1]
	s, _ := ReadRequestArg(r, "year_sem")
	ys, _ := s.(string)
	yearSem := strings.Split(ys, "-")
	if len(yearSem) == 3 {
		year = yearSem[0]
		sem = yearSem[2]
		if sem == "1" {
			sem = "3"
		}
		if sem == "2" {
			sem = "12"
		}
	}

	data, err := client.GetCourse(year, sem)
	if err != nil {
		logs.Error(err)
		Jwxt.Delete(getCacheKey(r, client.GetUsername())) //发生错误，从缓存中删除
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}

	for _, v := range data.CourseList {
		v.YearSem = ys
		v.StuID = client.GetUsername()
	}
	Response(w, r, data, http.StatusOK, "request ok")

	//====响应后的处理
	var userID int64
	if len(r.Cookies()) > 0 {
		userID, err = GetUserID(r)
		if err != nil {
			logs.Error(err, r.Cookies())
			return
		}
	}
	for _, v := range data.CourseList {
		v.CreatedBy = null.IntFrom(userID)
	}

	if len(data.CourseList) == 0 {
		return
	}

	//添加自动提醒，每个学生每学期，只自动设置上课提醒一次
	go func() {
		if !env.Conf.Kafka.Enabled {
			return
		}
		fm, _ := ReadRequestArg(r, "first_monday")
		firstMonday, _ := fm.(string)
		if firstMonday == "" {
			firstMonday = gzhu_jw.FirstMonday
			logs.Warn("use default firstMonday %s", firstMonday)
		}
		key := fmt.Sprintf("gzhupi:notify:course:stu:%d_%s_%s", userID, ys, client.GetUsername())
		_, err := env.RedisCli.Get(key).Result()
		if err == redis.Nil {
			err = AddCourseNotify(data.CourseList, firstMonday)
			if err != nil {
				return
			}
			err = env.RedisCli.Set(key, key, 120*24*time.Hour).Err()
			if err != nil {
				logs.Error(err)
				return
			}
		} else if err != nil {
			logs.Error(err, key)
			return
		}
		logs.Debug("key exists, skip add course notify key=%s", key)
	}()

}

func Exam(w http.ResponseWriter, r *http.Request) {
	//从context提取客户端
	c := r.Context().Value("client")
	if c == nil {
		Response(w, r, nil, http.StatusInternalServerError, "get nil client from context")
		return
	}
	client, ok := c.(pkg.Jwxt)
	if !ok {
		Response(w, r, nil, http.StatusInternalServerError, "get a wrong client from context")
		return
	}

	year, sem := gzhu_jw.Year, gzhu_jw.SemCode[1]
	s, _ := ReadRequestArg(r, "year_sem")
	ys, _ := s.(string)
	yearSem := strings.Split(ys, "-")
	if len(yearSem) == 3 {
		year = yearSem[0]
		sem = yearSem[2]
		if sem == "1" {
			sem = "3"
		}
		if sem == "2" {
			sem = "12"
		}
	}

	data, err := client.GetExam(year, sem)
	if err != nil {
		logs.Error(err)
		Jwxt.Delete(getCacheKey(r, client.GetUsername()))
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	Response(w, r, data, http.StatusOK, "request ok")
}

func Grade(w http.ResponseWriter, r *http.Request) {
	//从context提取客户端
	c := r.Context().Value("client")
	if c == nil {
		Response(w, r, nil, http.StatusInternalServerError, "get nil client from context")
		return
	}
	client, ok := c.(pkg.Jwxt)
	if !ok {
		Response(w, r, nil, http.StatusInternalServerError, "get a wrong client from context")
		return
	}

	data, err := client.GetAllGrade("", "")
	if err != nil {
		logs.Error(err)
		Jwxt.Delete(getCacheKey(r, client.GetUsername()))
		if err == gzhu_jw.AuthError {
			Response(w, r, nil, http.StatusUnauthorized, err.Error())
			return
		}
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}

	Response(w, r, data, http.StatusOK, "request ok")

	stuInfo := data.StuInfo
	var grades []*env.TGrade
	for _, v := range data.SemList {
		grades = append(grades, v.GradeList...)
	}
	//加入消息队列
	go func() {
		if !env.Conf.Kafka.Enabled {
			return
		}
		info, err := json.Marshal(stuInfo)
		if err != nil {
			logs.Error(err)
			return
		}
		err = env.Kafka.SendData(&kafka.ProduceData{
			Topic: env.QueueTopicStuInfo,
			Data:  info,
		})
		if err != nil {
			logs.Error(err)
			return
		}

		//====
		g, err := json.Marshal(grades)
		if err != nil {
			logs.Error(err)
			return
		}
		err = env.Kafka.SendData(&kafka.ProduceData{
			Topic: env.QueueTopicGrade,
			Data:  g,
		})
		if err != nil {
			logs.Error(err)
			return
		}
	}()

}

func EmptyRoom(w http.ResponseWriter, r *http.Request) {
	//从context提取客户端
	c := r.Context().Value("client")
	if c == nil {
		Response(w, r, nil, http.StatusInternalServerError, "get nil client from context")
		return
	}
	client, ok := c.(pkg.Jwxt)
	if !ok {
		Response(w, r, nil, http.StatusInternalServerError, "get a wrong client from context")
		return
	}

	data, err := client.GetEmptyRoom(r)
	if err != nil {
		logs.Error(err)
		Jwxt.Delete(getCacheKey(r, client.GetUsername()))
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	Response(w, r, data, http.StatusOK, "request ok")
}

func Achieve(w http.ResponseWriter, r *http.Request) {
	//从context提取客户端
	c := r.Context().Value("client")
	if c == nil {
		Response(w, r, nil, http.StatusInternalServerError, "get nil client from context")
		return
	}
	client, ok := c.(pkg.Jwxt)
	if !ok {
		Response(w, r, nil, http.StatusInternalServerError, "get a wrong client from context")
		return
	}

	data, err := client.GetAchieve()
	if err != nil {
		logs.Error(err)
		Jwxt.Delete(getCacheKey(r, client.GetUsername()))
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	Response(w, r, data, http.StatusOK, "request ok")
}

func AllCourse(w http.ResponseWriter, r *http.Request) {
	//从context提取客户端
	c := r.Context().Value("client")
	if c == nil {
		Response(w, r, nil, http.StatusInternalServerError, "get nil client from context")
		return
	}
	client, ok := c.(pkg.Jwxt)
	if !ok {
		Response(w, r, nil, http.StatusInternalServerError, "get a wrong client from context")
		return
	}

	year := r.URL.Query().Get("year")
	sem := r.URL.Query().Get("sem")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	count, _ := strconv.Atoi(r.URL.Query().Get("count"))

	data, csvData, err := client.SearchAllCourse(year, sem, page, count)
	if err != nil {
		logs.Error(err)
		Jwxt.Delete(getCacheKey(r, client.GetUsername()))
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}

	//导出文件
	if r.URL.Query().Get("action") == "export" {
		w.Header().Set("Content-Type", "application/csv")
		w.Header().Set("Content-Disposition",
			fmt.Sprintf("attachment; filename=%s", "export.csv"))
		_, _ = w.Write(csvData)
	} else {
		Response(w, r, data, http.StatusOK, "request ok")
	}
}

func Rank(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("username")
	user, err := VUserByCookies(r)
	if err != nil {
		Response(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}
	if user.StuID.String != "" && user.StuID.String != username {
		err = fmt.Errorf("Unauthorized ")
		logs.Error(err)
		Response(w, r, nil, http.StatusUnauthorized, err.Error())
		return
	}
	logs.Info("用户：%d 学号：%s 接口：%s IP: %s", user.ID, username, r.URL.Path, r.RemoteAddr)

	var data map[string]interface{}

	//====查询缓存
	key := fmt.Sprintf("gzhupi:rank:%s", username)
	val, err := env.RedisCli.Get(key).Result()
	if err != nil && err != redis.Nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	if err == redis.Nil {

		client := &gzhu_jw.JWClient{Username: username}
		data, err = client.GetRank(client.GetUsername())
		if err != nil {
			logs.Error(err)
			Response(w, r, nil, http.StatusInternalServerError, err.Error())
			return
		}
		//加入缓存
		logs.Debug("Set cache %s", key)
		buf, err := json.Marshal(&data)
		if err != nil {
			logs.Error(err)
			Response(w, r, nil, http.StatusInternalServerError, err.Error())
			return
		}
		//缓存一天
		err = env.RedisCli.Set(key, string(buf), 1*24*time.Hour).Err()
		if err != nil {
			logs.Error(err)
			Response(w, r, nil, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		//解析缓存
		logs.Debug("Hit cache %s", key)
		err = json.Unmarshal([]byte(val), &data)
		if err != nil {
			logs.Error(err)
			Response(w, r, nil, http.StatusInternalServerError, err.Error())
			return
		}
	}

	Response(w, r, data, http.StatusOK, "request ok")
}
