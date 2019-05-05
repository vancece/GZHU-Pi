const Page = require('../../../utils/sdk/ald-stat.js').Page;
var utils = require("../../../utils/utils.js")
var Data = require("../../../utils/data.js")
var Config = require("../../../utils/config.js")
Page({

  data: {
    hideTimeLine: true,
    showDetail: false,
    current: 1,
    dis: "none",
    today: new Date().getDate(), //日期
    week: utils.getSchoolWeek(), //周数
    schoolWeek: utils.getSchoolWeek(), //校历周
    weekDate: utils.setWeekDate(), //一周日期
    bg: Config.get("schedule_bg"), // 获取背景
    blur: Config.get("blur"), //高斯模糊

    weekDays: Data.weekDays,
    timeLine: Data.timeLine,
    colors: Data.colors,
    kbList: Data.course_sample
  },


  // 恢复校历周
  resetWeek() {
    let week = utils.getSchoolWeek()
    this.setData({
      week: week,
      weekDate: utils.setWeekDate(),
    })
    wx.showToast({
      title: "校历 " + String(week) + " 周",
      icon: "none",
      duration: 1000
    })
  },
  // 左右滑动切换周数
  switchWeek(e) {
    let value = e.detail.current - this.data.current
    let week
    if (value == 1 || value == -2) {
      // 下一周
      if (this.data.week + 1 > 20) {
        week = 0
      } else if (this.data.week < 0) {
        week = 1
      } else {
        week = this.data.week + 1
      }
    } else {
      // 上一周
      if (this.data.week - 1 < 0) {
        week = 20
      } else {
        week = this.data.week - 1
      }
    }
    this.setData({
      weekDate: utils.setWeekDate(week - this.data.schoolWeek),
      week: week,
      current: e.detail.current,
    })
    wx.showToast({
      title: "第 " + String(week) + " 周",
      icon: "none",
      duration: 1000
    })
  },
  // 展开时间轴
  tapSlideBar() {
    this.setData({
      hideTimeLine: !this.data.hideTimeLine,
    })
  },

  // 课程详情弹窗
  showDetail(e) {
    let that = this
    let id = Number(e.currentTarget.id)
    let day = this.data.kbList[id].weekday
    let start = this.data.kbList[id].start
    let detail = [this.data.kbList[id]]
    // 遍历课表，找出星期和开始节相同的课程
    this.data.kbList.forEach(function(item) {
      if (item.weekday == day && item.start == start) {
        if (that.data.kbList.indexOf(item) != id) detail.push(item)
      }
    })
    this.setData({
      detail: detail,
      showDetail: true,
      currentIndex: 0 //恢复滑动视图索引
    })
    this.showCourseId(0)
  },
  // 左右滑动切换课程
  switchCourse(e) {
    this.showCourseId(e.detail.current)
  },
  // 打开或者切换时更新显示的课程数组索引
  showCourseId(current) {
    let course = this.data.detail[current]
    for (let i = 0; i < this.data.kbList.length; i++) {
      if (course == this.data.kbList[i]) {
        this.data.openTarget = i
      }
    }
  },

  // 关闭课程详情弹窗
  cancelModal() {
    this.setData({
      showDetail: false,
    })
  },

  // 编辑课表
  navTo(e) {
    switch (e.currentTarget.id) {
      case "0": //编辑
        wx.navigateTo({
          url: '/pages/Campus/home/addCourse/addCourse',
        })
        break
      case "1": //添加
        wx.navigateTo({
          url: '/pages/Campus/evaluation/evaluation',
        })
        break
      case "2": //删除
        break
    }
  },

  onLoad: function(options) {

    wx.showLoading({
      title: 'Loading...',
    })
    this.data.className = options["classNmae"]
    this.data.year = options["year"]
    this.data.sem = options["sem"]
    this.getClassSchedule(options["classNmae"], options["year"], options["sem"])
  },

  onShareAppMessage: function() {
    return {
      title: this.data.className + " " + this.data.year + "-" + this.data.sem,
      desc: this.data.className + '的课表',
      path: '/pages/Campus/course/schedule?classNmae=' + this.data.className + "&year=" + this.data.year + "&sem=" + this.data.sem,
    }
  },

  // 获取班级课表并处理
  getClassSchedule(className = "软件181", year = "2018-2019", sem = "2") {
    let that = this
    let Obj = new wx.BaaS.TableObject("all_course")
    let query = new wx.BaaS.Query()
    query.compare('xn', '=', year) //学年
    query.compare('xq', '=', sem) //学期
    query.contains('jxbzc', className) //班级

    Obj.setQuery(query).find().then(res => {
      let kbList = that.handleSchedule(res.data.objects)
      that.setData({
        className: className,
        kbList: kbList,
        sjkList: []
      })
      wx.hideLoading()
    }, err => {
      wx.hideLoading()
    })
  },


  // 处理课程数据，生成完整课表
  handleSchedule(rawSchedule) {
    let set_arr = []
    let classSchedule = []

    rawSchedule.forEach(item => {
      if (set_arr.indexOf(item.kch_id) == -1) set_arr.push(item.kch_id)
      let newItem = {
        "credit": item.xf,
        "class_place": item.cdmc,
        "course_id": item.kch_id,
        "course_name": item.kcmc,
        "jgh_id": item.jgh_id,
        "teacher": item.xm,
        "course_time": item.skjc,
        "weeks": item.qsjsz,
        "weekday": item.xqj,
        "color": set_arr.indexOf(item.kch_id)
      }
      // 合并对象
      Object.assign(newItem, this.handleCourseTime(item.skjc, item.xqj))
      classSchedule.push(newItem)
    })
    return classSchedule
  },

  // 提取上课时间
  handleCourseTime(course_time, weekday) {

    let result = {}
    // 提取节次
    let sections = course_time.match(/\d+/g)
    result['start'] = Number(sections[0])
    result['last'] = Number(sections[1]) - Number(sections[0]) + 1
    // 转换星期
    let switcher = {
      "1": "星期一",
      "2": "星期二",
      "3": "星期三",
      "4": "星期四",
      "5": "星期五",
      "6": "星期六",
      "7": "星期日"
    }
    result['which_day'] = switcher[weekday]
    return result
  },
})