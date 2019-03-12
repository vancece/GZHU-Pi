var utils = require("../../../../../utils/utils.js")
var publicData = require("../../../../../utils/public_data.js")
Component({

  properties: {
    show: {
      type: Boolean,
      value: false
    },
  },

  data: {
    hideTimeLine: true,
    showDetail: false,
    current: 1,
    dis: "none",
    today: new Date().getDate(), //日期
    week: utils.getSchoolWeek(), //周数
    schoolWeek: utils.getSchoolWeek(), //校历周
    weekDate: utils.setWeekDate(), //一周日期

    weekDays: publicData.weekDays,
    timeLine: publicData.timeLine,
    colors: publicData.colors,
    kbList: publicData.course_sample
  },

  methods: {
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
      let detail = []
      // 遍历课表，找出星期和开始节相同的课程
      this.data.kbList.forEach(function(item) {
        if (item.weekday == day && item.start == start) {
          detail.push(item)
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
          this.deleteCourse()
          break
      }
    },

    // 删除课程
    deleteCourse() {
      let that = this
      let id = this.data.openTarget
      let obj = wx.getStorageSync('course')
      obj.course_list.splice(id, 1)
      wx.showModal({
        title: '提醒',
        content: '是否删除当前课程?',
        success: function(e) {
          if (e.confirm) {
            wx.setStorage({
              key: 'course',
              data: obj,
              success: function() {
                that.setData({
                  kbList: obj.course_list,
                  showDetail: false,
                })
              }
            })
          }
        }
      })
    },
  },

  lifetimes: {
    created: function() {},

    attached: function() {
      if (wx.getStorageSync('course') != "") {
        this.setData({
          kbList: wx.getStorageSync('course').course_list,
          sjkList: wx.getStorageSync('course').sjk_course_list
        })
      }
    },

    ready: function() {

    }
  },


})