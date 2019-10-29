package gzhu_second

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type SecondItem struct {
	ID string `json:"id"  remark:"项目id"`

	StuID      string `json:"stu_id" remark:"学号"`
	StuName    string `json:"stu_name" remark:"姓名"`
	MajorClass string `json:"major_class" remark:"专业班级"`
	ApplyYear  string `json:"apply_year" remark:"申报学年"`

	Year  string `json:"year" remark:"项目年度"`
	Type  string `json:"type" remark:"项目类型"`
	Name  string `json:"name" remark:"项目名称"`
	Level string `json:"level" remark:"项目级别"`
	Prize string `json:"prize" remark:"奖项等级"`
	Rank  string `json:"rank" remark:"项目排名"`
	Grade string `json:"grade" remark:"成绩"`

	ApplyCredit float64 `json:"apply_credit" remark:"申报学分"`
	ApplyTime   string  `json:"apply_time" remark:"申报时间"`
	AuditMark   string  `json:"audit_mark" remark:"审核时间"`
	AuditCredit float64 `json:"audit_credit" remark:"审核学分"`
}

//获取个人第二课堂申报列表
func (c *SecondClient) GetMySecond() (items []*SecondItem, err error) {
	resp, err := c.doRequest("GET", Urls["second-my"], urlencodedHeader, nil)
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

	items = []*SecondItem{}
	//提取一行表格
	doc.Find("tr").Each(func(i int, selection *goquery.Selection) {
		//匹配html代码内的表格行id，匹配失败说明是表头或者表尾，跳过
		htmlText, _ := selection.Html()
		r := regexp.MustCompile(`id=(\d+)&#`)
		match := r.FindStringSubmatch(htmlText)
		if len(match) < 2 {
			return
		}
		var tmp []string
		//提取单个表格
		selection.Find("td").Each(func(i int, s *goquery.Selection) {
			tmp = append(tmp, s.Text())
		})
		//项目申报页面共有16列，其中第一列和最后一列是不需要提取的
		if len(tmp) != 16 {
			return
		}
		item := &SecondItem{
			ID:         match[1],
			StuID:      "",
			StuName:    "",
			MajorClass: "",
			ApplyYear:  tmp[1],
			Year:       tmp[2],
			Type:       tmp[3],
			Name:       tmp[4],
			Level:      tmp[5],
			Prize:      tmp[6],
			Rank:       tmp[7],
			Grade:      tmp[8],

			AuditMark: tmp[13],
			ApplyTime: tmp[14],
		}
		item.ApplyCredit, _ = strconv.ParseFloat(tmp[11], 64)
		item.AuditCredit, _ = strconv.ParseFloat(tmp[12], 64)

		items = append(items, item)
	})
	return
}

//申报项目公示页面表格
//method GET说明是第一页，POST说明是翻页
func (c *SecondClient) Search(r *http.Request) (items []*SecondItem, err error) {

	if c.VIEWSTATE == "" || c.VIEWSTATEGENERATOR == "" || c.EVENTVALIDATION == "" {
		err = c.updateFormInfo()
		if err != nil {
			logs.Error(err)
			return
		}
	}

	form, err := c.getForm(r)
	if err != nil {
		return nil, err
	}
	requestBody := strings.NewReader(form.Encode())

	resp, err := c.doRequest("POST", Urls["second-search"], urlencodedHeader, requestBody)
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

	c.VIEWSTATE, _ = doc.Find("#__VIEWSTATE").Attr("value")
	c.VIEWSTATEGENERATOR, _ = doc.Find("#__VIEWSTATEGENERATOR").Attr("value")
	c.EVENTVALIDATION, _ = doc.Find("#__EVENTVALIDATION").Attr("value")

	items = []*SecondItem{}
	//记录最大页数
	maxPage := 1
	doc.Find("#MainContent_LabCountPage").Each(func(i int, selection *goquery.Selection) {
		maxPage, _ = strconv.Atoi(selection.Text())
	})
	//请求页数大于最大页数
	page := r.PostForm["page"]
	if len(page) != 0 {
		p, _ := strconv.Atoi(page[0])
		if p > maxPage && maxPage != 0 {
			return
		}
	}

	//提取一行表格
	doc.Find("tr").Each(func(i int, selection *goquery.Selection) {
		//匹配html代码内的表格行id，匹配失败说明是表头或者表尾，跳过
		htmlText, _ := selection.Html()
		r := regexp.MustCompile(`id=(\d+)&#`)
		match := r.FindStringSubmatch(htmlText)
		if len(match) < 2 {
			return
		}
		var tmp []string
		//提取单个表格
		selection.Find("td").Each(func(i int, s *goquery.Selection) {
			tmp = append(tmp, s.Text())
		})
		//查看所有申报项目页面共有18列，其中前两列和最后一列是不需要提取的
		if len(tmp) != 18 {
			return
		}
		item := &SecondItem{
			ID:         match[1],
			StuID:      tmp[2],
			StuName:    tmp[3],
			MajorClass: tmp[4],
			ApplyYear:  tmp[5],
			Year:       tmp[6],
			Type:       tmp[7],
			Name:       tmp[8],
			Level:      tmp[9],
			Prize:      tmp[10],
			Rank:       tmp[11],
			Grade:      tmp[12],

			ApplyTime: tmp[14],
			AuditMark: tmp[15],
		}
		item.ApplyCredit, _ = strconv.ParseFloat(tmp[13], 64)
		item.AuditCredit, _ = strconv.ParseFloat(tmp[16], 64)

		items = append(items, item)
	})

	return
}

//获取申报项目证明材料
func (c *SecondClient) GetImages(itemID string) (images []string, err error) {
	if itemID == "" {
		return
	}
	images = []string{}
	detailUrl := fmt.Sprintf(Urls["second-detail"], itemID)
	resp, err := c.doRequest("GET", detailUrl, urlencodedHeader, nil)
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
	r := regexp.MustCompile(`src="..(.*?)"[\s|\S]*?width`)
	res := r.FindAllStringSubmatch(string(body), -1)
	for _, v := range res {
		if len(v) >= 2 {
			images = append(images, baseUrl+v[1])
		}
	}
	return
}

//发送一次GET请求获取首页的表单信息
func (c *SecondClient) updateFormInfo() (err error) {
	resp, err := c.doRequest("GET", Urls["second-search"], urlencodedHeader, nil)
	if err != nil {
		logs.Error(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	//检查登录状态
	if strings.Contains(string(body), "登录") {
		return AuthError
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logs.Error(err)
		return
	}
	c.VIEWSTATE, _ = doc.Find("#__VIEWSTATE").Attr("value")
	c.VIEWSTATEGENERATOR, _ = doc.Find("#__VIEWSTATEGENERATOR").Attr("value")
	c.EVENTVALIDATION, _ = doc.Find("#__EVENTVALIDATION").Attr("value")
	return
}

func (c *SecondClient) getForm(r *http.Request) (form *url.Values, err error) {

	year := r.PostForm["year"]
	status_no := r.PostForm["status_no"]
	grade := r.PostForm["grade"]
	college_no := r.PostForm["college_no"]
	page := r.PostForm["page"]
	stu_name := r.PostForm["stu_name"]
	item_name := r.PostForm["item_name"]

	if len(year) == 0 {
		year = []string{""}
	}
	if len(status_no) == 0 {
		status_no = []string{""}
	}
	if len(grade) == 0 {
		grade = []string{""}
	}
	if len(college_no) == 0 {
		college_no = []string{""}
	}
	if len(page) == 0 {
		page = []string{"1"}
	}
	if len(stu_name) == 0 {
		stu_name = []string{""}
	}
	if len(item_name) == 0 {
		item_name = []string{""}
	}
	form = &url.Values{
		"__EVENTTARGET":   {""}, //事件来源
		"__EVENTARGUMENT": {""}, //事件参数
		"__LASTFOCUS":     {""},
		//前一个页面提取
		"__VIEWSTATE":          {c.VIEWSTATE},
		"__VIEWSTATEGENERATOR": {c.VIEWSTATEGENERATOR},
		"__EVENTVALIDATION":    {c.EVENTVALIDATION},

		"ctl00$MainContent$xn_list":   year, //*学年
		"ctl00$MainContent$nd_list":   {""},
		"ctl00$MainContent$shbz_list": status_no,  //审核状态
		"ctl00$MainContent$xs_f$nj":   grade,      //年级
		"ctl00$MainContent$xs_f$xy":   college_no, //*学院编号
		"ctl00$MainContent$xs_f$zy":   {""},       //专业编号
		"ctl00$MainContent$xs_f$bj":   {""},       //班级编号
		"ctl00$MainContent$b_xmmc":    item_name,  //*项目名称
		"ctl00$MainContent$b_xm":      stu_name,   //*查询姓名
		"ctl00$MainContent$px1":       {"sbsj"},   //排序类型
		"ctl00$MainContent$px_order":  {"1"},      //升序降序
		"ctl00$MainContent$findbtn":   {"项目申报查询"},
		"ctl00$MainContent$uidField":  {""},
		"ctl00$MainContent$uxmField":  {""},
		"ctl00$MainContent$topage":    page, //目标页数
	}
	//每次变更参数，前端都需要重置页数为1
	if !reflect.DeepEqual(page, []string{"1"}) {
		form.Add("ctl00$MainContent$gopagebtn", "GO")
	}
	return
}
