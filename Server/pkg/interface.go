package pkg

import (
	"GZHU-Pi/pkg/gzhu_jw"
	"net/http"
	"time"
)

type Jwxt interface {
	GetCourse(year, semester string) (courseData *gzhu_jw.CourseData, err error)
	GetExam(year, sem string) (exams []*gzhu_jw.Exam, err error)
	GetAllGrade(year, sem string) (gradeData *gzhu_jw.GradeData, err error)
	GetEmptyRoom(r *http.Request) (data *gzhu_jw.RoomData, err error)
	GetAchieve() (achieves []*gzhu_jw.Achieve, err error)
	SearchAllCourse(xnm, xqm string, page, count int) (data []gzhu_jw.RawCourse, csvData []byte, err error)
	GetExpiresAt() time.Time
	SetExpiresAt(time.Time)
	GetUsername() string
}

type Demo struct {
	//Jwxt
	Username  string
	Password  string
	ExpiresAt time.Time
	Client    *http.Client
}

func (d Demo) GetCourse(year, semester string) (courseData *gzhu_jw.CourseData, err error) {
	panic("implement me")
}

func (d Demo) GetExam(year, sem string) (exams []*gzhu_jw.Exam, err error) {
	panic("implement me")
}

func (d Demo) GetAllGrade(year, sem string) (gradeData *gzhu_jw.GradeData, err error) {
	panic("implement me")
}

func (d Demo) GetEmptyRoom(r *http.Request) (data *gzhu_jw.RoomData, err error) {
	panic("implement me")
}

func (d Demo) GetAchieve() (achieves []*gzhu_jw.Achieve, err error) {
	panic("implement me")
}

func (d Demo) SearchAllCourse(xnm, xqm string, page, count int) (data []gzhu_jw.RawCourse, csvData []byte, err error) {
	panic("implement me")
}

func (d Demo) GetExpiresAt() time.Time {
	panic("implement me")
}

func (d Demo) SetExpiresAt(time.Time) {
	panic("implement me")
}

func (d Demo) GetUsername() string {
	panic("implement me")
}
