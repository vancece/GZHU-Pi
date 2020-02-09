package routers

import (
	"GZHU-Pi/models"
	"GZHU-Pi/pkg/gzhu_jw"
	"fmt"
	"github.com/astaxie/beego/logs"
	"net/http"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	o, err := ReadRequestArg(r, "open_id")
	u, err1 := ReadRequestArg(r, "username")
	p, err2 := ReadRequestArg(r, "password")
	if err != nil || err1 != nil || err2 != nil {
		logs.Error(err, err1, err2)
		Response(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}
	openID, _ := o.(string)
	username, _ := u.(string)
	password, _ := p.(string)

	if openID == "" && (username == "" || password == "") {
		Response(w, r, nil, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var user models.TUser
	if openID != "" {
		logs.Info("auth by open_id:",openID)

		db := models.GetGorm()
		db.Where("open_id = ?", openID).First(&user)
		logs.Info(user)

		if user.ID <= 0 {
			err = fmt.Errorf("user not found")
			logs.Error(err)
			Response(w, r, nil, http.StatusUnauthorized, err.Error())
			return
		}
		cookieStr, err := NewCookie(user.ID)
		if err != nil {
			logs.Error(err)
			Response(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		w.Header().Set("Set-Cookie", cookieStr)
		Response(w, r, user, http.StatusOK, "")
		return
	}

	//使用学号认证 不设置cookie
	if username != "" && password != "" {
		logs.Info("auth by stu_id")
		client, err := gzhu_jw.BasicAuthClient(username, password)
		if err != nil {
			logs.Error(err)
			Response(w, r, nil, http.StatusUnauthorized, err.Error())
			return
		}
		//将客户端存入缓存
		JWClients[username] = client
		//TODO get stu_info
		Response(w, r, nil, http.StatusOK, "")
		return
	}

}
