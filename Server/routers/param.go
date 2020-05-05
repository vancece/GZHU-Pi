/**
 * @File: param
 * @Author: Shaw
 * @Date: 2020/5/5 5:26 PM
 * @Desc

 */

package routers

import (
	"GZHU-Pi/env"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
	"net/http"
	"time"
)

func Param(w http.ResponseWriter, r *http.Request) {

	var user env.TUser
	var err error
	if len(r.Cookies()) == 0 {
		Response(w, r, nil, http.StatusOK, "request ok")
		return
	}
	user.ID, err = GetUserID(r)
	if err != nil || user.ID <= 0 {
		Response(w, r, nil, http.StatusUnauthorized, fmt.Sprintf("%v", err))
		return
	}

	key := fmt.Sprintf("gzhupi:param:action:user:%d", user.ID)
	_, err = env.RedisCli.Get(key).Result()
	if err != nil && err != redis.Nil {
		logs.Error(err)
		Response(w, r, nil, http.StatusUnauthorized, err.Error())
		return
	}
	var data map[string]interface{}
	if err == redis.Nil {
		data, err = getActionParam()
		var store = map[string]interface{}{
			"timestamp": time.Now().Unix() * 1000,
			"data":      data,
			"times":     1,
		}
		js, err := json.Marshal(store)
		if err != nil {
			logs.Error(err)
			Response(w, r, nil, http.StatusInternalServerError, err.Error())
			return
		}
		err = env.RedisCli.Set(key, string(js), 0).Err()
		if err != nil {
			Response(w, r, nil, http.StatusInternalServerError, err.Error())
			logs.Error(err)
			return
		}
	}

	Response(w, r, data, http.StatusOK, "request ok")

}

func getActionParam() (data map[string]interface{}, err error) {
	key := fmt.Sprintf("gzhupi:param:action")
	val, err := env.RedisCli.Get(key).Result()
	if err != nil && err != redis.Nil {
		logs.Error(err)
		return
	}
	if err == redis.Nil {
		//example := map[string]interface{}{
		//
		//	"modal": map[string]interface{}{
		//		"valid":   false,
		//		"cancel":  false,
		//		"confirm": false,
		//
		//		"img":       "",
		//		"title":     "活动推荐",
		//		"sub_title": "义工",
		//		"content":   []string{"经历了一个漫长的假期,你一定想出去走走，外面的世界那么大"},
		//		"btn_text":  "我想去看看",
		//		"nav_to":    "/pages/Setting/webview/webview?src=https://mp.weixin.qq.com/s/NbVEBpvlPgpfbAf_tRh09w",
		//	},
		//}
		//js, _ := json.Marshal(example)
		//err = env.RedisCli.Set(key, string(js), 0).Err()
		//if err != nil {
		//	logs.Error(err)
		//	return
		//}
		//logs.Warn("action param is not set")
		return
	}
	data = make(map[string]interface{})
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		logs.Error(err)
		return
	}

	return
}
