package gzhu_jw

import (
	"GZHU-Pi/models"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
)

type Rank struct {
	StuID       string  `json:"stu_id"       db:"stu_id" `
	Year        string  `json:"year"         db:"year"`
	YearSem     string  `json:"year_sem"     db:"year_sem"`
	CourseName  string  `json:"course_name"  db:"course_name"`
	CourseID    string  `json:"course_id"    db:"course_id"`
	GradeValue  float64 `json:"grade_value"  db:"grade_value"`
	Gpa         float64 `json:"gpa"          db:"gpa"`
	CollegeRank int64   `json:"college_rank" db:"college_rank"`
	MajorRank   int64   `json:"major_rank"   db:"major_rank"`
	ClassRank   int64   `json:"class_rank"   db:"class_rank"`
}

var (
	//查询大学GPA排名
	gpaTpl = `
            select * from (SELECT stu_id,gpa,college_rank,major_rank,
                         ROW_NUMBER() over (order by gpa desc) as class_rank
                  FROM (SELECT class_id, stu_id, gpa, college_rank, ROW_NUMBER()
                               over (order by gpa desc) as major_rank
                        FROM (SELECT major_id, class_id, stu_id, gpa, ROW_NUMBER()
                                     over (order by gpa desc) as college_rank
                              FROM (select major_id, class_id, stu_id,
                                           cast(sum(credit * course_gpa) / sum(credit) AS decimal(3, 2)) as gpa
                                    from v_grade where stu_id in (select distinct stu_id --学院同级所有学生
                                            from v_grade where college_id = (select college_id from t_stu_info where stu_id = '%s')
                                                       and admit_year = (select admit_year from t_stu_info where stu_id = '%s'))
                                      and invalid != '是' and course_gpa != 0
                                    group by major_id, class_id, stu_id) as a) as b
                        where major_id = (select distinct major_id from t_stu_info where stu_id = '%s')) as c
                  where class_id = (select distinct class_id from t_stu_info where stu_id = '%s')) as d
            where stu_id = '%s';
        `
	//学年GPA排名
	yearTpl = `
            select * from (SELECT year, stu_id, gpa, college_rank, major_rank,
                         ROW_NUMBER() over (PARTITION By year order by gpa desc) as class_rank
                  FROM (SELECT year, class_id, stu_id, gpa, college_rank,
                               ROW_NUMBER() over (PARTITION By year order by gpa desc) as major_rank
                        FROM (SELECT year, major_id, class_id, stu_id, gpa,
                                     ROW_NUMBER() over (PARTITION By year order by gpa desc) as college_rank
                              FROM (select year, major_id, class_id, stu_id,
                                           cast(sum(credit * course_gpa) / sum(credit) AS decimal(3, 2)) as gpa
                                    from v_grade where stu_id in (select distinct stu_id from v_grade
                                                where college_id = (select college_id from t_stu_info where stu_id = '%s')
                                                  and admit_year = (select admit_year from t_stu_info where stu_id = '%s'))
                                        and invalid != '是' and course_gpa != 0
                                    group by year, major_id, class_id, stu_id) as v) as a
                        where major_id = (select major_id from t_stu_info where stu_id = '%s')) as b
                  where class_id = (select class_id from t_stu_info where stu_id = '%s') order by stu_id) as c
            where stu_id = '%s' order by year desc;
        `

	//查询每学期GPA排名
	semTpl = `
            select * from (SELECT year_sem,stu_id,gpa,college_rank,major_rank,
                         ROW_NUMBER() over (PARTITION By year_sem order by gpa desc) as class_rank
                  FROM (SELECT year_sem, class_id, stu_id, gpa, college_rank,
                               ROW_NUMBER() over (PARTITION By year_sem order by gpa desc) as major_rank
                        FROM (SELECT year_sem, major_id, class_id, stu_id, gpa,
                                     ROW_NUMBER() over (PARTITION By year_sem order by gpa desc) as college_rank
                              FROM (select year_sem, major_id, class_id, stu_id,
                                           cast(sum(credit * course_gpa) / sum(credit) AS decimal(3, 2)) as gpa
                                    from v_grade where stu_id in (select distinct stu_id --学院同级所有学生
                                        from v_grade where college_id = (select college_id from t_stu_info where stu_id = '%s')
                                             and admit_year = (select admit_year from t_stu_info where stu_id = '%s'))
                                        and invalid != '是' and course_gpa != 0
									group by year_sem, major_id, class_id, stu_id) as v) as a
                        where major_id = (select major_id from t_stu_info where stu_id = '%s')) as b
                  where class_id = (select class_id from t_stu_info where stu_id = '%s') order by stu_id) as c
            where stu_id = '%s' order by year_sem desc;
        `

	//查询某学号的所有科目成绩及其排名
	courseTpl = `
            select year_sem, course_name, course_id, grade_value, college_rank, major_rank, class_rank
            from (select stu_id, course_name, course_id, grade_value, year_sem, college_rank, major_rank,
                         ROW_NUMBER() over (PARTITION By course_id order by grade_value desc) as class_rank
                  from (select stu_id, class_id, course_name, course_id, grade_value, year_sem, college_rank,
                               ROW_NUMBER() over (PARTITION By course_id order by grade_value desc) as major_rank
                        from (SELECT stu_id, major_id, class_id, course_name, course_id, grade_value, year_sem,
                                     ROW_NUMBER() over (PARTITION By course_id order by grade_value desc) as college_rank
                              FROM v_grade,
                                   (select distinct college_id, admit_year from t_stu_info where stu_id = '%s') as tmp
                              where v_grade.college_id = tmp.college_id
                                and v_grade.admit_year = tmp.admit_year) as a
                        where major_id = (select distinct major_id from t_stu_info where stu_id = '%s')) as b
                  where class_id = (select distinct class_id from t_stu_info where stu_id = '%s')) as c
            where stu_id = '%s'
            order by year_sem desc;
        `
	//人数统计
	stuCountTpl = `
        select (select count(*) from t_stu_info where college_id = a.college_id and admit_year = a.admit_year) as college_count,
        (select count(*) from t_stu_info where major_id = a.major_id and admit_year = a.admit_year)     as major_count,
        (select count(*) from t_stu_info where class_id = a.class_id and admit_year = a.admit_year)     as class_count
        from (select * from t_stu_info where stu_id = '%s') as a;   
        `
)

func (c *JWClient) GetRank(stuID string) (result map[string]interface{}, err error) {

	var (
		gpaSQL    = fmt.Sprintf(gpaTpl, stuID, stuID, stuID, stuID, stuID)
		yearSQL   = fmt.Sprintf(yearTpl, stuID, stuID, stuID, stuID, stuID)
		semSQL    = fmt.Sprintf(semTpl, stuID, stuID, stuID, stuID, stuID)
		courseSQL = fmt.Sprintf(courseTpl, stuID, stuID, stuID, stuID)
		countSQL  = fmt.Sprintf(stuCountTpl, stuID)

		gpaRanks    []Rank
		yearRanks   []Rank
		semRanks    []Rank
		courseRanks []Rank
	)

	result = make(map[string]interface{})

	gpaRanks, err = getGpaRank(gpaSQL)
	if err != nil {
		logs.Error(err)
		return
	}
	yearRanks, err = getGpaRank(yearSQL)
	if err != nil {
		logs.Error(err)
		return
	}
	semRanks, err = getGpaRank(semSQL)
	if err != nil {
		logs.Error(err)
		return
	}
	courseRanks, err = getGpaRank(courseSQL)
	if err != nil {
		logs.Error(err)
		return
	}

	sqlxDB := models.GetSqlx()
	var collegeCount, majorCount, classCount int64
	err = sqlxDB.QueryRowx(countSQL).Scan(&collegeCount, &majorCount, &classCount)
	if err != nil {
		logs.Error(err)
		return
	}
	countMap := make(map[string]int64)
	countMap["college_count"] = collegeCount
	countMap["major_count"] = majorCount
	countMap["class_count"] = classCount

	if len(gpaRanks) > 0 {
		result["gpa_rank"] = gpaRanks[0]
	}
	result["year_rank"] = yearRanks
	result["sem_rank"] = semRanks
	result["course_rank"] = courseRanks
	result["stu_count"] = countMap

	return
}

func getGpaRank(query string) (list []Rank, err error) {
	sqlxDB := models.GetSqlx()

	var rows *sqlx.Rows
	rows, err = sqlxDB.Queryx(query)
	if err != nil {
		logs.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var v Rank
		err = rows.StructScan(&v)
		if err != nil {
			logs.Error(err)
			return
		}
		list = append(list, v)
	}
	return
}
