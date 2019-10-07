const Page = require('../../../utils/sdk/ald-stat.js').Page;
Page({

  data: {

  },
  onLoad: function (options) {
    let account = wx.getStorageSync("account")
    if (account != "") {
      this.setData({
        account: account
      })
    }
    // let agree = wx.getStorageSync("agree")
    let agree = true
    if (agree != true) {
      wx.showModal({
        title: '未经授权',
        content: '请打开成绩查询页面授权用户协议',
        success(res) {
          if (res.confirm) {
            wx.reLaunch({
              url: '/pages/Campus/grade/grade',
            })
          } else if (res.cancel) {
            wx.reLaunch({
              url: '/pages/Campus/grade/grade',
            })
          }
        }
      })
    } else {
      this.setData({
        showAgree: true
      })
      this.syncData()
      // this.setData({
      //   rank: wx.getStorageSync("rank")
      // })
    }
  },
  onShareAppMessage: function () {
    return {
      title: '成绩排名统计',
      desc: '',
      // path: '路径',
      imageUrl: "https://cos.ifeel.vip/gzhu-pi/images/pic/rank.png",
      success: function (res) {
        // 转发成功
        wx.showToast({
          title: '分享成功',
          icon: "none"
        });
      },
      fail: function (res) {
        // 转发失败
        wx.showToast({
          title: '分享失败',
          icon: "none"
        })
      }
    }
  },

  onShow: function () {
    if (wx.getStorageSync("account") == "") {
      wx.showToast({
        title: '请绑定学号',
        icon: "none",
        duration: 1500
      })
      wx.reLaunch({
        url: '/pages/Setting/login/bindStudent',
      })
    }
  },
  syncData() {
    let that = this
    wx.showLoading({
      title: '更新排名...',
    })
    wx.request({
      method: "POST",
      url: "https://1171058535813521.cn-shanghai.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/Spider/rank",
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      data: this.data.account,
      success: function (res) {
        if (res.statusCode != 200) {
          wx.showModal({
            title: '错误提示',
            content: '服务器响应错误',
          })
          return
        }
        if (res.data.statusCode != 200) {
          wx.showModal({
            title: '错误提示',
            content: '账号或密码错误',
          })
          return
        }
        that.setData({
          rank: res.data.data,
        })
        // wx.setStorageSync("rank", res.data.data)
      },
      fail: function (err) {
        console.log("err:", err)
        wx.showToast({
          title: "请求失败",
          icon: "none"
        })
      },
      complete: function (res) {
        wx.hideLoading()
        if (res.statusCode == 502) {
          wx.showToast({
            title: "访问超时 " + res.statusCode,
            icon: "none"
          })
        }
      }
    })
  },
})