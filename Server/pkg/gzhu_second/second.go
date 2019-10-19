/**
广州大学第二课堂学分系统客户端接口
*/

package gzhu_second

import (
	"crypto/tls"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	LoginError       = fmt.Errorf("认证失败，账号或密码错误")
	AuthError        = fmt.Errorf("认证失败，可能是缓存过期，请重试")
	urlencodedHeader = http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}
)

const baseUrl = "http://172.17.1.123"

var Urls = map[string]string{
	"second-login":  "https://cas.gzhu.edu.cn/cas_server/login?service=" + baseUrl + "/Login.aspx",
	"second-my":     baseUrl + "/XS/XMSB.aspx",
	"second-search": baseUrl + "/JWC/view.aspx",
	"second-detail": baseUrl + "/XS/XMSB_Detail_view.aspx?action=show&id=%s",
}

type SecondClient struct {
	Username  string
	Password  string
	ExpiresAt time.Time //客户端过期时间
	Client    *http.Client

	//记录查询页面的一些便当参数，下一次请求需要使用到
	VIEWSTATE          string
	VIEWSTATEGENERATOR string
	EVENTVALIDATION    string
}

func (c *SecondClient) doRequest(method, url string, header http.Header, body io.Reader) (*http.Response, error) {
	t1 := time.Now().UnixNano() / 1000000

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header[k] = v
	}
	resp, err := c.Client.Do(req)

	logs.Debug("请求耗时：", time.Now().UnixNano()/1000000-t1, "ms", url)
	return resp, err
}

func newClient(username, password string) *SecondClient {
	// Allocate a new cookie jar to mimic the browser behavior:
	cookieJar, _ := cookiejar.New(nil)
	// Fill up basic data:
	c := &SecondClient{
		Username: username,
		Password: password,
	}
	//设置客户端20分钟后过期
	c.ExpiresAt = time.Now().Add(20 * time.Minute)
	// Whn initializing the http.Client, copy default values from http.DefaultClient
	// Pass a pointer to the cookie jar that was created earlier:
	c.Client = &http.Client{
		CheckRedirect: http.DefaultClient.CheckRedirect,
		Jar:           cookieJar,
		Timeout:       http.DefaultClient.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	return c
}

func BasicAuthClient(username, password string) (client *SecondClient, err error) {
	if username == "" {
		return nil, fmt.Errorf("not init username or password")
	}
	c := newClient(username, password)
	//发送get请求，获取登录页面信息
	resp, err := c.doRequest("GET", Urls["second-login"], nil, nil)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	//提取登录表单
	r1, _ := regexp.Compile(`name="lt" value="(.*?)" />`)
	lt := r1.FindStringSubmatch(string(body))

	r2, _ := regexp.Compile(`name="execution" value="(.*?)" />`)
	execution := r2.FindStringSubmatch(string(body))

	if len(lt) < 2 || len(execution) < 2 {
		return nil, fmt.Errorf("get login form failed")
	}

	postValue := url.Values{
		"username":  {c.Username},
		"password":  {c.Password},
		"captcha":   {""},
		"warn":      {"true"},
		"lt":        {lt[1]},
		"execution": {execution[1]},
		"_eventId":  {"submit"},
		"submit":    {"登录"},
	}
	//编码表单
	postString := postValue.Encode()

	resp, err = c.doRequest("POST", Urls["second-login"], urlencodedHeader, strings.NewReader(postString))
	if err != nil {
		return nil, err
	}
	body, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	//判断是否登录成功

	ok, _ := regexp.MatchString("账号或密码错误", string(body))
	if ok {
		return nil, LoginError
	}
	return c, nil
}
