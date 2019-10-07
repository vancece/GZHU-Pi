const Page = require('../../../utils/sdk/ald-stat.js').Page;
Page({

  data: {
    noCover: "cloud://gzhu-pi-f63be3.677a-gzhu-pi-f63be3/images/icon/book.svg",
    exist: false, //豆瓣藏书状态
    douban: {},
    holdings: []
  },

  onLoad: function(options) {

    // 从收藏页面进入 options.id != undefined 
    if (options.id != undefined) {
      let fav = wx.getStorageSync("fav_books")
      let fav_book = fav[options.id]

      this.setData({
        book: fav_book
      })
      this.checkFav()
      this.douban(fav_book.ISBN)
      this.getHoldings(fav_book.id, fav_book.source)
      return
    }

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

    this.checkFav()
  },

  checkFav() {
    let fav = wx.getStorageSync("fav_books")
    if (fav == "") {
      this.setData({
        favorite: false
      })
      return
    }
    let flag = false
    for (let i = 0; i < fav.length; i++) {
      if (fav[i].ISBN == this.data.book.ISBN) {
        flag = true
      }
    }
    this.setData({
      favorite: flag
    })
  },

  // 点击图标收藏
  favorite() {
    if (this.data.book.ISBN == "") {
      wx.showToast({
        title: '无ISBN，无法收藏',
        icon: 'none'
      })
      return
    }

    let book = {
      id: this.data.book.id,
      source: this.data.book.source,
      ISBN: this.data.book.ISBN,
      book_name: this.data.book.book_name,
      call_No: this.data.book.call_No,
      author: this.data.book.author,
      publisher: this.data.exist ? this.data.douban.publisher : this.data.book.publisher,
      image: this.data.exist ? this.data.douban.image : this.data.noCover
    }

    let fav = wx.getStorageSync("fav_books")
    let flag = -1

    if (fav == "") {
      wx.setStorageSync("fav_books", [book])
    } else {
      for (let i = 0; i < fav.length; i++) {
        if (fav[i].ISBN == this.data.book.ISBN) {
          flag = i
        }
      }

      if (flag != -1) {
        // 取消收藏
        fav.splice(flag, 1)
        wx.setStorageSync("fav_books", fav)
      } else {
        // 加入收藏
        fav[fav.length] = book
        wx.setStorageSync("fav_books", fav)
      }
    }
    this.checkFav()
  },


  // 豆瓣获取图书信息
  douban(ISBN) {
    if (ISBN == "") {
      return
    }
    this.setData({
      loading: true
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
        that.setData({
          loading: false
        })
      }
    })
  },

  // 获取馆藏信息
  getHoldings(id, source) {
    var time = new Date()
    if (time.getHours() >= 0 && time.getHours() < 7) {
      wx.showToast({
        title: '当前数据馆藏查询不可用~',
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