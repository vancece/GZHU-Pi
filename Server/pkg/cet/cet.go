/**
 * @File: cet
 * @Author: Shaw
 * @Date: 2020/5/27 12:52 AM
 * @Desc

 */

package cet

import (
	"crypto/tls"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var cache = sync.Map{}

func NewCetClient(id, name, captcha string) *CetClient {
	c, ok := cache.Load(id)
	if ok {
		cli := c.(*CetClient)
		cli.Name = name
		cli.Captcha = captcha
		return cli
	}
	cli := newClient(id, name)
	cli.Captcha = captcha
	cache.Store(id, cli)
	return cli
}

type CetClient struct {
	Client  *http.Client `json:"-"`
	Captcha string       `json:"captcha"` //验证码

	Title      string `json:"title"`       //考试标题
	ID         string `json:"id"`          //准考证
	Name       string `json:"name"`        //姓名
	School     string `json:"school"`      //学校
	Total      int    `json:"total"`       //总分
	Listening  int    `json:"listening"`   //听力
	Reading    int    `json:"reading"`     //阅读
	Writing    int    `json:"writing"`     //写作
	VoiceLevel string `json:"voice_level"` //口语等级
}

var cetHeader = http.Header{
	"Referer": []string{"http://cet.neea.edu.cn/cet/"},
}

var transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

func newClient(id, name string) *CetClient {
	cookieJar, _ := cookiejar.New(nil)
	c := &CetClient{
		ID:   id,
		Name: name,
	}
	c.Client = &http.Client{
		Jar:       cookieJar,
		Timeout:   time.Minute,
		Transport: transport,
	}
	return c
}

func (this *CetClient) doRequest(method, url string, header http.Header, body io.Reader) (*http.Response, error) {
	t1 := time.Now()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header[k] = v
	}
	resp, err := this.Client.Do(req)

	logs.Debug("请求耗时：", time.Since(t1), url)
	return resp, err
}

func (this *CetClient) GetCaptcha() (imgUrl string, err error) {

	u := fmt.Sprintf("http://cache.neea.edu.cn/Imgs.do?c=CET&ik=%s&t=%.16f", this.ID, rand.Float64())

	resp, err := this.doRequest("GET", u, cetHeader, nil)
	if err != nil {
		logs.Error(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)

	//获取验证码id
	reg := regexp.MustCompile(`result.imgs\("(.*)"\)`)
	match := reg.FindStringSubmatch(string(body))
	if len(match) != 2 || match[1] == "" {
		err = fmt.Errorf("验证码获取 match failed: %s", string(body))
		logs.Error(err)
		return
	}

	imgUrl = fmt.Sprintf("http://cet.neea.edu.cn/imgs/%s.png", match[1])

	return
}

func (this *CetClient) GetCetInfo() (err error) {

	if len(this.ID) < 15 {
		err = fmt.Errorf("准考证号小于15位：%s", this.ID)
		logs.Error(err)
		return
	}

	level := "4"
	this.Title = "全国大学英语四级考试(CET4)"
	if this.ID[9] == '2' {
		level = "6"
		this.Title = "全国大学英语六级考试(CET6)"
	}

	year := this.ID[6:9] //年份场次 如192
	examCode := fmt.Sprintf("CET%s_%s_DANGCI", level, year)

	form := url.Values{
		"data": []string{fmt.Sprintf("%s,%s,%s", examCode, this.ID, this.Name)},
		"v":    []string{this.Captcha},
	}

	u := "http://cache.neea.edu.cn/cet/query?" + form.Encode()

	resp, err := this.doRequest("GET", u, cetHeader, nil)
	if err != nil {
		logs.Error(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	text := string(body)
	if !strings.Contains(text, this.ID) {
		text = strings.ReplaceAll(text, "result.callback({", "")
		text = strings.ReplaceAll(text, "});", "")
		err = fmt.Errorf(text)
		logs.Error(err)
		return
	}
	reg := regexp.MustCompile(`result.callback\({(.*)}\);`)
	match := reg.FindStringSubmatch(text)
	if len(match) != 2 || match[1] == "" {
		err = fmt.Errorf("成绩查询失败 match failed: %s", text)
		logs.Error(err)
		return
	}

	//logs.Info(match[1])
	//形如 z:'准考证',n:'姓名',x:'学校',s:.00,t:0,id:'',l:0,r:0,w:0,kyz:'--',kys:'口语等级'

	res := strings.Split(match[1], ",")

	for _, v := range res {

		sp := strings.Split(v, ":")
		if len(sp) != 2 {
			err = fmt.Errorf("获取成绩信息失败: %s", v)
			logs.Error(err)
			return
		}
		val := strings.ReplaceAll(sp[1], "'", "")

		switch {
		case strings.Contains(v, "x"):
			this.School = val
		case strings.Contains(v, "kys"):
			this.VoiceLevel = val
		case strings.Contains(v, "t"):
			this.Total, err = strconv.Atoi(val)
		case strings.Contains(v, "l"):
			this.Listening, err = strconv.Atoi(val)
		case strings.Contains(v, "r"):
			this.Reading, err = strconv.Atoi(val)
		case strings.Contains(v, "w"):
			this.Writing, err = strconv.Atoi(val)
		}
		if err != nil {
			logs.Error(err)
			return
		}
	}
	if this.School == "" {
		err = fmt.Errorf("查询成绩失败")
		logs.Error(err, text)
		return
	}
	return
}

func (this *CetClient) DelCache() {
	cache.Delete(this.ID)
}
