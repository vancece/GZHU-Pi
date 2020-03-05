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

type CourseExport struct {
	Xq      string `json:"xq"`
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
	Kkxy    string `json:"kkxy"`
	Rwzxs   string `json:"rwzxs"`
	Skjc    string `json:"skjc"`
	Sksj    string `json:"sksj"`
	Xbmc    string `json:"xbmc"`
	Xf      string `json:"xf"`
	Xkrs    int    `json:"xkrs"`
	Xm      string `json:"xm"`
	Xnm     string `json:"xnm"`
	Xn      string `json:"xn"`
	XqhID   string `json:"xqh_id"`
	Xqj     int    `json:"xqj"`
	Xqmc    string `json:"xqmc"`
	Zcmc    string `json:"zcmc"`
	Zgxl    string `json:"zgxl"`
	Zhxs    string `json:"zhxs"`
	Zyzc    string `json:"zyzc"`
}

//查询全校课表
func (c *JWClient) SearchAllCourse(xnm, xqm string, page, count int) (data []CourseExport, csvData []byte, err error) {

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
	data = ParseAllCourse(body)
	csvData = ToCsvFormat(data)
	return
}

func ParseAllCourse(body []byte) (all []CourseExport) {

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	courseList := json.Get(body, "items")

	for i := 0; true; i++ {
		v := CourseExport{}
		v.KchID = courseList.Get(i).Get("kch_id").ToString()
		if v.KchID == "" {
			break
		}
		v.Xq = courseList.Get(i).Get("xq").ToString()
		v.Cdlbmc = courseList.Get(i).Get("cdlbmc").ToString()
		v.Cdmc = courseList.Get(i).Get("cdmc").ToString()
		v.Cdqsjsz = courseList.Get(i).Get("cdqsjsz").ToString()
		v.Cdskjc = courseList.Get(i).Get("cdskjc").ToString()
		v.Jgh = courseList.Get(i).Get("jgh").ToString()
		v.JghID = courseList.Get(i).Get("jgh_id").ToString()
		v.Jslxdh = courseList.Get(i).Get("jslxdh").ToString()
		v.Jsxy = courseList.Get(i).Get("jsxy").ToString()
		v.JxbID = courseList.Get(i).Get("jxb_id").ToString()
		v.Jxbmc = courseList.Get(i).Get("jxbmc").ToString()
		v.Jxbrs = courseList.Get(i).Get("jxbrs").ToInt()
		v.Jxbzc = courseList.Get(i).Get("jxbzc").ToString()
		v.Jxdd = courseList.Get(i).Get("jxdd").ToString()
		v.Jxlmc = courseList.Get(i).Get("jxlmc").ToString()
		v.Kch = courseList.Get(i).Get("kch").ToString()
		v.KchID = courseList.Get(i).Get("kch_id").ToString()
		v.Kcmc = courseList.Get(i).Get("kcmc").ToString()
		v.Kcxzmc = courseList.Get(i).Get("kcxzmc").ToString()
		v.Kkxy = courseList.Get(i).Get("kkxy").ToString()
		v.Rwzxs = courseList.Get(i).Get("rwzxs").ToString()
		v.Skjc = courseList.Get(i).Get("skjc").ToString()
		v.Sksj = courseList.Get(i).Get("sksj").ToString()
		v.Xbmc = courseList.Get(i).Get("xbmc").ToString()
		v.Xf = courseList.Get(i).Get("xf").ToString()
		v.Xkrs = courseList.Get(i).Get("xkrs").ToInt()
		v.Xm = courseList.Get(i).Get("xm").ToString()
		v.Xnm = courseList.Get(i).Get("xnm").ToString()
		v.Xn = courseList.Get(i).Get("xn").ToString()
		v.XqhID = courseList.Get(i).Get("xqh_id").ToString()
		v.Xqj = courseList.Get(i).Get("xqj").ToInt()
		v.Xqmc = courseList.Get(i).Get("xqmc").ToString()
		v.Zcmc = courseList.Get(i).Get("zcmc").ToString()
		v.Zgxl = courseList.Get(i).Get("zgxl").ToString()
		v.Zhxs = courseList.Get(i).Get("zhxs").ToString()
		v.Zyzc = courseList.Get(i).Get("zyzc").ToString()

		all = append(all, v)
	}
	return
}

func ToCsvFormat(all []CourseExport) (data []byte) {

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
