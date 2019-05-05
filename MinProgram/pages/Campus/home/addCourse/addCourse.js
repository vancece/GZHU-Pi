const Page = require('../../../../utils/sdk/ald-stat.js').Page;
Page({

  data: {
    colorValue: 4,
    weekIndex: [0, 15],
    class_timeIndex: [6, 0, 1],
    tempIndex: 0,
    tempIndex0: 0,

    colorArrays: ["#86b0fe", "#71eb55", "#f7c156", "#76e9eb", "#ff9dd8", "#80f8e6", "#eaa5f7", "#86b3a5", "#85B8CF", "#90C652"],
    multiWeeks: [
      [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20],
      ["至1周", "至2周", "至3周", "至4周", "至5周", "至6周", "至7周", "至8周", "至9周", "至10周", "至11周", "至12周", "至13周", "至14周", "至15周", "至16周", "至17周", "至18周", "至19周", "至20周", ]
    ],
    multiTime: [
      ["周一", "周二", "周三", "周四", "周五", "周六", "周日"],
      [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11],
      ["至1节", "至2节", "至3节", "至4节", "至5节", "至6节", "至7节", "至8节", "至9节", "至10节", "至11节"]
    ]
  },


  onLoad: function(options) {
    let that = this
    wx.getStorage({
      key: 'course',
      success: function(res) {
        that.setData({
          course: res.data
        })
      },
      fail: function(res) {
        let course = {
          "course_list": []
        }
        that.setData({
          course: course
        })
      }
    })
  },

  // 颜色单选器
  radioChange(e) {
    this.setData({
      colorValue: Number(e.detail.value)
    })
  },

  // 周选择器改变
  multiWeeksColumnChange(e) {
    let index = e.detail.value
    // 第一列
    if (e.detail.column == 0) {
      this.setData({
        weekIndex: [index, index],
        tempIndex: index
      })
    }
    // 第二列
    else {
      if (this.data.tempIndex > index) {
        this.setData({
          weekIndex: [index, index],
          tempIndex: index
        })
      } else {
        this.setData({
          weekIndex: [this.data.tempIndex, index]
        })
      }
    }
  },

  // 节次选择器改变
  multiTimeColumnChange(e) {
    console.log(e)
    let index = e.detail.value

    if (e.detail.column == 0) {
      this.data.tempIndex0 = e.detail.value
      this.setData({
        class_timeIndex: [this.data.tempIndex0, this.data.class_timeIndex[1], this.data.class_timeIndex[2]],
      })
    }
    if (e.detail.column == 1) {
      this.setData({
        class_timeIndex: [this.data.tempIndex0, index, index],
        tempIndex: index
      })
    } else if (e.detail.column == 2) {
      if (this.data.tempIndex > index) {
        this.setData({
          class_timeIndex: [this.data.tempIndex0, index, index],
          tempIndex: index
        })
      } else {
        this.setData({
          class_timeIndex: [this.data.tempIndex0, this.data.tempIndex, index]
        })
      }
    }
  },

  formSubmit(e) {
    let obj = e.detail.value
    if (obj.course_name == "" || obj.class_place == "" || obj.teacher == "") {
      wx.showToast({
        title: '未填写完整',
        icon: "none",
        duration: 2000
      })
    } else {

      let weekday = obj.class_time[0] + 1 //星期代码
      let which_day = this.data.multiTime[0][obj.class_time[0]] //星期几
      let class_start = obj.class_time[1] + 1 //开始节次
      let class_last = obj.class_time[2] - obj.class_time[1] + 1 //持续节数
      let class_time = (obj.class_time[1]) + 1 + "-" + (obj.class_time[2] + 1) + "节" //节次信息
      let weeks = (obj.weeks[0] + 1) + "-" + (obj.weeks[1] + 1) + "周" //周次信息
      let course_name = obj.course_name + "*" //课程名称

      let courseItem = {
        "check_type": obj.check_type,
        "last": class_last,
        "class_place": obj.class_place,
        "start": class_start,
        "course_time": class_time,
        "color": Number(obj.color),
        "course_id": "add",
        "course_name": course_name,
        "teacher": obj.teacher,
        "weekday": weekday,
        "weeks": weeks,
        "which_day": which_day
      }
      // 加入课程并缓存
      let course = this.data.course
      course.course_list.push(courseItem)
      console.log("545",course)
      wx.setStorage({
        key: 'course',
        data: course,
        success: function() {
         wx.reLaunch({
           url: "/pages/Campus/home/home"
         })
        }
      })
    }
  },


})