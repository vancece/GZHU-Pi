/**
广州大学教务系统客户端接口
*/

package gzhu_jw

import (
	"crypto/tls"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/rpc"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type JWClient struct {
	Username  string
	Password  string
	ExpiresAt time.Time //客户端过期时间
	Client    *http.Client
}

func (c *JWClient) doRequest(method, url string, header http.Header, body io.Reader) (*http.Response, error) {
	t1 := time.Now()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header[k] = v
	}
	resp, err := c.Client.Do(req)

	logs.Debug("请求耗时：", time.Since(t1), url)
	return resp, err
}

func newClient(username, password string) *JWClient {
	// Allocate a new cookie jar to mimic the browser behavior:
	cookieJar, _ := cookiejar.New(nil)
	// Fill up basic data:
	c := &JWClient{
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

func BasicAuthClient(username, password string) (client *JWClient, err error) {
	if username == "" {
		return nil, fmt.Errorf("not init username or password")
	}
	c := newClient(username, password)
	//发送get请求，获取登录页面信息
	resp, err := c.doRequest("GET", Urls["jw-login"], nil, nil)
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

	captcha := c.GetCaptcha()
	if captcha == "" {
		logs.Debug("验证码识别为空")
	}

	postValue := url.Values{
		"username":  {c.Username},
		"password":  {c.Password},
		"captcha":   {captcha},
		"warn":      {"true"},
		"lt":        {lt[1]},
		"execution": {execution[1]},
		"_eventId":  {"submit"},
		"submit":    {"登录"},
	}
	//编码表单
	postString := postValue.Encode()

	resp, err = c.doRequest("POST", Urls["jw-login"], urlencodedHeader, strings.NewReader(postString))
	if err != nil {
		return nil, err
	}
	body, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	//判断是否登录成功

	ok, _ := regexp.MatchString("验证码不正确", string(body))
	if ok {
		return nil, fmt.Errorf("验证码识别不正确，请重试")
	}
	ok, _ = regexp.MatchString("必须输入验证码", string(body))
	if ok {
		return nil, fmt.Errorf("验证码识别失败，请重试")
	}
	ok, _ = regexp.MatchString("账号或密码错误", string(body))
	if ok {
		return nil, LoginError
	}
	ok, _ = regexp.MatchString(`用户名[\s\S]*密码`, string(body))
	if ok {
		return nil, fmt.Errorf("不知为啥，就是没有登录进去，请告知开发者")
	}
	return c, nil
}

var rpcClient *rpc.Client

func (c *JWClient) GetCaptcha() (capture string) {

	var err error
	if rpcClient == nil {
		rpcClient, err = rpc.DialHTTP("tcp", "ifeel.vip:7201")
		if err != nil {
			logs.Error(err)
			return
		}
	}

	const maxTry = 3
	for range [maxTry]int{} {

		resp, err := c.doRequest("GET", "https://cas.gzhu.edu.cn/cas_server/captcha.jsp", nil, nil)
		if err != nil {
			logs.Error(err)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logs.Error(err)
			return
		}

		err = rpcClient.Call("OcrService.Capture", body, &capture)
		if err != nil {
			logs.Error(err)
			return
		}

		reg := regexp.MustCompile(`\d+`)
		res := reg.FindStringSubmatch(capture)
		if len(res) == 0 || len(res[0]) != 4 {
			logs.Debug("capture ocr failed: ", res)
			continue
		}
		capture = res[0]
		break
	}
	logs.Info("验证码：", capture)
	return
}
