const Page = require('../../../utils/sdk/ald-stat.js').Page;
Page({

  data: {
    achieve: wx.getStorageSync("achieve"),
    colors: ['cyan', 'blue', 'purple', 'mauve', 'pink', 'brown', 'red', 'orange', 'olive', 'green'],

  },
  onLoad: function(options) {

    let form = wx.getStorageSync("account")
    if (form == "") {
      wx.showToast({
        title: '未绑定学号',
        icon: "none"
      })
      return
    }
    this.getData(form)
  },


  onShareAppMessage: function() {},

  navTo(e) {
    let type = e.currentTarget.dataset.type
    wx.navigateTo({
      url: '/pages/Campus/grade/list?type=' + type,
    })
  },

  getData(form) {
    let that = this
    wx.request({
      url: "https://gzhu.ifeel.vip/api/v1/jwxt/achieve",
      method: "POST",
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      data: form,
      success: function(res) {
        if (res.data.status != 200) {
          wx.showToast({
            title: res.data.msg,
            icon: "none"
          })
          return
        }
        wx.setStorageSync("achieve", res.data.data)
        that.setData({
          achieve: res.data.data
        })
      },
      fail: function(err) {
        wx.showModal({
          title: '请求失败',
          content: "错误信息:" + err.errMsg,
        })
      },
      complete(res) {
        console.log(res.data)
      }
    })
  }

})