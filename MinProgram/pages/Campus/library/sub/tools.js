Page({

  data: {
    id: "open",
    title: "开放时间"
  },

  onLoad: function(options) {
    let title = {
      room: "库室分布",
      open: "开放时间",
      visit: "进馆数据",
      overall: "校园全景"
    }
    this.setData({
      id: options.id,
      title: title[options.id]
    })
    this.getVisit()
  },

  onShareAppMessage: function() {

  },

  getVisit() {
    let that = this
    // 更新图书馆进馆信息
    wx.request({
      url: 'https://myapi.iego.net:5000/library',
      method: "GET",
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      success: function(res) {
        that.setData({
          lib: res.data
        })
      }
    })
  },

  preview() {
    let imgurl = "cloud://gzhu-pi-f63be3.677a-gzhu-pi-f63be3/images/res/qrcode.jpg"
    wx.previewImage({
      urls: [imgurl]
    })
  }

})