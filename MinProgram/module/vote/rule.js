import Utils from '../utils.js';
Utils.initSdk();
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
 
    console.log(options)
    this.setData({
      url: options.url
    })
  },

  navBack() {
    wx.navigateBack({
      delta: 1
    })
  }
})