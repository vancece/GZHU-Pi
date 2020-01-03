package gzhu_jw

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Room struct {
	Bz     string `json:"bz" remark:"备注"`
	Cdbh   string `json:"cdbh" remark:"场地编号"`
	Cdjylx string `json:"cdjylx" remark:"场地借用类型"`
	CdlbID string `json:"cdlb_id" remark:"场地类别id"`
	Cdlbmc string `json:"cdlbmc" remark:"场地类别"`
	Cdmc   string `json:"cdmc" remark:"场地名称"`
	Jxlmc  string `json:"jxlmc" remark:"教学楼名称"`
	Kszws1 string `json:"kszws1" remark:"考试座位数"`
	Sydxmc string `json:"sydxmc" remark:"使用部门"`
	Xqmc   string `json:"xqmc" remark:"校区名称"`
	Zws    string `json:"zws" remark:"座位数"`
}

type RoomData struct {
	Items []*Room `json:"items"`
	Count int     `json:"count"`
}

func (c *JWClient) GetEmptyRoom(r *http.Request) (data *RoomData, err error) {

	if len(r.PostForm["jcd"]) < 1 || len(r.PostForm["zcd"]) < 1 || len(r.PostForm["xqm"]) < 1 {
		return nil, fmt.Errorf("illegal argument")
	}
	var xqm string
	if r.PostForm["xqm"][0] == "1" {
		xqm = "3"
	} else if r.PostForm["xqm"][0] == "2" {
		xqm = "12"
	} else {
		xqm = r.PostForm["xqm"][0]
	}

	nd := time.Now().Unix() * 1000 //时间戳
	form := url.Values{
		"xqh_id":                 r.PostForm["xqh_id"],                 // 校区号
		"xnm":                    r.PostForm["xnm"],                    // 学年名
		"xqm":                    {xqm},                                // 学期名
		"cdlb_id":                r.PostForm["cdlb_id"],                // 场地类别
		"qszws":                  r.PostForm["qszws"],                  // 最小座位号
		"jszws":                  r.PostForm["jszws"],                  // 最大座位号
		"cdmc":                   r.PostForm["cdmc"],                   // 场地名称
		"lh":                     r.PostForm["lh"],                     // 楼号
		"jcd":                    {powHandler(r.PostForm["jcd"][0])},   // 节次
		"queryModel.currentPage": r.PostForm["queryModel.currentPage"], // 前往页面数
		"nd":                     {strconv.Itoa(int(nd))},              // 生成时间戳
		"xqj":                    r.PostForm["xqj"],                    // 星期
		"zcd":                    {powHandler(r.PostForm["zcd"][0])},   // 周次

		//default form
		"fwzt":                 {"cx"},
		"cdejlb_id":            {""},
		"qssd":                 {""},
		"jssd":                 {""},
		"qssj":                 {""},
		"jssj":                 {""},
		"jyfs":                 {"0"},
		"cdjylx":               {""},
		"_search":              {"false"},
		"queryModel.showCount": {"30"},
		"queryModel.sortName":  {"cdbh"},
		"queryModel.sortOrder": {"asc"},
		"time":                 {"1"},
	}

	resp, err := c.doRequest("POST", Urls["empty-room"], urlencodedHeader, strings.NewReader(form.Encode()))
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
	rooms, total := ParseRoom(body)
	data = &RoomData{
		Items: rooms,
		Count: total,
	}
	return
}

//周次，节次累加求2的幂
func powHandler(target string) string {
	result := 0
	numList := strings.Split(target, ",")
	for _, v := range numList {
		n, err := strconv.Atoi(v)
		if err != nil {
			return ""
		}
		result = result + int(math.Pow(2, float64(n)-1))
	}
	res := strconv.Itoa(result)
	return res
}

func ParseRoom(body []byte) (rooms []*Room, total int) {

	rooms = []*Room{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	roomList := json.Get(body, "items")
	total = json.Get(body, "totalCount").ToInt()

	for i := 0; true; i++ {
		r := &Room{}
		r.Cdbh = roomList.Get(i).Get("cdbh").ToString()
		if r.Cdbh == "" {
			break
		}
		r.Bz = roomList.Get(i).Get("bz").ToString()
		r.Cdjylx = roomList.Get(i).Get("cdjylx").ToString()
		r.CdlbID = roomList.Get(i).Get("cdlb_id").ToString()
		r.Cdlbmc = roomList.Get(i).Get("cdlbmc").ToString()
		r.Cdmc = roomList.Get(i).Get("cdmc").ToString()
		r.Jxlmc = roomList.Get(i).Get("jxlmc").ToString()
		r.Kszws1 = roomList.Get(i).Get("kszws1").ToString()
		//r.Sydxmc = roomList.Get(i).Get("sydxmc").ToString()
		r.Xqmc = roomList.Get(i).Get("xqmc").ToString()
		r.Zws = roomList.Get(i).Get("zws").ToString()

		rooms = append(rooms, r)
	}
	return
}
