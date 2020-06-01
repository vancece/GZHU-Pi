package routers

import (
	"GZHU-Pi/pkg/gzhu_second"
	"context"
	"github.com/astaxie/beego/logs"
	"github.com/mo7zayed/reqip"
	"net/http"
	"strings"
	"time"
)

var SecondClients = make(map[string]*gzhu_second.SecondClient)

//教务系统统一中间件，做一些准备客户端的公共操作
func SecondMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if strings.ToUpper(r.Method) == "GET" {
			username := r.URL.Query().Get("username")
			if username == "" {
				Response(w, r, nil, http.StatusUnauthorized, "Unauthorized")
				return
			}
			client, ok := SecondClients[username]
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
		client, ok := SecondClients[username]
		if !ok || client == nil || time.Now().After(client.ExpiresAt) {
			if client != nil && time.Now().After(client.ExpiresAt) {
				logs.Debug("客户端在 %ds 前过期了", time.Now().Unix()-client.ExpiresAt.Unix())
			}
			client, err = gzhu_second.BasicAuthClient(username, password)
			if err != nil {
				logs.Error(err)
				Response(w, r, nil, http.StatusUnauthorized, err.Error())
				return
			}
			//将客户端存入缓存
			SecondClients[username] = client
		}
		if client != nil && !time.Now().After(client.ExpiresAt) {
			logs.Debug("客户端正常 %ds 后过期", client.ExpiresAt.Unix()-time.Now().Unix())
		}
		//如果客户端不发生错误而被删除，则更新过期时间
		defer func() {
			if client, ok := SecondClients[username]; ok || client != nil {
				client.ExpiresAt = time.Now().Add(20 * time.Minute)
			}
		}()
		logs.Info("用户：%s IP: %s 接口：%s ", username, reqip.GetClientIP(r), r.URL.Path)
		//把客户端通过context传递给下一级
		ctx := context.WithValue(r.Context(), "client", client)
		// 创建新的请求
		r = r.WithContext(ctx)

		next(w, r)
	}
}

func MySecond(w http.ResponseWriter, r *http.Request) {
	//从context提取客户端
	c := r.Context().Value("client")
	if c == nil {
		Response(w, r, nil, http.StatusInternalServerError, "get nil client from context")
		return
	}
	client, ok := c.(*gzhu_second.SecondClient)
	if !ok {
		Response(w, r, nil, http.StatusInternalServerError, "get a wrong client from context")
		return
	}

	data, err := client.GetMySecond()
	if err != nil {
		logs.Error(err)
		delete(SecondClients, client.Username)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	Response(w, r, data, http.StatusOK, "request ok")
}

func SecondSearch(w http.ResponseWriter, r *http.Request) {
	//从context提取客户端
	c := r.Context().Value("client")
	if c == nil {
		Response(w, r, nil, http.StatusInternalServerError, "get nil client from context")
		return
	}
	client, ok := c.(*gzhu_second.SecondClient)
	if !ok {
		Response(w, r, nil, http.StatusInternalServerError, "get a wrong client from context")
		return
	}
	data, err := client.Search(r)
	if err != nil {
		logs.Error(err)
		delete(SecondClients, client.Username)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	Response(w, r, data, http.StatusOK, "request ok")
}

//根据申报项目id获取图片
func SecondImage(w http.ResponseWriter, r *http.Request) {
	//从context提取客户端
	c := r.Context().Value("client")
	if c == nil {
		Response(w, r, nil, http.StatusInternalServerError, "get nil client from context")
		return
	}
	client, ok := c.(*gzhu_second.SecondClient)
	if !ok {
		Response(w, r, nil, http.StatusInternalServerError, "get a wrong client from context")
		return
	}
	id, err := ReadRequestArg(r, "id")
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusBadRequest, "illegal request form")
		return
	}
	itemID, _ := id.(string)
	data, err := client.GetImages(itemID)
	if err != nil {
		logs.Error(err)
		delete(SecondClients, client.Username)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	Response(w, r, data, http.StatusOK, "request ok")
}
