const Page = require('../../../utils/sdk/ald-stat.js').Page;
Page({

  data: {
    achieve: wx.getStorageSync("achieve"),
    colors: ['cyan', 'blue', 'purple', 'mauve', 'pink', 'brown', 'red', 'orange', 'olive', 'green'],

  },
  onLoad: function (options) {

    if (this.data.achieve == "") {
      let form = wx.getStorageSync("account")
      if (form == "") {
        wx.showToast({
          title: '未绑定学号',
          icon: "none"
        })
        return
      }
      this.getData(form)
    }
  },

  // 下拉刷新
  onPullDownRefresh: function () {
    let form = wx.getStorageSync("account")
    if (form == "") {
      wx.showToast({
        title: '未绑定学号',
        icon: "none"
      })
      return
    }
    this.getData(form)
    setTimeout(function () {
      wx.stopPullDownRefresh()
    }, 3000)
  },
  onShareAppMessage: function () { },

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
    var time = new Date()
    if (time.getHours() >= 0 && time.getHours() < 7) {
      wx.showToast({
        title: '00:00~07:00不可同步',
        icon: "none"
      })
      return
    }
    let that = this
    wx.$ajax({
      url: "/jwxt/achieve",
      data: form,
      loading: true
    })
      .then(res => {
        wx.setStorageSync("achieve", res.data)
        that.setData({
          achieve: res.data
        })
      })
      .catch((e) => {})
  }

})