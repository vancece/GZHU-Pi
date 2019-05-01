var app = getApp()
var Config = require("../../../utils/config.js")
var Setting = require("../../../utils/setting.js")

Page({
  data: {
    schedule: Config.get("schedule_mode") == "week" ? true : false,
    navColor: Config.get("schedule_mode") == "week" ? "rgba(221, 221, 221, 0.7)" : "",
    out: "ami",
    showUpdate: false,
    showDrawer: false,
    arrowUrl: "https://cos.ifeel.vip/gzhu-pi/images/icon/right-arrow.svg",
  },

  onLoad: function(options) {
    wx.showLoading({
      title: 'Loading...',
    })
  },
  onShareAppMessage: function() {

  },
  onReady() {
    wx.hideLoading()
    this.updateCheck()
  },

  updateCheck(){
    let version=Config.get("version")
    if(version <"0.8.7.20190418"){
      this.setData({
        showUpdate:true
      })
      Config.set("version","0.8.7.20190418")
    }
  },
  
  formSubmit(e) {
    console.log(e)
    wx.BaaS.wxReportTicket(e.detail.formId)
  },
  // 切换课表模式，点解悬浮图标
  switchModel() {
    if (this.data.schedule) {
      this.setData({
        schedule: !this.data.schedule,
        navColor: "",
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

  // 抽屉选项
  tapDrawer(e) {
    let drawerItem = e.currentTarget.id
    const schedule = this.selectComponent('#schedule')
    switch (drawerItem) {
      case "changeBg":
      case "changeMode":
        this.setData({
          drawerItem: drawerItem == this.data.drawerItem ? null : drawerItem,
          checkedBlur: Config.get("blur"),
          mode: Config.get("schedule_mode")
        })
        break
      case "selectImg":
        Setting.setBg().then(res => {
          Config.set("schedule_bg", res)
          schedule.updateBg()
        })
        break
      case "white,white":
      case "#ddd,#ddd":
      case "#d299c2,#fef9d7":
      case "#a8edea,#fed6e3":
        Config.set("schedule_bg", drawerItem)
        schedule.updateBg()
        break
      case "navToAbout":
        wx.navigateTo({
          url: '/pages/Setting/about/about',
        })
        break
      case "navToSync":
        wx.navigateTo({
          url: "/pages/Setting/login/sync",
        })
    }
  },

  // 开启关闭高斯模糊
  switchChange(e) {
    if (e.detail.value) Config.set("blur", 8)
    else Config.set("blur", 0)
    const schedule = this.selectComponent('#schedule')
    schedule.updateBg()
  },

  // 切换课表模式
  radioChange(e) {
    Config.set("schedule_mode", e.detail.value)

    if (e.detail.value == "day") this.data.schedule = true
    else this.data.schedule = false
    this.switchModel()
  },
  catchtap(e) {},


})