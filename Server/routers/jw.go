package routers

import (
	"GZHU-Pi/pkg/gzhu_jw"
	"context"
	"github.com/astaxie/beego/logs"
	"net/http"
	"time"
)

var JWClients = make(map[string]*gzhu_jw.JWClient)

//教务系统统一中间件，做一些准备客户端的公共操作
func JWMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := ReadRequestArg(r, "username")
		p, err0 := ReadRequestArg(r, "password")
		if err != nil || err0 != nil {
			logs.Error(err, err0)
			Response(w, r, nil, http.StatusBadRequest, "illegal request form")
			return
		}
		username, _ := u.(string)
		password, _ := p.(string)
		if username == "" || password == "" {
			Response(w, r, nil, http.StatusUnauthorized, "Unauthorized")
			return
		}
		//从缓存中获取客户端，不存在或者过期则创建
		client, ok := JWClients[username]
		if !ok || client == nil || time.Now().After(client.ExpiresAt) {
			if client != nil && time.Now().After(client.ExpiresAt) {
				logs.Debug("客户端在 %ds 前过期了", time.Now().Unix()-client.ExpiresAt.Unix())
			}
			client, err = gzhu_jw.BasicAuthClient(username, password)
			if err != nil {
				logs.Error(err)
				Response(w, r, nil, http.StatusUnauthorized, err.Error())
				return
			}
			//将客户端存入缓存
			JWClients[username] = client
		}
		if client != nil && !time.Now().After(client.ExpiresAt) {
			logs.Debug("客户端正常 %ds 后过期", client.ExpiresAt.Unix()-time.Now().Unix())
		}
		//如果客户端不发生错误而被删除，则更新过期时间
		defer func() {
			if client, ok := JWClients[username]; ok || client != nil {
				client.ExpiresAt = time.Now().Add(20 * time.Minute)
			}
		}()
		logs.Info("用户：%s 接口：%s",username, r.URL.Path)
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
	client, ok := c.(*gzhu_jw.JWClient)
	if !ok {
		Response(w, r, nil, http.StatusInternalServerError, "get a wrong client from context")
		return
	}

	data, err := client.GetCourse(gzhu_jw.Year, gzhu_jw.SemCode[0])
	if err != nil {
		logs.Error(err)
		delete(JWClients, client.Username) //发生错误，从缓存中删除
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	Response(w, r, data, http.StatusOK, "request ok")
}

func Exam(w http.ResponseWriter, r *http.Request) {
	//从context提取客户端
	c := r.Context().Value("client")
	if c == nil {
		Response(w, r, nil, http.StatusInternalServerError, "get nil client from context")
		return
	}
	client, ok := c.(*gzhu_jw.JWClient)
	if !ok {
		Response(w, r, nil, http.StatusInternalServerError, "get a wrong client from context")
		return
	}

	data, err := client.GetExam(gzhu_jw.Year, gzhu_jw.SemCode[0])
	if err != nil {
		logs.Error(err)
		delete(JWClients, client.Username)
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
	client, ok := c.(*gzhu_jw.JWClient)
	if !ok {
		Response(w, r, nil, http.StatusInternalServerError, "get a wrong client from context")
		return
	}

	data, err := client.GetAllGrade("", "")
	if err != nil {
		logs.Error(err)
		delete(JWClients, client.Username)
		if err == gzhu_jw.AuthError {
			//TODO 重试处理

			Response(w, r, nil, http.StatusUnauthorized, err.Error())
			return
		}
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	Response(w, r, data, http.StatusOK, "request ok")
}


func EmptyRoom(w http.ResponseWriter, r *http.Request) {
	//从context提取客户端
	c := r.Context().Value("client")
	if c == nil {
		Response(w, r, nil, http.StatusInternalServerError, "get nil client from context")
		return
	}
	client, ok := c.(*gzhu_jw.JWClient)
	if !ok {
		Response(w, r, nil, http.StatusInternalServerError, "get a wrong client from context")
		return
	}

	data, err := client.GetEmptyRoom(r)
	if err != nil {
		logs.Error(err)
		delete(JWClients, client.Username)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	Response(w, r, data, http.StatusOK, "request ok")
}