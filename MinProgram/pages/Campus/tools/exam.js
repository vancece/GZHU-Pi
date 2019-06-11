const Page = require('../../../utils/sdk/ald-stat.js').Page;
var app = getApp()
var tableID = 57516
Page({
  data: {

  },

  onLoad: function (options) {
    console.log("page启动：", options)
    if (JSON.stringify(options) != "{}") {
      this.getRecord(options.recordID)
    } else {
      this.checkLocal()
    }

  },

  onShareAppMessage: function () {
    let recordID = wx.getStorageSync("exam").id
    return {
      title: this.data.exam[0].major_class + " - 考试安排",
      path: '/pages/Campus/tools/exam?recordID=' + recordID,
    }
  },


  // 下拉刷新
  onPullDownRefresh: function () {
    this.getExam()
    setTimeout(function () {
      wx.stopPullDownRefresh()
    }, 3000)
  },


  // 检查绑定状态和本地缓存
  checkLocal() {
    let account = wx.getStorageSync("account")
    this.data.account = account
    if (account == "") {
      wx.showToast({
        title: '未绑定学号',
        icon: "none",
        duration: 1500
      })
      wx.navigateTo({
        url: '/pages/Setting/login/bindStudent',
      })
    } else {
      let exam = wx.getStorageSync("exam")
      if (exam == "") {
        this.getExam()
      } else {
        this.setData({
          exam: exam.exam_list
        })
      }
    }
  },


  // 创建记录
  setRecord(exam) {
    let Obj = new wx.BaaS.TableObject("exam_share")
    let record = Obj.create()
    let info = wx.getStorageSync("student_info")
    let data = {
      sem: exam[0].year + "-" + exam[0].sem,
      major_class: exam[0].major_class,
      student: info == "" ? "" : info.name,
      exam_list: exam,
    }
    record.set(data).save().then(res => {
      wx.setStorageSync("exam", res.data)
    }, err => {
      console.log(err)
    })
  },

  // 获取记录
  getRecord(recordID) {
    let that = this
    let Obj = new wx.BaaS.TableObject("exam_share")
    Obj.get(recordID).then(res => {
      console.log(res)
      that.setData({
        exam: res.data.exam_list
      })
    }, err => {
      console.log(err)
    })
  },

  // 更新记录
  updateRecord(recordID, exam) {
    let Obj = new wx.BaaS.TableObject("exam_share")
    let record = Obj.getWithoutData(recordID)

    record.set('exam_list', exam)
    record.update().then(res => {
      wx.setStorageSync("exam", res.data)
    }, err => { })
  },


  // 处理数据
  handler(exam) {

    let tmp = wx.getStorageSync("exam")
    // 没有缓存就创建记录
    if (tmp == "") {
      this.setRecord(exam)
    } else {
      // 有缓存就更新记录
      let id = tmp.id
      this.updateRecord(id, exam)
    }
  },


  // 获取考试数据，成功后调用handler函数
  getExam() {
    let that = this
    wx.showLoading({
      title: '获取中~',
    })
    wx.request({
      url: 'https://1171058535813521.cn-shanghai.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/Spider/exam',
      method: "POST",
      data: that.data.account,
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      success: function (res) {
        console.log("exam", res)
        if (res.statusCode != 200) {
          wx.showToast({
            title: "服务器响应错误",
            icon: "none"
          })
          return
        }
        if (res.data.statusCode != 200) {
          wx.showToast({
            title: "账号或密码错误",
            icon: "none"
          })
          return
        }
        // 登录成功，渲染成绩
        that.setData({
          exam: res.data.data
        })

        if (res.data.data.length != 0)
          that.handler(res.data.data)
      },
      // 请求失败
      fail: function (err) {
        wx.showToast({
          title: "访问超时",
          icon: "none"
        })
      },
      complete: function () {
        wx.hideLoading()
      }
    })
  },


})