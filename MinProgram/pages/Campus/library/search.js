// pages/Campus/library/search.js
Page({


  data: {

  },


  onLoad: function(options) {
    let that = this
    wx.request({
      url: 'https://myapi.iego.net:5000/library',
      method: "GET",
      success: function(res) {
        // console.log("图书馆信息", res.data)
        that.setData({
          lib: res.data
        })
      }
    })
  },

  formSubmit(e) {
    let query = e.detail.value.query
    // wx.BaaS.wxReportTicket(e.detail.formId)
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

  nav(){
    wx.navigateTo({
      url: '/pages/Campus/library/sub/overview',
    })
  }
})