Page({

  data: {
    noCover: "cloud://gzhu-pi-f63be3.677a-gzhu-pi-f63be3/images/icon/book.svg",
    page: 1
  },

  onLoad: function(options) {
    this.getBooks(options.query)
    this.setData({
      query: options.query
    })

    let curPage = getCurrentPages()
    let prePage = curPage[curPage.length - 2]
    console.log(curPage,555)
    console.log(curPage, 555)
    console.log(122,prePage.data)

  },

  formSubmit(e) {
    let query = e.detail.value.query
    if (query == "") {
      wx.showToast({
        title: '请输入书名',
        icon: "none"
      })
      return
    }
    this.setData({
      query: query
    })
    this.getBooks(query)
  },

  next() {
    this.data.page = this.data.page + 1
    this.getBooks(this.data.query, this.data.page)
  },


  navToDetail(e) {
    let id = Number(e.currentTarget.id)
    let arg = JSON.stringify(this.data.books[id])
    arg = arg.replace("=", "")
    wx.navigateTo({
      url: '/pages/Campus/library/detail?arg=' + arg,
    })

  },

  // 发送GET请求
  getBooks(query, page = 1) {
    let that = this
    wx.showLoading({
      title: '...',
    })

    let url = "https://1171058535813521.cn-shanghai.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/Spider/"
    let url2 = "http://192.168.8.1:5000/"
    wx.request({
      url: url + 'library/search?query=' + query + "&page=" + page,
      method: "get",
      success: function(res) {
        console.log(res.data.data)
        that.setData({
          books: res.data.data.books
        })
      },
      complete: function() {
        wx.hideLoading()
      }
    })
  },

})