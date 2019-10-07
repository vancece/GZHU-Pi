package routers

import (
	"GZHU-Pi/pkg/gzhu_jw"
	"github.com/astaxie/beego/logs"
	"net/http"
)

var JWClients = make(map[string]*gzhu_jw.JWClient)

func Course(w http.ResponseWriter, r *http.Request) {
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
	client, ok := JWClients[username]
	if !ok {
		client, err = gzhu_jw.BasicAuthClient(username, password)
		if err != nil {
			logs.Error(err)
			Response(w, r, nil, http.StatusUnauthorized, "Unauthorized")
			return
		}
	}
	//更新客户端
	defer func() {
		JWClients[username] = client
	}()

	data, err := client.GetCourse(gzhu_jw.Year, gzhu_jw.SemCode[0])
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, "Unauthorized")
		return
	}
	Response(w, r, data, http.StatusOK, "request ok")
}

func Exam(w http.ResponseWriter, r *http.Request) {
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
	client, ok := JWClients[username]
	if !ok {
		client, err = gzhu_jw.BasicAuthClient(username, password)
		if err != nil {
			logs.Error(err)
			Response(w, r, nil, http.StatusUnauthorized, "Unauthorized")
			return
		}
	}
	//更新客户端
	defer func() {
		JWClients[username] = client
	}()

	data, err := client.GetExam(gzhu_jw.Year, gzhu_jw.SemCode[0])
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, "Unauthorized")
		return
	}
	Response(w, r, data, http.StatusOK, "request ok")
}

func Grade(w http.ResponseWriter, r *http.Request) {
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
	client, ok := JWClients[username]
	if !ok {
		client, err = gzhu_jw.BasicAuthClient(username, password)
		if err != nil {
			logs.Error(err)
			Response(w, r, nil, http.StatusUnauthorized, "Unauthorized")
			return
		}
	}
	//更新客户端
	defer func() {
		JWClients[username] = client
	}()

	data, err := client.GetAllGrade("", "")
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, "Unauthorized")
		return
	}
	Response(w, r, data, http.StatusOK, "request ok")
}
