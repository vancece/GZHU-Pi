package gzhu_jw

import (
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Grade struct {
	CourseGpa  float64 `json:"course_gpa" remark:"课程绩点"`
	CourseID   string  `json:"course_id" remark:"课程ID"`
	CourseName string  `json:"course_name" remark:"课程名称"`
	CourseType string  `json:"course_type" remark:"课程类型"`
	Credit     int     `json:"credit" remark:"学分"`
	ExamType   string  `json:"exam_type" remark:"考试类型"`
	Grade      string  `json:"grade" remark:"成绩"`
	GradeValue int     `json:"grade_value" remark:"成绩分数"`
	Invalid    string  `json:"invalid" remark:"是否作废"`
	JxbID      string  `json:"jxb_id" remark:"教学班ID"`
	Semester   string  `json:"semester" remark:"学期"`
	StuID      string  `json:"stu_id" remark:"学号"`
	Teacher    string  `json:"teacher" remark:"教师"`
	Year       string  `json:"year" remark:"学年2018-2019"`
	YearSem    string  `json:"year_sem" remark:"学年学期"`
}

func (c *JWClient) GetAllGrade(year, sem string) (grades []*Grade, err error) {

	nd := time.Now().Unix() * 1000 //时间戳
	var form = url.Values{
		"xh_id":                  {c.Username},
		"xnm":                    {year}, //year sem为空时获取所有成绩
		"xqm":                    {sem},
		"_search":                {"false"},
		"nd":                     {strconv.Itoa(int(nd))},
		"queryModel.showCount":   {"200"},
		"queryModel.currentPage": {"1"},
		"queryModel.sortName":    {""},
		"queryModel.sortOrder":   {"asc"},
		"time":                   {"0"},
	}

	resp, err := c.doRequest("POST", Urls["grade"], urlencodedHeader, strings.NewReader(form.Encode()))
	if err != nil {
		logs.Error(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	grades = ParseGrade(body)
	return grades, nil
}

//解析提取成绩信息
func ParseGrade(body []byte) (grades []*Grade) {

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	gradeList := json.Get(body, "items")
	//遍历所有事件课程
	for i := 0; true; i++ {
		g := &Grade{}
		g.CourseID = gradeList.Get(i).Get("kch_id").ToString()
		if g.CourseID == "" {
			break
		}
		g.CourseGpa = gradeList.Get(i).Get("jd").ToFloat64()
		g.CourseName = gradeList.Get(i).Get("kcmc").ToString()
		g.CourseType = gradeList.Get(i).Get("kcxzmc").ToString()
		g.Credit = gradeList.Get(i).Get("xf").ToInt()
		g.ExamType = gradeList.Get(i).Get("ksxz").ToString()
		g.Grade = gradeList.Get(i).Get("cj").ToString()
		g.GradeValue = gradeList.Get(i).Get("bfzcj").ToInt()
		g.Invalid = gradeList.Get(i).Get("cjsfzf").ToString()
		g.JxbID = gradeList.Get(i).Get("jxb_id").ToString()
		g.Semester = gradeList.Get(i).Get("xqmmc").ToString()
		g.StuID = gradeList.Get(i).Get("xh").ToString()
		g.Teacher = gradeList.Get(i).Get("jsxm").ToString()
		g.Year = gradeList.Get(i).Get("xnmmc").ToString()
		//g.YearSem = gradeList.Get(i).Get("").ToString()

		grades = append(grades, g)
	}
	return
}
