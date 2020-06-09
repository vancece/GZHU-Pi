package routers

import (
	"GZHU-Pi/env"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/mo7zayed/reqip"
	"gopkg.in/guregu/null.v3"
	"io/ioutil"
	"net/http"
	"time"
)

//使用open_id认证，不存在则创建新用户
func Auth(w http.ResponseWriter, r *http.Request) {

	if r.URL.Query().Get("type") == "gzhu" || r.URL.Query().Get("type") == "school" ||
		r.URL.Query().Get("school") != "" {
		AuthBySchool(w, r)
		return
	}
	if r.URL.Query().Get("type") == "bind_mp" {
		logs.Info("bind mp_open_id")
		BindMpOpenID(w, r)
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

	var u env.TUser
	err = json.Unmarshal(body, &u)
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}
	if u.OpenID.String == "" || len(u.OpenID.String) != 28 {
		err = fmt.Errorf("must give openid and with length 28")
		logs.Error(err)
		Response(w, r, nil, http.StatusUnauthorized, err.Error())
		return
	}
	logs.Info("auth by open_id:", u.OpenID.String)

	//防抖检测
	if isDebounce("auth:"+u.OpenID.String, 10*time.Second) {
		Response(w, r, nil, http.StatusOK, "debounce auth")
		return
	}

	var vUser env.VUser
	db := env.GetGorm()

	key := fmt.Sprintf("gzhupi:vuser:%s", u.OpenID.String)
	defer func() {
		if err != nil {
			_, _ = env.RedisCli.Del(key).Result()
		}
	}()
	//查询缓存
	val, err := env.RedisCli.Get(key).Result()
	if err != nil && err != redis.Nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	if err == redis.Nil {
		err = db.Where("open_id = ?", u.OpenID.String).First(&vUser).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logs.Error(err)
			Response(w, r, nil, http.StatusInternalServerError, err.Error())
			return
		}
		if err == nil && vUser.ID > 0 {
			//加入缓存
			logs.Debug("Set cache %s", key)
			buf, err := json.Marshal(&vUser)
			if err != nil {
				logs.Error(err)
				Response(w, r, nil, http.StatusInternalServerError, err.Error())
				return
			}
			err = env.RedisCli.Set(key, string(buf), 30*24*time.Hour).Err()
			if err != nil {
				logs.Error(err)
				Response(w, r, nil, http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			logs.Debug(u.OpenID.String, err)
		}
	} else {
		//解析缓存
		logs.Debug("Hit cache %s", key)
		err = json.Unmarshal([]byte(val), &vUser)
		if err != nil {
			logs.Error(err)
			Response(w, r, nil, http.StatusInternalServerError, err.Error())
			return
		}
	}

	//创建新用户
	if vUser.ID <= 0 || err == gorm.ErrRecordNotFound {
		logs.Info("new vUser, create with open_id: %v", u)

		if u.MinappID.Int64 <= 0 {
			err = fmt.Errorf("must provide minapp_id")
			logs.Error(err)
			Response(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		if u.Nickname.String == "" {
			t := fmt.Sprint(time.Now().UnixNano())
			u.Nickname = null.StringFrom(t[len(t)-4:])
		}
		if u.Nickname.String == "" || u.Avatar.String == "" {
			err = fmt.Errorf("must provide nickname and avatar")
			logs.Error(err)
			Response(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		if u.Phone.String != "" && !verifyPhone(u.Phone.String) {
			err = fmt.Errorf("%s not a valid phone number", u.Phone.String)
			logs.Error(err)
			Response(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		seed := vUser.StuID.String + vUser.OpenID.String + vUser.Nickname.String + fmt.Sprint(time.Now().Unix())
		u.ProfilePic = null.StringFrom(RandomAvatar(seed))

		err = db.Create(&u).Error
		if err != nil {
			u.ProfilePic.String = ""
			logs.Error(err, u)
			Response(w, r, nil, http.StatusInternalServerError, err.Error())
			return
		}
		logs.Info("创建用户：%v", u)
		err = db.Where("id = ?", u.ID).First(&vUser).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logs.Error(err)
			Response(w, r, nil, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if vUser.MinappID.Int64 != u.MinappID.Int64 {
			err = fmt.Errorf("auth failed with wrong minapp_id")
			logs.Error(err)
			Response(w, r, nil, http.StatusUnauthorized, err.Error())
			return
		}
		err = invalidZeroNullValue(&u)
		if err != nil {
			logs.Error(err)
			Response(w, r, nil, http.StatusBadRequest, err.Error())
			return
		}
		u.ID = vUser.ID
		//更新用户信息
		if (vUser.UnionID.String != u.UnionID.String || vUser.StuID.String != u.StuID.String ||
			vUser.Avatar.String != u.Avatar.String || vUser.Nickname.String != u.Nickname.String ||
			vUser.City.String != u.City.String || vUser.Province.String != u.Province.String ||
			vUser.Country.String != u.Country.String || vUser.Gender.Int64 != u.Gender.Int64 ||
			vUser.Language.String != u.Language.String || vUser.Phone.String != u.Phone.String) &&
			u.Nickname.String != "" {

			err = db.Model(&u).Update(u).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				logs.Error(err)
				Response(w, r, nil, http.StatusInternalServerError, err.Error())
				return
			}

			err = db.First(&vUser).Error
			if err != nil {
				logs.Error(err)
				Response(w, r, nil, http.StatusInternalServerError, err.Error())
				return
			}
			//更新缓存
			logs.Debug("Update cache %s", key)
			buf, err := json.Marshal(&vUser)
			if err != nil {
				logs.Error(err)
				Response(w, r, nil, http.StatusInternalServerError, err.Error())
				return
			}
			err = env.RedisCli.Set(key, string(buf), 30*24*time.Hour).Err()
			if err != nil {
				logs.Error(err)
				Response(w, r, nil, http.StatusInternalServerError, err.Error())
				return
			}
		}
		//创建随机头像
		if vUser.ProfilePic.String == "" {
			seed := vUser.StuID.String + vUser.OpenID.String + vUser.Nickname.String + fmt.Sprint(time.Now().Unix())
			go func() {
				u.ProfilePic = null.StringFrom(RandomAvatar(seed))
				err = db.Model(&u).Update(u).Error
				if err != nil && err != gorm.ErrRecordNotFound {
					logs.Error(err)
					return
				}
				logs.Debug("Del key : %s", key)
				_, err = env.RedisCli.Del(key).Result()
				if err != nil && err != gorm.ErrRecordNotFound {
					logs.Error(err)
					return
				}
			}()
		}
	}
	vUser.ProfilePic.String = ""
	logs.Info(vUser)
	cookieStr, err := NewCookie(vUser.ID)
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Set-Cookie", cookieStr)
	Response(w, r, vUser, http.StatusOK, "")
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
	client, err := newJWClient(r, username, password)
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusUnauthorized, err.Error())
		return
	}
	//将客户端存入缓存
	Jwxt.Store(getCacheKey(r, username), client)

	logs.Info("用户：%s IP: %s 接口：%s ", username, reqip.GetClientIP(r), r.URL.Path)
	Response(w, r, nil, http.StatusOK, "request ok")

	//绑定账号
	vUser, err := VUserByCookies(r)
	if err != nil {
		logs.Error(err)
		//Response(w, r, nil, http.StatusUnauthorized, err.Error())
		return
	}
	vUser.StuID = null.StringFrom(username)

	err = env.GetGorm().Model(&env.TUser{ID: vUser.ID}).UpdateColumns(env.TUser{StuID: vUser.StuID}).Error
	if err != nil {
		logs.Error(err)
		//Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}

	//更新缓存
	key := fmt.Sprintf("gzhupi:vuser:%s", vUser.OpenID.String)
	logs.Debug("Update cache %s", key)
	buf, err := json.Marshal(&vUser)
	if err != nil {
		logs.Error(err)
		//Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	err = env.RedisCli.Set(key, string(buf), 30*24*time.Hour).Err()
	if err != nil {
		logs.Error(err)
		//Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}

	//Response(w, r, vUser, http.StatusOK, "request ok")

}

func VUserByCookies(r *http.Request) (user *env.VUser, err error) {

	if len(r.Cookies()) == 0 {
		err = fmt.Errorf("微信授权信息无效，请退出小程序重新打开")
		logs.Error(err, r.URL)
		return
	}

	user = &env.VUser{}
	user.ID, err = GetUserID(r)
	if err != nil {
		if err.Error() == "Token is expired" {
			err = fmt.Errorf("登录信息过期，请重新打开小程序")
		}
		logs.Error(err)
		return
	}
	err = env.GetGorm().First(user).Error
	if err != nil {
		logs.Error(err)
		return
	}
	return
}

//绑定微信公众号openid
func BindMpOpenID(w http.ResponseWriter, r *http.Request) {

	mpOpenID := r.URL.Query().Get("mp_open_id")
	if mpOpenID == "" || len(mpOpenID) != 28 {
		err := fmt.Errorf("公众号openid无效: %s", mpOpenID)
		logs.Error(err)
		Response(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}
	vUser, err := VUserByCookies(r)
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusUnauthorized, err.Error())
		return
	}
	vUser.MpOpenID = null.StringFrom(mpOpenID)

	err = env.GetGorm().Model(&env.TUser{ID: vUser.ID}).UpdateColumns(env.TUser{MpOpenID: vUser.MpOpenID}).Error
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}

	//更新缓存
	key := fmt.Sprintf("gzhupi:vuser:%s", vUser.OpenID.String)
	logs.Debug("Update cache %s", key)
	buf, err := json.Marshal(&vUser)
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}
	err = env.RedisCli.Set(key, string(buf), 30*24*time.Hour).Err()
	if err != nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusInternalServerError, err.Error())
		return
	}

	Response(w, r, vUser, http.StatusOK, "")
}

//防抖检测，存在该key则返回true，否则设置key在指定时间过期
func isDebounce(key string, expiration time.Duration) (ok bool) {

	if expiration < 3 {
		expiration = 3
	}

	key = fmt.Sprintf("gzhupi:debounce:%s", key)
	_, err := env.RedisCli.Get(key).Result()
	if err != nil && err != redis.Nil {
		logs.Error(err)
		return
	}
	if err == redis.Nil {
		err = env.RedisCli.Set(key, "", expiration).Err()
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	logs.Warn("Hit debounce cache %s", key)
	return true
}
