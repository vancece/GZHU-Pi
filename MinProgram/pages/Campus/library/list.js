const Page = require('../../../utils/sdk/ald-stat.js').Page;
Page({

  data: {
    noCover: "https://cos.ifeel.vip/gzhu-pi/images/icon/book.svg",
    page: 1,
    pages: 1,
    books: []
  },

  onLoad: function(options) {
    this.getBooks(options.query)
    this.setData({
      query: options.query
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
    if (query == "") {
      wx.showToast({
        title: '请输入书名',
        icon: "none"
      })
      return
    }
    this.setData({
      query: query,
      books: [],
      page: 1
    })
    this.getBooks(query)
  },


  navToDetail(e) {
    let index = e.currentTarget.id
    wx.navigateTo({
      url: '/pages/Campus/library/detail?index=' + index,
    })
  },

  loadMore() {
    let page = this.data.page + 1
    if (page > this.data.pages) {
      wx.showToast({
        title: '没有更多啦！',
        icon: "none"
      })
      return
    }
    this.getBooks(this.data.query, page)
  },

  // 发送GET请求
  getBooks(query, page = 1) {
    let that = this
    this.setData({
      loading:true
    })

    let url = "https://1171058535813521.cn-shanghai.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/Spider/"
    wx.request({
      url: url + 'library/search?query=' + query + "&page=" + page,
      method: "get",
      success: function(res) {
        console.log(res.data.data)
        if (res.data.data.total == 0) {
          wx.showToast({
            title: '无结果',
            icon: "none"
          })
          return
        }
        that.setData({
          books: that.data.books.concat(res.data.data.books),
          pages: res.data.data.pages == 0 ? 1 : res.data.data.pages,
          page: page
        })
      },
      complete: function() {
        that.setData({
          loading: false
        })
      }
    })
  },

})