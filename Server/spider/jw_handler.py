from lxml import html
from jsonpath_rw import parse
import json
import re
import time

"""

教务系统相关数据处理函数集合

"""


def get_login_form(text, username, password):
    """
    获取登录表单
    :param text: 登录页面html文本
    :param username: 用户名
    :param password: 密码
    :return: 用于POST的登录表单
    """

    selector = html.fromstring(text)  # 将html文件转换为xpath可以识别的结构
    target = selector.xpath('//div[@class="row btn-row"]/input/@value')
    lt = target[0]
    execution = target[1]

    form_data = {
        "username": username,
        "password": password,
        "captcha": "",
        "warn": "true",
        "lt": lt,
        "execution": execution,
        "_eventId": "submit",
        "submit": "登录"
    }
    return form_data


def get_student_info(text):
    """
    获取学生信息
    :param text: 学生信息页面html文本
    :return: 爬取的学生基础信息
    """

    selector = html.fromstring(text)  # 将html文件转换为xpath可以识别的结构
    name = selector.xpath('//*[@id="ajaxForm"]/div/div[1]/div/div[2]/div/div/p')[0].text
    student_id = selector.xpath('//*[@id="ajaxForm"]/div/div[1]/div/div[1]/div/div/p')[0].text
    year = selector.xpath('//*[@id="col_njdm_id"]/p')[0].text
    college = selector.xpath('//*[@id="col_jg_id"]/p')[0].text
    major = selector.xpath('//*[@id="col_zyh_id"]/p')[0].text
    major_class = selector.xpath('//*[@id="col_bh_id"]/p')[0].text

    student_info = {
        "name": name,
        "student_id": student_id,
        "year": year.replace("\t", "").replace("\n", "").replace("\r", ""),
        "college": college.replace("\t", "").replace("\n", "").replace("\r", ""),
        "major": major.replace("\t", "").replace("\n", "").replace("\r", ""),
        "major_class": major_class.replace("\t", "").replace("\n", "").replace("\r", "")
    }
    return student_info


def get_course(text):
    """
    获取课表信息
    :param text: 获取的课表JSON文本
    :return: 处理过的课表数据 dict
    """

    kb_json = json.loads(text)  # 转换成json,dict类型

    # 用jsonpath选取课程信息，类型为list
    course_id = parse('$.kbList[*].kch_id').find(kb_json)  # 课程ID
    course_name = parse('$.kbList[*].kcmc').find(kb_json)  # 课程名称
    class_place = parse('$.kbList[*].cdmc').find(kb_json)  # 上课地点
    which_day = parse('$.kbList[*].xqjmc').find(kb_json)  # 星期几
    course_time = parse('$.kbList[*].jc').find(kb_json)  # 上课时间（节数）
    weeks = parse('$.kbList[*].zcd').find(kb_json)  # 周数
    teacher = parse('$.kbList[*].xm').find(kb_json)  # 教师姓名
    check_type = parse('$.kbList[*].khfsmc').find(kb_json)  # 考核类型
    # 实践课程，课表底部
    sjk_course_name = parse('$.sjkList[*].kcmc').find(kb_json)  # 课程名称
    sjk_weeks = parse('$.sjkList[*].qsjsz').find(kb_json)  # 周数
    sjk_teacher = parse('$.sjkList[*].xm').find(kb_json)  # 教师姓名

    course_list = []
    for i, item in enumerate(course_id):
        course = {}
        course["course_id"] = course_id[i].value
        course["course_name"] = course_name[i].value
        course["class_place"] = class_place[i].value
        course["which_day"] = which_day[i].value
        course["course_time"] = course_time[i].value
        course["weeks"] = weeks[i].value
        course["teacher"] = teacher[i].value
        course["check_type"] = check_type[i].value
        course_list.append(course)

    sjk_course_list = []
    for i, item in enumerate(sjk_course_name):
        course = {}
        course["sjk_course_name"] = sjk_course_name[i].value
        course["sjk_weeks"] = sjk_weeks[i].value
        course["sjk_teacher"] = sjk_teacher[i].value
        sjk_course_list.append(course)

    courses = {"course_list": course_list, "sjk_course_list": sjk_course_list}
    courses = handle_course(courses)
    return courses


def add_credit(text, courses):
    """
    给对应课程添加学分信息,课表页面获取的课程不包含学分信息
    :param text: 个人信息--选课页面html文本
    :param courses: 处理过的课程数据
    :return:增加了学分的课程数据
    """

    kb_json = json.loads(text)
    course_id = parse('$.items[*].kch').find(kb_json)  # 课程ID
    credit = parse('$.items[*].xf').find(kb_json)  # 课程名称
    print(len(credit))
    for i, item in enumerate(course_id):
        for i2, item2 in enumerate(courses["course_list"]):
            if item.value == item2["course_id"]:
                item2["credit"] = credit[i].value
    return courses


def get_grade(text):
    """
    获取成绩信息
    :param text: 获取的成绩JSON文本
    :return: 筛选处理的成绩数据
    """

    grade_json = json.loads(text)
    # 筛选数据
    year = parse('$.items[*].xnmmc').find(grade_json)  # 学年 2017~2018
    semester = parse('$.items[*].xqmmc').find(grade_json)  # 学期 1/2
    course_id = parse('$.items[*].kch_id').find(grade_json)  # 课程代码
    course_name = parse('$.items[*].kcmc').find(grade_json)  # 课程名称
    credit = parse('$.items[*].xf').find(grade_json)  # 学分
    grade_value = parse('$.items[*].bfzcj').find(grade_json)  # 成绩分数
    grade = parse('$.items[*].cj').find(grade_json)  # 成绩
    course_gpa = parse('$.items[*].jd').find(grade_json)  # 绩点
    course_type = parse('$.items[*].kcxzmc').find(grade_json)  # 课程性质
    exam_type = parse('$.items[*].ksxz').find(grade_json)  # 考试性质 正常/补考/重修
    invalid = parse('$.items[*].cjsfzf').find(grade_json)  # 成绩是否作废

    total_count = parse('$.totalCount').find(grade_json)[0].value  # 成绩总条数

    grade_list = []
    for i, item in enumerate(course_id):
        temp = {}
        temp["year"] = year[i].value
        temp["semester"] = semester[i].value
        temp["course_id"] = course_id[i].value
        temp["course_name"] = course_name[i].value
        temp["credit"] = credit[i].value
        temp["grade_value"] = grade_value[i].value
        temp["grade"] = grade[i].value
        temp["course_gpa"] = course_gpa[i].value
        temp["course_type"] = course_type[i].value
        temp["exam_type"] = exam_type[i].value
        temp["invalid"] = invalid[i].value
        grade_list.append(temp)

    grade = handle_grade(grade_list, total_count)
    return grade


def get_exam(text):
    """
    获取考试信息
    :param text: 获取的考试JSON文本
    :return: 筛选处理的考试数据
    """

    exam_json = json.loads(text)  # 转换成json,dict类型
    # 筛选数据
    exam_course = parse('$.items[*].kcmc').find(exam_json)  # 课程名称
    exam_time = parse('$.items[*].kssj').find(exam_json)  # 考试时间
    exam_room = parse('$.items[*].cdmc').find(exam_json)  # 考试地点
    exam_major = parse('$.items[*].bj').find(exam_json)  # 专业班级
    exam_class = parse('$.items[*].jxbzc').find(exam_json)  # 教学班级
    year = parse('$.items[*].xnmc').find(exam_json)  # 学年
    sem = parse('$.items[*].xqmmc').find(exam_json)  # 学期
    credit = parse('$.items[*].xf').find(exam_json)  # 学分

    exam_list = []
    for idx, item in enumerate(exam_course):
        temp = {}
        temp["exam_course"] = exam_course[idx].value
        temp["exam_time"] = exam_time[idx].value
        temp["exam_room"] = exam_room[idx].value
        temp["exam_major"] = exam_major[idx].value
        temp["exam_class"] = exam_class[idx].value
        temp["year"] = year[idx].value
        temp["sem"] = sem[idx].value
        temp["xf"] = credit[idx].value
        exam_list.append(temp)
    return exam_list


"""
辅助数据处理函数
"""


def handle_course(courses):
    """
    处理课表数据，适配前端
    :param courses: 基础课表数据
    :return: 增加星期、开始持续节数的课表数据
    """

    set_list = set(())  # 定义空集合，记录所有不同的课程
    for item in courses['course_list']:
        set_list.add(item["course_id"])  # 生成id集合

        class_time = item['course_time']

        reg = "\d+"
        class_res = re.findall(reg, class_time)

        # 生成开始节和持续节数
        if len(class_res) == 2:
            item['start'] = int(class_res[0])
            item['last'] = int(class_res[1]) - int(class_res[0]) + 1
        else:
            item['start'] = int(class_res[0])
            item['last'] = 1

        # 转换星期几至数字
        switcher = {
            "星期一": 1,
            "星期二": 2,
            "星期三": 3,
            "星期四": 4,
            "星期五": 5,
            "星期六": 6,
            "星期日": 7,
            "星期天": 7,
        }
        item['weekday'] = switcher.get(item['which_day'], "未安排")

    # 给每种不同的课程标号，相同课程标号相同
    for item1 in courses['course_list']:
        for item2 in set_list:
            if item1["course_id"] == item2:
                item1["color"] = list(set_list).index(item2)

    return courses


def handle_grade(grade_list, total_count):
    """
    处理成绩信息，计算学分绩点，分学期学年整理
    :param grade_list: 初步处理的成绩数据
    :param total_count: 成绩数据总条数
    :return:
    """

    if total_count == 0:
        grade = {"update_time": time.strftime("%Y-%m-%d %H:%M:%S"), "sem_list": [], "totalCount": total_count}
        return grade

    jd_xf, xf = 0, 0  # 绩点 x 学分，学分

    list_year = []  # 定义空列表，有序记录所有不同的学年
    list_sem = []  # 有序记录所有不同的 学年--学期

    for item in grade_list:
        if item["year"] not in list_year:
            list_year.append(item["year"])
        # 去除不及格和作废成绩
        if item["course_gpa"] != "0.00" and item["invalid"] == "否":
            xf = xf + float(item["credit"])  # 总学分，分母
            jd_xf = jd_xf + float(item["course_gpa"]) * float(item["credit"])

    if xf == 0:
        GPA = round(0, 2)
    else:
        GPA = round(jd_xf / xf, 2)  # 大学总绩点
    grade = {"GPA": GPA, "total_credit": xf, "totalCount": total_count,
             "update_time": time.strftime("%Y-%m-%d %H:%M:%S")}

    # 添加 学年-学期  如2017-2018-2
    for set_item in list_year:
        for item in grade_list:
            if item["year"] == set_item:
                if item["semester"] == "1":
                    item["year_sem"] = item["year"] + "-1"
                    if item["year_sem"] not in list_sem:
                        list_sem.append(item["year_sem"])
                else:
                    item["year_sem"] = item["year"] + "-2"
                    if item["year_sem"] not in list_sem:
                        list_sem.append(item["year_sem"])

    temp_sem_list = []  # 所有学期的成绩存放于一个列表
    for set_item in list_sem:
        jd_xf, xf = 0, 0

        temp_sem = {}  # 每个学期的成绩存放于一个字典
        for item in grade_list:
            if item["year_sem"] == set_item:
                temp_sem["year_sem"] = item["year_sem"]
                temp_sem["year"] = item["year"]
                temp_sem["semester"] = item["semester"]
                if item["course_gpa"] != "0.00" and item["invalid"] == "否":
                    xf = xf + float(item["credit"])  # 总学分，分母
                    jd_xf = jd_xf + float(item["course_gpa"]) * float(item["credit"])
        if xf == 0:
            sem_gpa = round(0, 2)
        else:
            sem_gpa = round(jd_xf / xf, 2)
        temp_sem["sem_credit"] = xf  # 学期总学分
        temp_sem["sem_gpa"] = sem_gpa  # 学期绩点
        temp_sem_list.append(temp_sem)

    sum_arr = []
    for sem_item in temp_sem_list:
        temp = []
        for item in grade_list:
            if item["year_sem"] == sem_item["year_sem"]:
                temp.append(item)
            sem_item["grade_list"] = temp

        tmp_arr = sem_item["year_sem"].split("-")  # 拆分2017-2018-1
        tmp_arr = [int(item) for item in tmp_arr]  # 转换int求和
        sum = 0
        for i in tmp_arr:
            sum += i
        sum_arr.append(sum)

    sum_arr.sort(reverse=True)  # 求和降序
    # 按学年学期排序
    sem_list = []
    for x in sum_arr:
        for sem_item in temp_sem_list:
            tmp_arr = sem_item["year_sem"].split("-")  # 拆分2017-2018-1
            tmp_arr = [int(item) for item in tmp_arr]  # 转换int求和
            sum = 0
            for i in tmp_arr:
                sum += i
            if x == sum:
                sem_list.append(sem_item)

    grade["sem_list"] = sem_list

    return grade
