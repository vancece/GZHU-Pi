
Page({

  data: {

  },

  onLoad: function(options) {
    let that = this
    wx.request({
      url: 'https://1171058535813521.cn-shanghai.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/Spider/library/visit',
      method: "GET",
      success: function(res) {
        // console.log("图书馆", res.data)
        that.setData({
          lib: res.data.data
        })
      }
    })
  },

  formSubmit(e) {
    var time = new Date()
    if (time.getHours() >= 0 && time.getHours() < 7) {
      wx.showToast({
        title: '当前时间段不可用~',
        icon: "none"
      })
      return
    }

    let query = e.detail.value.query
    wx.BaaS.wxReportTicket(e.detail.formId)
    if (query == "") {
      wx.showToast({
        title: '请输入书名',
        icon: "none"
      })
      return
    }
    wx.navigateTo({
      url: "/pages/Campus/library/list?query=" + query,
    })
  },
  onShareAppMessage: function() {

  },

  nav(e) {
    let id = e.currentTarget.id
    if (e.currentTarget.id == "") return
    wx.navigateTo({
      url: '/pages/Campus/library/sub/tools?id=' + id,
    })
  }
})