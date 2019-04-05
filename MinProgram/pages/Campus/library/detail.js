Page({

  data: {
    noCover: "cloud://gzhu-pi-f63be3.677a-gzhu-pi-f63be3/images/icon/book.svg",
    exist: false, //豆瓣藏书状态
    douban: {},
    holdings: []
  },

  onLoad: function(options) {
    // 根据点击获取上一页面对应书本的信息
    let curPage = getCurrentPages()
    let prePage = curPage[curPage.length - 2]
    let book = prePage.data.books[Number(options.index)]
    this.setData({
      book: book
    })

    if (book.ISBN != "")
      this.douban(book.ISBN)
    if (book.copies != 0)
      this.getHoldings(book.id, book.source)
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
      // url: 'https://douban.uieee.com/v2/book/isbn/' + ISBN,
      url: "https://api.douban.com/v2/book/isbn/" + ISBN + "?apikey=0b2bdeda43b5688921839c8ecb20399b",
      method: "get",
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      success: function(res) {
        console.log("豆瓣", res.data)
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
    var time = new Date()
    if (time.getHours() >= 0 && time.getHours() < 7) {
      wx.showToast({
        title: '当前时间段不可用~',
        icon: "none"
      })
      return
    }
    
    let that = this
    let url = "https://1171058535813521.cn-shanghai.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/Spider/"
    wx.request({
      url: url + 'library/holdings?id=' + id + "&source=" + source,
      method: "get",
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      success: function(res) {
        console.log("馆藏信息", res.data)
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