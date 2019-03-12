var app = getApp()
// var schedule = true

Page({

  data: {
    navColor: "rgba(221, 221, 221, 0.7)",
    schedule: true,
    out: "ami",
    showDrawer: false
  },

  onLoad: function(options) {


  },
  onShareAppMessage: function() {

  },

  // 切换课表模式
  switchModel() {
    if (this.data.schedule) {
      this.setData({
        schedule: !this.data.schedule,
        navColor: "white",
      })
    } else {
      this.setData({
        schedule: !this.data.schedule,
        navColor: "rgba(221, 221, 221, 0.7)",
      })
    }

  },

  // 打开抽屉弹窗
  openDrawer() {
    this.setData({
      showDrawer: true
    })
  },






})