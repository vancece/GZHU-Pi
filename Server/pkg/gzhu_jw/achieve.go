package gzhu_jw

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

//学生学业情况数据

type Achieve struct {
	Type     string  `json:"type" remark:"课程类型"`
	Required string  `json:"required" remark:"要求学分"`
	Acquired string  `json:"acquired" remark:"获得学分"`
	Remained string  `json:"remained" remark:"未获得学分"`
	Node     bool    `json:"node" remark:"是否中间节点"`
	Items    []*Item `json:"items" remark:"课程列表"`

	FormID string `json:"-"` //请求表单id
}

type Item struct {
	CJ string  `json:"CJ" remark:"成绩"`
	JD float64 `json:"JD" remark:"绩点"`

	//JYXDXNM  string `json:"JYXDXNM" remark:"建议修读学年"`
	//JYXDXNMC string `json:"JYXDXNMC" remark:"建议修读学年度"`
	//JYXDXQM  string `json:"JYXDXQM" remark:"建议修读学期代码"`
	//JYXDXQMC string `json:"JYXDXQMC" remark:"建议修读学期"`

	KCH    string  `json:"KCH" remark:"课程号"`
	KCHID  string  `json:"KCH_ID" remark:"课程号id"`
	KCLBDM string  `json:"KCLBDM" remark:"课程类别代码"`
	KCLBMC string  `json:"KCLBMC" remark:"课程类别"`
	KCMC   string  `json:"KCMC" remark:"课程名称"`
	KCXZMC string  `json:"KCXZMC" remark:"课程性质"`
	KCYWMC string  `json:"KCYWMC" remark:"课程英文名称"`
	KCZT   float64 `json:"KCZT" remark:"课程状态？"`
	MAXCJ  string  `json:"MAXCJ" remark:"最大成绩"`
	SFJHKC string  `json:"SFJHKC" remark:""`
	XDZT   string  `json:"XDZT" remark:"修读状态"` //在修1、未过2、未修3、已修4
	XF     string  `json:"XF" remark:"学分"`
	XNM    string  `json:"XNM" remark:"学年"`
	XNMC   string  `json:"XNMC" remark:"学年度"`
	XQM    string  `json:"XQM" remark:"学期代码"`
	XQMMC  string  `json:"XQMMC" remark:"学期"`
	XSXXXX string  `json:"XSXXXX" remark:"学时备注"`
}

//学业情况查询
func (c *JWClient) GetAchieve() (achieves []*Achieve, err error) {

	resp, err := c.doRequest("GET", Urls["achieve-get"], urlencodedHeader, nil)
	if err != nil {
		logs.Error(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	//检查登录状态
	if strings.Contains(string(body), "登录") {
		return nil, AuthError
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logs.Error(err)
		return
	}
	//从html页面提取表单（表单未完整）
	var form = url.Values{}
	doc.Find("[type=hidden]").Each(func(i int, selection *goquery.Selection) {
		id, ok := selection.Attr("id")
		if !ok {
			return
		}
		value, ok := selection.Attr("value")
		if !ok {
			return
		}
		form[id] = []string{value}
	})
	achieves = GetOverViewInfo(body)
	//并发请求各类型课程列表,注意指针变量的传递
	var wg = sync.WaitGroup{}
	for _, v := range achieves {

		form["xfyqjd_id"] = []string{v.FormID}
		postValue := form.Encode()

		wg.Add(1)

		go func(v *Achieve) {
			resp, err := c.doRequest("POST", Urls["achieve-post"], urlencodedHeader, strings.NewReader(postValue))
			if err != nil {
				logs.Error(err)
				wg.Done()
				return
			}
			body, _ = ioutil.ReadAll(resp.Body)
			v.Items = ParseOverView(body)
			wg.Done()
		}(v)

	}
	wg.Wait()
	return
}

//从html页面提取各类别课程修读汇总信息，以及表单id
func GetOverViewInfo(body []byte) (achieves []*Achieve) {

	achieves = []*Achieve{}

	//提取类别代码，用于提交表单
	r1 := regexp.MustCompile(`xfyqjd_id='(.*?)' jdkcsx`)
	//课程类别名称
	r2 := regexp.MustCompile(`"(.*?)&nbsp;要求学分`)
	//学分要求
	r3 := regexp.MustCompile(`要求学分:([\d.]+)&nbsp;`)
	//获得学分
	r4 := regexp.MustCompile(`获得学分:([\d.]+)&nbsp;`)
	//未获得学分
	r5 := regexp.MustCompile(`未获得学分:([\d.]+)&nbsp;`)

	r := regexp.MustCompile(`</span></p><span`)
	indexes := r.FindAllIndex(body, -1)

	bodyStr := string(body)
	for _, v := range indexes {
		if v[0] < 500 {
			return
		}
		//提取正则r匹配出来的索引前几百个字符作为一个片段
		//正常是12个段落对应12g蓝色方块
		index := v[0]
		section := bodyStr[index-500 : index]

		var ov = &Achieve{}

		res1 := r1.FindStringSubmatch(section)
		if len(res1) >= 2 {
			ov.FormID = res1[1]
		}
		res2 := r2.FindStringSubmatch(section)
		if len(res2) >= 2 {
			ov.Type = res2[1]
		}
		res3 := r3.FindStringSubmatch(section)
		if len(res3) >= 2 {
			ov.Required = res3[1]
		}
		res4 := r4.FindStringSubmatch(section)
		if len(res4) >= 2 {
			ov.Acquired = res4[1]
		}
		res5 := r5.FindStringSubmatch(section)
		if len(res5) >= 2 {
			ov.Remained = res5[1]
		}
		//包含该字符串，认定为可以展开的中间节点
		if strings.Contains(section, "xfyqzjdgx='1'") {
			ov.Node = true
		}
		achieves = append(achieves, ov)
	}
	return
}

//解析学业情况信息
func ParseOverView(body []byte) (items []*Item) {

	items = []*Item{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	itemList := json.Get(body)
	//遍历提取
	for i := 0; true; i++ {
		e := &Item{}
		e.KCHID = itemList.Get(i).Get("KCH_ID").ToString()
		if e.KCHID == "" {
			break
		}
		e.CJ = itemList.Get(i).Get("CJ").ToString()
		e.JD = itemList.Get(i).Get("JD").ToFloat64()

		e.KCH = itemList.Get(i).Get("KCH").ToString()
		e.KCLBDM = itemList.Get(i).Get("KCLBDM").ToString()
		e.KCLBMC = itemList.Get(i).Get("KCLBMC").ToString()
		e.KCMC = itemList.Get(i).Get("KCMC").ToString()
		e.KCXZMC = itemList.Get(i).Get("KCXZMC").ToString()
		e.KCYWMC = itemList.Get(i).Get("KCYWMC").ToString()
		e.KCZT = itemList.Get(i).Get("KCZT").ToFloat64()
		e.MAXCJ = itemList.Get(i).Get("MAXCJ").ToString()
		e.SFJHKC = itemList.Get(i).Get("SFJHKC").ToString()
		e.XDZT = itemList.Get(i).Get("XDZT").ToString()
		e.XF = itemList.Get(i).Get("XF").ToString()
		e.XSXXXX = itemList.Get(i).Get("XSXXXX").ToString()

		e.XNM = itemList.Get(i).Get("XNM").ToString()
		e.XNMC = itemList.Get(i).Get("XNMC").ToString()
		e.XQM = itemList.Get(i).Get("XQM").ToString()
		e.XQMMC = itemList.Get(i).Get("XQMMC").ToString()

		if e.XNM == "" {
			e.XNM = itemList.Get(i).Get("JYXDXNM").ToString()
			e.XNMC = itemList.Get(i).Get("JYXDXNMC").ToString()
			e.XQM = itemList.Get(i).Get("JYXDXQM").ToString()
			e.XQMMC = itemList.Get(i).Get("JYXDXQMC").ToString()
		}
		items = append(items, e)
	}
	return
}
