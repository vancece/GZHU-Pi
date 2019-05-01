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
      overall: "校园全景",
      favorite: "我的收藏"
    }
    this.setData({
      id: options.id,
      title: title[options.id]
    })

    if (options.id == "visit") this.getVisit()

    if (options.id == "favorite") this.getFav()
  },

  onShareAppMessage: function() {

  },

  getFav() {
    let fav = wx.getStorageSync("fav_books")
    let favourite
    if (fav == "" || fav.length == 0) fav = []
    this.setData({
      fav: fav
    })
  },

  navToDetail(e) {
    let index = e.currentTarget.id
    wx.navigateTo({
      url: '/pages/Campus/library/detail?id=' + index,
    })
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
    let imgurl = "https://cos.ifeel.vip/gzhu-pi/images/resource/qrcode.jpg"
    wx.previewImage({
      urls: [imgurl]
    })
  }

})