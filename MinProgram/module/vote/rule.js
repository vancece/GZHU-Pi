// const Page = require('../sdk/ald-stat.js').Page;
// module/vote/rule.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    statusBarHeight: wx.getSystemInfoSync().statusBarHeight,
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function(options) {
    let that = this
    require('../sdk/sdk-v2.0.6-a')
    let ClientID = 'd5add948fe00fbdd6cdf'
    wx.BaaS.init(ClientID, {
      autoLogin: true
    })

    wx.showLoading({
      title: 'Loading...',
    })
    let tableName = 'temp'
    let recordID = '5cc98999574c645dd667fcb2'

    let Product = new wx.BaaS.TableObject(tableName)

    Product.get(recordID).then(res => {
      console.log(res)
      this.setData({
        url: res.data.data.url
      })
      wx.hideLoading()
    }, err => {
      wx.hideLoading()
    })

  },

  navBack() {
    wx.navigateBack({
      delta: 1
    })
  }
})