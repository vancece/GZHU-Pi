const Page = require('../../../utils/sdk/ald-stat.js').Page;
var Request = require("../../../utils/request.js")
Page({

  data: {
    current: 1,
    sem_list: ["2018-2019-1", "2018-2019-2", "2019-2020-1", "2019-2020-2"],
    pickerIndex: 2,
    loading: false,
    exp_btn: "同步实验",
    kb_btn: "同步课表",

    update_time: wx.getStorageSync("course").update_time + " +8h",
    cet: wx.getStorageSync("cet"),
    account: wx.getStorageSync("account"),
    exp_account: wx.getStorageSync("exp_account"),
  },


  onLoad: function(options) {
    // 根据来源的不同显示不同界面
    if (JSON.stringify(options) != "{}") {
      this.setData({
        current: Number(options.id)
      })
    }
  },

  switch (e) {
    this.setData({
      current: Number(e.target.id)
    })
  },
  swiperChange(e) {
    this.setData({
      current: e.detail.current
    })
  },
  pickerChange(e) {
    this.setData({
      pickerIndex: Number(e.detail.value)
    })
  },

  // 提交表单
  formSubmit(e) {
    let that = this
    let id = e.detail.target.id
    // 上报formId
    wx.BaaS.wxReportTicket(e.detail.formId)
    if (e.detail.value.username == "" || e.detail.value.password == "") {
      wx.showToast({
        title: '用户名、密码不能为空',
        icon: "none"
      })
      return
    }
    // 防止多次点击
    if (this.data.loading) {
      return
    }
    this.setData({
      loading: true
    })

    switch (id) {

      // 课表
      case "kb":
        wx.showModal({
          title: '提示',
          content: '同步将会覆盖当前课表',
          success: function(res) {
            if (res.confirm) {
              Request.sync(e.detail.value.username, e.detail.value.password, "course", "account", that.data.sem_list[that.data.pickerIndex]).then(res => {
                wx.showToast({
                  title: res,
                  icon: res == "同步完成" ? "success" : "none"
                })
                that.setData({
                  loading: false,
                  update_time: wx.getStorageSync("course").update_time + " +8h",
                  kb_btn: res == "同步完成" ? "已同步" : "同步课表",
                })
              })
            } else {
              that.setData({
                loading: false
              })
            }
          }
        })
        break


        // 实验
      case "exp":
        Request.sync(e.detail.value.username, e.detail.value.password, "exp", "exp_account").then(res => {
          wx.showToast({
            title: res,
            icon: res == "同步完成" ? "success" : "none"
          })
          that.setData({
            loading: false,
            exp_btn: res == "同步完成" ? "已同步" : "同步实验",
          })
        })
        break


        // 四六级
      case "cet":
        wx.setStorageSync("cet", e.detail.value.cet)
        wx.showToast({
          title: '预存成功！',
        })
        that.setData({
          loading: false
        })
        break
    }
  },

})