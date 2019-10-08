package routers

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
)

// token 为微信开发者平台上提交的字符串
const token = "token"

func WeChatCheck(w http.ResponseWriter, r *http.Request) {

	logs.Info("Url : %s", r.URL)

	_ = r.ParseForm()
	if len(r.Form) > 0 {
		for k, v := range r.Form {
			logs.Debug("%s=%s", k, v[0])
		}
	}
	ok := checkSignature(r.Form)
	if ok {
		logs.Info("check signature success")
		_, _ = w.Write([]byte(r.Form.Get("echostr")))
	} else {
		logs.Info("check signature fail")
		Response(w, r, nil, http.StatusUnauthorized, "check signature fail")
	}
}

//校验签名，检查请求是否来自微信后台
func checkSignature(values url.Values) bool {

	array := []string{values.Get("nonce"), values.Get("timestamp"), token}

	sort.Sort(sort.StringSlice(array))

	tmpStr := array[0] + array[1] + array[2]
	data := []byte( tmpStr )
	hash := sha1.Sum(data)
	hashCode := hex.EncodeToString(hash[0:20])
	if hashCode != values.Get("signature") {
		return false
	}
	return true
}


func WeChat(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	logs.Info(string(body), err)
	w.Write([]byte("success"))
}
