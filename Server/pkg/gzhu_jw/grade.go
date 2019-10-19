package gzhu_jw

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type GradeData struct {
	BaseInfo    *BaseInfo   `json:"base_info" remark:"基本信息"`
	GPA         float64     `json:"GPA" remark:"平均绩点"`
	SemList     []*SemGrade `json:"sem_list" remark:"学期列表"`
	TotalCredit float64     `json:"total_credit" remark:"总学分"`
	UpdateTime  string      `json:"update_time" remark:"更新时间"`
}

type BaseInfo struct {
	AdmitYear  string `json:"admit_year" remark:"年级"`
	ClassID    string `json:"class_id" remark:"班级id"`
	College    string `json:"college" remark:"学院"`
	CollegeID  string `json:"college_id" remark:"学院id"`
	Major      string `json:"major" remark:"专业"`
	MajorClass string `json:"major_class" remark:"专业班级"`
	MajorID    string `json:"major_id" remark:"专业id"`
	StuID      string `json:"stu_id" remark:"学号"`
	StuName    string `json:"stu_name" remark:"姓名"`
}

type SemGrade struct {
	GradeList []*Grade `json:"grade_list"  remark:"学期成绩列表"`
	SemCredit float64  `json:"sem_credit" remark:"学期学分"`
	SemGpa    float64  `json:"sem_gpa" remark:"学期绩点"`
	Semester  string   `json:"semester" remark:"学期"`
	Year      string   `json:"year" remark:"学年2018-2019"`
	YearSem   string   `json:"year_sem" remark:"学年学期"`

	GpaCredit float64 `json:"-" remark:"学分*绩点 忽略字段"`
}

type Grade struct {
	CourseGpa  float64 `json:"course_gpa" remark:"课程绩点"`
	CourseID   string  `json:"course_id" remark:"课程ID"`
	CourseName string  `json:"course_name" remark:"课程名称"`
	CourseType string  `json:"course_type" remark:"课程类型"`
	Credit     float64 `json:"credit" remark:"学分"`
	ExamType   string  `json:"exam_type" remark:"考试类型"`
	Grade      string  `json:"grade" remark:"成绩"`
	GradeValue float64 `json:"grade_value" remark:"成绩分数"`
	Invalid    string  `json:"invalid" remark:"是否作废"`
	JxbID      string  `json:"jxb_id" remark:"教学班ID"`
	Semester   string  `json:"semester" remark:"学期"`
	StuID      string  `json:"stu_id" remark:"学号"`
	Teacher    string  `json:"teacher" remark:"教师"`
	Year       string  `json:"year" remark:"学年2018-2019"`
	YearSem    string  `json:"year_sem" remark:"学年学期"`
}

func (c *JWClient) GetAllGrade(year, sem string) (gradeData *GradeData, err error) {

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
	//检查登录状态
	if strings.Contains(string(body), "登录") {
		return nil, AuthError
	}

	gradeData = &GradeData{}
	//提取成绩列表及基本信息
	grades, baseInfo := ParseGrade(body)
	gradeData.BaseInfo = baseInfo

	//根基成绩列表统计所有成绩信息，传址
	CountGpa(grades, gradeData)

	gradeData.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

	return gradeData, nil
}

//统计GPA信息，指针传递
func CountGpa(grades []*Grade, gradeData *GradeData) {

	var (
		sumCredit    float64 = 0                         //大学总学分绩点
		sumGpaCredit float64 = 0                         //大学总学分*绩点
		semData              = make(map[string]SemGrade) //各学期成绩数据
	)

	for k, v := range grades {

		if v.Invalid == "是" || v.CourseGpa == 0 || v.Credit == 0 {
			logs.Debug("作废或者不及格成绩，跳过统计", k, v)
			semData[v.YearSem] = SemGrade{
				GradeList: append(semData[v.YearSem].GradeList, v),
				Semester:  v.Semester,
				Year:      v.Year,
				YearSem:   v.YearSem,
			}
			continue
		}
		//累计学分绩点
		sumCredit = sumCredit + v.Credit
		sumGpaCredit = sumGpaCredit + v.Credit*v.CourseGpa

		//计算各个学期的学分绩点
		semData[v.YearSem] = SemGrade{
			SemCredit: semData[v.YearSem].SemCredit + v.Credit,
			GpaCredit: semData[v.YearSem].GpaCredit + v.Credit*v.CourseGpa,

			GradeList: append(semData[v.YearSem].GradeList, v),
			Semester:  v.Semester,
			Year:      v.Year,
			YearSem:   v.YearSem,
		}
		tmp := semData[v.YearSem]
		if semData[v.YearSem].SemCredit != 0 {
			tmp.SemGpa, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", semData[v.YearSem].GpaCredit/semData[v.YearSem].SemCredit), 2)
		}
		semData[v.YearSem] = tmp
	}
	//大学总学分绩点
	gradeData.TotalCredit = sumCredit
	if sumCredit != 0 {
		gradeData.GPA, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", sumGpaCredit/sumCredit), 2)
	}
	//把map转换成切片
	var tmpList []*SemGrade
	for _, v := range semData {
		gradeList := v
		tmpList = append(tmpList, &gradeList)
	}
	//按学期倒序排列，新学期在前
	for i := 0; i < len(tmpList); i++ {
		for j := i + 1; j < len(tmpList); j++ {
			if bytes.Compare([]byte(tmpList[i].YearSem), []byte(tmpList[j].YearSem)) == -1 {
				tmpList[i], tmpList[j] = tmpList[j], tmpList[i]
			}
		}
	}
	gradeData.SemList = append(gradeData.SemList, tmpList...)
}

//解析提取成绩信息，同时填充学生基本信息
func ParseGrade(body []byte) (grades []*Grade, info *BaseInfo) {

	grades = []*Grade{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	gradeList := json.Get(body, "items")
	//遍历所有事件课程
	for i := 0; true; i++ {
		g := &Grade{}
		g.CourseID = gradeList.Get(i).Get("kch_id").ToString()
		if g.CourseID == "" {
			break
		}
		if i == 0 {
			info = &BaseInfo{
				AdmitYear:  gradeList.Get(i).Get("njdm_id").ToString(),
				ClassID:    gradeList.Get(i).Get("bh_id").ToString(),
				College:    gradeList.Get(i).Get("jgmc").ToString(),
				CollegeID:  gradeList.Get(i).Get("jg_id").ToString(),
				Major:      gradeList.Get(i).Get("zymc").ToString(),
				MajorClass: gradeList.Get(i).Get("bj").ToString(),
				MajorID:    gradeList.Get(i).Get("zyh_id").ToString(),
				StuID:      gradeList.Get(i).Get("xh").ToString(),
				StuName:    gradeList.Get(i).Get("xm").ToString(),
			}
		}
		g.CourseGpa = gradeList.Get(i).Get("jd").ToFloat64()
		g.CourseName = gradeList.Get(i).Get("kcmc").ToString()
		g.CourseType = gradeList.Get(i).Get("kcxzmc").ToString()
		g.Credit = gradeList.Get(i).Get("xf").ToFloat64()
		g.ExamType = gradeList.Get(i).Get("ksxz").ToString()
		g.Grade = gradeList.Get(i).Get("cj").ToString()
		g.GradeValue = gradeList.Get(i).Get("bfzcj").ToFloat64()
		g.Invalid = gradeList.Get(i).Get("cjsfzf").ToString()
		g.JxbID = gradeList.Get(i).Get("jxb_id").ToString()
		g.Semester = gradeList.Get(i).Get("xqmmc").ToString()
		g.StuID = gradeList.Get(i).Get("xh").ToString()
		g.Teacher = gradeList.Get(i).Get("jsxm").ToString()
		g.Year = gradeList.Get(i).Get("xnmmc").ToString()
		g.YearSem = g.Year + "-" + g.Semester

		grades = append(grades, g)
	}
	return
}
