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
    let data = that.data.account
    data["year_sem"] = wx.$param.school["year_sem"],

    wx.$ajax({
      url: "/jwxt/exam",
      data: data,
      loading: true
    })
      .then(res => {
        that.setData({
          exam: res.data
        })
        if (res.data.length != 0)
          that.handler(res.data)

      })
  },


})