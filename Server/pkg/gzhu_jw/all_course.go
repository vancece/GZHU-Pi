package gzhu_jw

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type RawCourse struct {
	Cdbh    string `json:"cdbh"`
	Cdlbmc  string `json:"cdlbmc"`
	Cdmc    string `json:"cdmc"`
	Cdqsjsz string `json:"cdqsjsz"`
	Cdskjc  string `json:"cdskjc"`
	Jgh     string `json:"jgh"`
	JghID   string `json:"jgh_id"`
	Jslxdh  string `json:"jslxdh"`
	Jsxy    string `json:"jsxy"`
	JxbID   string `json:"jxb_id"`
	Jxbmc   string `json:"jxbmc"`
	Jxbrs   int    `json:"jxbrs"`
	Jxbzc   string `json:"jxbzc"`
	Jxdd    string `json:"jxdd"`
	Jxlmc   string `json:"jxlmc"`
	Kch     string `json:"kch"`
	KchID   string `json:"kch_id"`
	Kcmc    string `json:"kcmc"`
	Kcxzmc  string `json:"kcxzmc"`
	KkbmID  string `json:"kkbm_id"`
	Kkxy    string `json:"kkxy"`
	Qsjsz   string `json:"qsjsz"`
	Rwzxs   string `json:"rwzxs"`
	Skjc    string `json:"skjc"`
	Sksj    string `json:"sksj"`
	Xbmc    string `json:"xbmc"`
	Xf      string `json:"xf"`
	Xkrs    int    `json:"xkrs"`
	Xm      string `json:"xm"`
	Xnm     string `json:"xnm"`
	Xn      string `json:"xn"`
	Xq      string `json:"xq"`
	XqhID   string `json:"xqh_id"`
	Xqj     int    `json:"xqj"`
	Xqm     string `json:"xqm"`
	Xqmc    string `json:"xqmc"`
	Zcmc    string `json:"zcmc"`
	Zgxl    string `json:"zgxl"`
	Zhxs    string `json:"zhxs"`
	Zjxh    string `json:"zjxh"`
	Zyzc    string `json:"zyzc"`
	Zcd     int    `json:"zcd"`
	Jc      int    `json:"jc"`
	Cdjc    string `json:"cdjc"`
	Zws     int    `json:"zws"`
	Lch     int    `json:"lch"`
}

//查询全校课表
func (c *JWClient) SearchAllCourse(xnm, xqm string, page, count int) (data []RawCourse, csvData []byte, err error) {

	if xnm == "" {
		year := time.Now().Year()
		month := time.Now().Month()
		if month < 8 {
			year = year - 1
		}
		xnm = fmt.Sprint(year)
	}

	if page <= 0 {
		page = 1
	}
	if count <= 0 {
		count = 15
	}

	if xqm == "1" || xqm == "3" {
		xqm = "3"
	} else {
		xqm = "12"
	}

	nd := time.Now().Unix() * 1000 //时间戳
	var form = url.Values{
		"xnm":                    {xnm}, //2019
		"xqm":                    {xqm}, //3 是第一学期，12 是第二学期
		"_search":                {"false"},
		"nd":                     {strconv.Itoa(int(nd))},
		"queryModel.showCount":   {strconv.Itoa(count)},
		"queryModel.currentPage": {strconv.Itoa(page)},
		"queryModel.sortName":    {""},
		"queryModel.sortOrder":   {"asc"},
	}

	resp, err := c.doRequest("POST", Urls["all-course"], urlencodedHeader, strings.NewReader(form.Encode()))
	if err != nil {
		logs.Error(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	//检查登录状态
	if strings.Contains(string(body), "登录") {
		return nil, nil, AuthError
	}

	type Data struct {
		Items []RawCourse `json:"items"`
	}
	var d Data

	json1 := jsoniter.ConfigCompatibleWithStandardLibrary
	err = json1.Unmarshal(body, &d)
	if err != nil {
		logs.Error(err)
		return
	}
	data = d.Items

	csvData = ToCsvFormat(data)

	return
}

func ToCsvFormat(all []RawCourse) (data []byte) {

	header := []string{"xq", "cdlbmc", "cdmc", "cdqsjsz", "cdskjc", "jgh", "jgh_id", "jslxdh", "jsxy",
		"jxb_id", "jxbmc", "jxbrs", "jxbzc", "jxdd", "jxlmc", "kch", "kch_id", "kcmc", "kcxzmc", "kkxy",
		"rwzxs", "skjc", "sksj", "xbmc", "xf", "xkrs", "xm", "xnm", "xn", "xqh_id", "xqj", "xqmc", "zcmc",
		"zgxl", "zhxs", "zyzc"}

	var lines []string
	lines = append(lines, strings.Join(header, ","))

	for _, v := range all {
		var values []string
		values = append(values, v.Xq, v.Cdlbmc, v.Cdmc, v.Cdqsjsz, v.Cdskjc, v.Jgh, v.JghID, v.Jslxdh, v.Jsxy,
			v.JxbID, v.Jxbmc, fmt.Sprint(v.Jxbrs), v.Jxbzc, v.Jxdd, v.Jxlmc, v.Kch, v.KchID, v.Kcmc, v.Kcxzmc,
			v.Kkxy, v.Rwzxs, v.Skjc, v.Sksj, v.Xbmc, v.Xf, fmt.Sprint(v.Xkrs), v.Xm, v.Xnm, v.Xn, v.XqhID,
			fmt.Sprint(v.Xqj), v.Xqmc, v.Zcmc, v.Zgxl, v.Zhxs, v.Zyzc)

		for k, v := range values {
			if strings.Contains(v, ",") {
				values[k] = fmt.Sprintf(`"%s"`, v) //去除字符串内部逗号对csv的影响
			}
		}

		csvLine := strings.Join(values, ",")
		lines = append(lines, csvLine)
	}
	csv := strings.Join(lines, "\n")

	data = []byte(csv)
	return
}
