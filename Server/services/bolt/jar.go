package bolt

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	dbname    = "cookie.db"
	bucket    = "cookieBucket"
	cookieKey = "jw_cookie:%s"
)

var cookieDb *BlotClient

func init() {
	cookieDb = newBlotClient(dbname, bucket)
	if cookieDb == nil {
		log.Fatal("init cookie boltdb failed")
	}
}

//保存cookieJar
func SaveJar(jar http.CookieJar, username, rawUrl string) (err error) {
	if jar == nil || username == "" || rawUrl == "" {
		return fmt.Errorf("invalid arguments")
	}
	u, err := url.Parse(rawUrl)
	if err != nil {
		logs.Error(err)
		return
	}
	cookies := jar.Cookies(u)
	data, err := json.Marshal(cookies)
	if err != nil {
		logs.Error(err)
		return
	}
	key := fmt.Sprintf(cookieKey, username)

	err = cookieDb.Put([]byte(key), data)
	if err != nil {
		logs.Error(err)
		return
	}
	return
}

//提取恢复cookieJar
func GetJar(username string, rawUrl string) (jar http.CookieJar, err error) {
	if username == "" || rawUrl == "" {
		return nil, fmt.Errorf("invalid arguments")
	}

	//读取缓存cookie
	key := fmt.Sprintf(cookieKey, username)
	data, err := cookieDb.Get(key)
	if err != nil {
		logs.Error(err)
		return
	}

	u, err := url.Parse(rawUrl)
	if err != nil {
		logs.Error(err)
		return
	}
	//解析cookie
	c := &[]*http.Cookie{}
	err = json.Unmarshal(data, c)
	if err != nil {
		logs.Error(err)
		return
	}
	//创建新jar
	jar, _ = cookiejar.New(nil)
	jar.SetCookies(u, *c)

	return
}
