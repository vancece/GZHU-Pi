const Page = require('../../../utils/sdk/ald-stat.js').Page
Page({

  /**
   * 页面的初始数据
   */
  data: {
    pics: [
      ["https://cos.ifeel.vip/gzhu-pi/images/campus/12.jpg"],
      ["https://cos.ifeel.vip/gzhu-pi/images/campus/2.jpg"],
      ["https://cos.ifeel.vip/gzhu-pi/images/campus/1.jpg"]
    ]
  },

  onLoad: function() {

    wx.showLoading({
      title: 'Loading...',
    })

    let tableName = 'temp'
    let recordID = '5d030697f05d4f0d4c4338a0'

    let Product = new wx.BaaS.TableObject(tableName)

    Product.get(recordID).then(res => {
      console.log(res)
      this.setData({
        pics: res.data.data.pics
      })
      setTimeout(function() {
        wx.hideLoading()
      }, 300)
    }, err => {
      // err
    })

  },

  preview(e) {
    let id = Number(e.currentTarget.id)
    let url = [this.data.pics[id].url]
    wx.previewImage({
      urls: url,
    })
  },


  onShareAppMessage: function() {

  }

})