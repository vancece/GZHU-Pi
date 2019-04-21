// pages/Campus/navTools/calendar.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    urls: [
      ["https://cos.ifeel.vip/gzhu-pi/images/campus/1.jpg"],
      ["https://cos.ifeel.vip/gzhu-pi/images/campus/2.jpg"],
      ["https://cos.ifeel.vip/gzhu-pi/images/campus/calendar_new.png"]
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