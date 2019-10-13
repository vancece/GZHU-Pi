package gzhu_jw

import (
	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Exam struct {
	CourseID   string `json:"course_id" remark:"课程ID"`
	ExamCourse string `json:"exam_course" remark:"考试科目"`
	ExamRoom   string `json:"exam_room" remark:"考试地点"`
	ExamTime   string `json:"exam_time" remark:"考试时间"`
	Major      string `json:"major" remark:"专业"`
	MajorClass string `json:"major_class" remark:"专业班级"`
	Sem        string `json:"sem" remark:"学期"`
	Credit     string `json:"credit" remark:"学分"`
	Year       string `json:"year" remark:"学年2018-2019"`
}

func (c *JWClient) GetExam(year, sem string) (exams []*Exam, err error) {

	nd := time.Now().Unix() * 1000 //时间戳
	var form = url.Values{
		"xnm":                    {year},
		"xqm":                    {sem},
		"ksmcdmb_id":             {""},
		"kch":                    {""},
		"kc":                     {""},
		"ksrq":                   {""},
		"_search":                {"false"},
		"nd":                     {strconv.Itoa(int(nd))},
		"queryModel.showCount":   {"15"},
		"queryModel.currentPage": {"1"},
		"queryModel.sortName":    {""},
		"queryModel.sortOrder":   {"asc"},
		"time":                   {"0"},
	}

	resp, err := c.doRequest("POST", Urls["exam"], urlencodedHeader, strings.NewReader(form.Encode()))
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

	exams = ParseExam(body)

	return
}

//解析提取考试信息
func ParseExam(body []byte) (exams []*Exam) {

	exams = []*Exam{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	examList := json.Get(body, "items")
	//遍历所有事件课程
	for i := 0; true; i++ {
		e := &Exam{}
		e.CourseID = examList.Get(i).Get("kch").ToString()
		if e.CourseID == "" {
			break
		}
		e.ExamCourse = examList.Get(i).Get("kcmc").ToString()
		e.ExamRoom = examList.Get(i).Get("cdmc").ToString()
		e.ExamTime = examList.Get(i).Get("kssj").ToString()
		e.Major = examList.Get(i).Get("zymc").ToString()
		e.MajorClass = examList.Get(i).Get("bj").ToString()
		e.Sem = examList.Get(i).Get("xqmmc").ToString()
		e.Credit = examList.Get(i).Get("xf").ToString()
		e.Year = examList.Get(i).Get("xnmc").ToString()

		exams = append(exams, e)
	}
	return
}
