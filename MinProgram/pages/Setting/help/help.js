// pages/Life/wall/wall.js
Page({

  onLoad: function () {

    wx.showLoading({
      title: 'Loading...',
    })

    let tableName = 'temp'
    let recordID = '5ce561007583606f4ffe331e'

    let Product = new wx.BaaS.TableObject(tableName)

    Product.get(recordID).then(res => {
      this.setData({
        url: res.data.data.url
      })
      setTimeout(function () {
        wx.hideLoading()
      }, 500)
    }, err => {
      // err
    })

  },
  view() {
    wx.previewImage({
      urls: this.data.url,
    })
  },

  onShareAppMessage: function () {

  }
})