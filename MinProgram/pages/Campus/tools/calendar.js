const Page = require('../../../utils/sdk/ald-stat.js').Page
Page({

  /**
   * 页面的初始数据
   */
  data: {
    urls: [
      ["https://cos.ifeel.vip/gzhu-pi/images/campus/12.jpg"],
      ["https://cos.ifeel.vip/gzhu-pi/images/campus/2.jpg"],
      ["https://cos.ifeel.vip/gzhu-pi/images/campus/1.jpg"]
    ]
  },

  preview(e) { 
    let id = Number(e.currentTarget.id)
    wx.previewImage({
      urls: this.data.urls[id],
    })

  },
  onShareAppMessage: function() {

  }
})