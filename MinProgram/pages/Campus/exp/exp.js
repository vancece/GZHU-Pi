var Utils = require("../../../utils/utils.js")
var Data = require("../../../utils/data.js")
var Config = require("../../../utils/config.js")
var Setting = require("../../../utils/setting.js")
Page({


  data: {
    colors: Data.colors,
  },


  onLoad: function(options) {},
  onShareAppMessage: function() {},
  onShow: function() {
    this.setData({
      exp: this.reSort(),
      checked: Config.get("showExp")
    })
  },

  navToSync() {
    wx.navigateTo({
      url: "/pages/Setting/login/sync?id=0",
    })
  },

  switchChange(e) {
    let that = this
    if (this.data.exp == "") {
      wx.showToast({
        title: '未同步实验',
        icon: 'none'
      })
      this.setData({
        checked: false
      })
      return
    }
    if (this.data.checked) {
      Config.set("showExp", false)
      this.data.checked = false
      return
    }
    wx.showModal({
      title: '提示',
      content: '加到课表中有可能会覆盖其它课程或被其它课程覆盖,如课表上有对应实验大课程，可先将其删除！',
      success: function(res) {
        if (res.confirm) {
          Config.set("showExp", true)
          that.setData({
            checked: true
          })
        } else {
          that.setData({
            checked: false
          })
        }
      }
    })
  },

  // 课程重新排序
  reSort() {
    let exp = wx.getStorageSync("exp")
    if (exp == "") return exp

    let time = new Date()
    time = Utils.formatTime(time).split(" ")[0]

    for (let i = 0; i < exp.length; i++) {
      if (time > exp[0].date) {
        let tmp = exp.splice(0, 1)[0]
        tmp.color = 0
        exp = exp.concat(tmp)
      }
    }

    return exp
  }

})