// pages/Campus/navTools/calendar.js
Page({

  /**
   * 页面的初始数据
   */
  data: {

  },

  preview(e) {
    let urls = [
      ["https://677a-gzhu-pi-f63be3-1258677591.tcb.qcloud.la/images/campus/1.jpg"],
      ["https://677a-gzhu-pi-f63be3-1258677591.tcb.qcloud.la/images/campus/2.jpg"]
      
    ]

    if (e.currentTarget.id == "0") {
      wx.previewImage({
        urls: urls[0],
      })
    }
    if (e.currentTarget.id == "1") {
      wx.previewImage({
        urls: urls[1],
      })
    }
  },
  onShareAppMessage: function () {

  }
})