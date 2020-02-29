package routers

import (
	"GZHU-Pi/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

//使用open_id认证，不存在则创建新用户
func Auth(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logs.Error(err)
		return
	}
	defer r.Body.Close()
	if len(body) == 0 {
		err = fmt.Errorf("Call api by post with empty body ")
		logs.Error(err)
		return
	}
	var u models.TUser
	err = json.Unmarshal(body, &u)
	if err != nil {
		logs.Error(err)
		return
	}

	if u.OpenID.String == "" || len(u.OpenID.String) != 28 {
		err = fmt.Errorf("must give openid and with length 28")
		Response(w, r, nil, http.StatusUnauthorized, err.Error())
		return
	}
	logs.Info("auth by open_id:", u.OpenID.String)

	var user models.TUser
	db := models.GetGorm()
	result := db.Where("open_id = ?", u.OpenID.String).First(&user)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	if user.ID <= 0 {
		logs.Info("new user, create with open_id: %s", u.OpenID.String)

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
		user = u
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
