const Page = require('../../../utils/sdk/ald-stat.js').Page;
Page({

  data: {
    achieve: wx.getStorageSync("achieve"),
    colors: ['cyan', 'blue', 'purple', 'mauve', 'pink', 'brown', 'red', 'orange', 'olive', 'green'],

  },
  onLoad: function(options) {

    let form = wx.getStorageSync("account")
    // let form = {
    //   username: "1916160044",
    //   password: "171611"
    // }
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
    if (type == "") {
      wx.showToast({
        title: '该类别为父级节点',
        icon: "none",
        duration: 1500
      })
      return
    }
    wx.navigateTo({
      url: '/pages/Campus/grade/list?type=' + type,
    })
  },

  getData(form) {
    wx.showLoading({
      title: '加载中...',
    })
    let that = this
    wx.request({
      url: "https://1171058535813521.cn-shenzhen.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/go/api/v1/jwxt/achieve",
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
        wx.hideLoading()
      }
    })
  }

})