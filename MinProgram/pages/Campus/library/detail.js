Page({

  data: {
    noCover: "cloud://gzhu-pi-f63be3.677a-gzhu-pi-f63be3/images/icon/book.svg",
    exist: false, //豆瓣藏书状态
    douban: {},
    // holdings: wx.getStorageSync("detail").book_meg //馆藏信息
  },

  onLoad: function(options) {
    console.log(options)
    let arg = JSON.parse(options.arg)
    this.setData({
      book: arg
    })
    this.douban(arg.ISBN)
    this.getHoldings(arg.id, arg.source)
  },

  // 豆瓣获取图书信息
  douban(ISBN) {
    if (ISBN == "") {
      return
    }
    wx.showLoading({
      title: '...',
    })
    let that = this
    wx.request({
      url: 'https://douban.uieee.com/v2/book/isbn/' + ISBN,
      method: "get",
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      success: function(res) {
        console.log(res)
        if (res.statusCode != 200) return
        // 豆瓣查找成功
        if (res.data.code == undefined) {
          that.setData({
            douban: res.data,
            exist: true
          })
        } else {
          that.setData({
            douban: {},
            exist: false
          })
        }
      },
      fail: function(res) {},
      complete: function(res) {
        wx.hideLoading()
      }
    })
  },

  // 获取馆藏信息
  getHoldings(id, source) {
    let that = this

    let url = "https://1171058535813521.cn-shanghai.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/Spider/"
    let url2 = "http://192.168.8.1:5000/"
    wx.request({
      url: url + 'library/holdings?id=' + id + "&source=" + source,
      method: "get",
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      success: function(res) {
        console.log(res)
        if (res.statusCode != 200) return
        that.setData({
          holdings: res.data.data,
        })
      },
      fail: function(res) {},
      complete: function(res) {
        wx.hideLoading()
      }
    })
  },

})