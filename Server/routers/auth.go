package routers

import (
	"GZHU-Pi/models"
	"GZHU-Pi/pkg/gzhu_jw"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

//使用open_id认证，不存在则创建新用户
func Auth(w http.ResponseWriter, r *http.Request) {

	if r.URL.Query().Get("type") == "gzhu" {
		AuthBySchool(w, r)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	if len(body) == 0 {
		err = fmt.Errorf("Call api by post with empty body ")
		logs.Error(err)
		Response(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}

	var u models.TUser
	err = json.Unmarshal(body, &u)
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}
	if u.OpenID.String == "" || len(u.OpenID.String) != 28 {
		err = fmt.Errorf("must give openid and with length 28")
		Response(w, r, nil, http.StatusUnauthorized, err.Error())
		return
	}
	logs.Info("auth by open_id:", u.OpenID.String)

	var user models.VUser
	db := models.GetGorm()
	result := db.Where("open_id = ?", u.OpenID.String).First(&user)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}

	//创建新用户
	if user.ID <= 0 || result.Error == gorm.ErrRecordNotFound {
		logs.Info("new user, create with open_id: %s", u.OpenID.String)

		if u.MinappID.Int64 <= 0 {
			err = fmt.Errorf("must provide minapp_id")
			Response(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}

		if u.Nickname.String == "" || u.Avatar.String == "" {
			err = fmt.Errorf("must provide nickname and avatar")
			Response(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		if u.Phone.String != "" && !verifyPhone(u.Phone.String) {
			err = fmt.Errorf("%s not a valid phone number", u.Phone.String)
			Response(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		err = db.Create(&u).Error
		if err != nil {
			logs.Error(err)
			Response(w, r, nil, http.StatusInternalServerError, err.Error())
			return
		}
		logs.Info("创建用户：%v", u)
		result := db.Where("id = ?", u.ID).First(&user)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			logs.Error(err)
			Response(w, r, nil, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if user.MinappID.Int64 != u.MinappID.Int64 {
			err = fmt.Errorf("auth failed with wrong minapp_id")
			Response(w, r, nil, http.StatusUnauthorized, err.Error())
			return
		}
		err = invalidZeroNullValue(&u)
		if err != nil {
			logs.Error(err)
			Response(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		//更新用户信息
		if user.UnionID.String != u.UnionID.String || user.StuID.String != u.StuID.String ||
			user.Avatar.String != u.Avatar.String || user.Nickname.String != u.Nickname.String ||
			user.City.String != u.City.String || user.Province.String != u.Province.String ||
			user.Country.String != u.Country.String || user.Gender.Int64 != u.Gender.Int64 ||
			user.Language.String != u.Language.String || user.Phone.String != u.Phone.String {
			u.ID = user.ID
			result := db.Model(&u).Update(u)
			if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
				logs.Error(err)
				Response(w, r, nil, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	logs.Info(user)
	cookieStr, err := NewCookie(user.ID)
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Set-Cookie", cookieStr)
	Response(w, r, user, http.StatusOK, "")
	return
}

func AuthBySchool(w http.ResponseWriter, r *http.Request) {
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
	client, err := gzhu_jw.BasicAuthClient(username, password)
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusUnauthorized, err.Error())
		return
	}
	//将客户端存入缓存
	JWClients[username] = client

	logs.Info("用户：%s 接口：%s", username, r.URL.Path)
	Response(w, r, nil, http.StatusOK, "request ok")
}
