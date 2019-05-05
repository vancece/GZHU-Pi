const Page = require('../../../utils/sdk/ald-stat.js').Page;
var app = getApp()
Page({

  data: {
    hideSyncTip: true,
    hideLoginBtn1: false,
    hideLoginBtn2: true,
    hideLogin: false,
    hideSuccess: true,
    api: "https://1171058535813521.cn-shanghai.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/Spider/"
  },

  onLoad: function(options) {

    this.setData({
      show: !app.globalData.isAuthorized,
      hideLogin: app.globalData.bindStatus,
      hideSuccess: !app.globalData.bindStatus,
      account: options
    })

    // 用户迁移绑定
    if (!app.globalData.isAuthorized || 　JSON.stringify(options) == "{}") return
    if (!app.globalData.bindStatus && options.username != "undefined") {
      wx.showLoading({
        title: '迁移绑定...',
      })
      this.login()
    }
  },
  onShow() {
    let that = this
    var time = new Date()
    if (time.getHours() >= 0 && time.getHours() < 7) {
      this.setData({
        hideSyncTip: false
      })
    }
  },
  onReady() {
    if (app.globalData.bindStatus) {
      wx.showToast({
        title: '您已绑定学号',
        icon: "none"
      })
      this.setData({
        show: !app.globalData.isAuthorized,
        hideLogin: app.globalData.bindStatus,
        hideSuccess: !app.globalData.bindStatus,
      })
    }

  },

  userInfoHandler(data) {
    let that = this
    wx.showLoading({
      title: '授权中...',
    })

    wx.BaaS.auth.loginWithWechat(data, {
      createUser: true
    }).then(user => {
      console.log(user)
    })

    if (data.detail.errMsg == "getUserInfo:ok") {
      console.log(" 授权", data)
      wx.hideLoading()
      app.globalData.isAuthorized = true
      that.setData({
        show: false
      })
      // 用户迁移绑定
      if (JSON.stringify(this.data.account) == "{}") return
      if (!app.globalData.bindStatus && this.data.account.username != "undefined") {
        wx.showLoading({
          title: '迁移绑定...',
        })
        that.login()
      }
    } else {
      console.log("拒绝授权", data)
      wx.hideLoading()
      wx.showToast({
        title: '授权失败，可退出重试',
        icon: "none"
      })
      that.setData({
        hideLoginBtn1: false,
        hideLoginBtn2: true,
        show: false,
        showGuide: true
      })

      if (JSON.stringify(this.data.account) == "{}") return
      if (!app.globalData.bindStatus && this.data.account.username != "undefined") {
        wx.showLoading({
          title: '迁移绑定...',
        })
        that.login()
      }
    }



  },

  // 提交登录表单
  formSubmit(e) {
    // 上报formId
    wx.BaaS.wxReportTicket(e.detail.formId)
    let account = {
      username: e.detail.value.username,
      password: e.detail.value.password
    }
    if (e.detail.value.username == "" || e.detail.value.password == "") {
      wx.showToast({
        title: '用户名和密码不能为空',
        icon: 'none'
      })
    } else {
      wx.showLoading({
        title: '绑定中...',
      })
      this.setData({
        account: account,
        hideLoginBtn1: true,
        hideLoginBtn2: false,
      })
      console.log("登录", account)
      this.login() // 登录请求
    }
  },


  // 登录绑定学号
  login() {
    let that = this
    wx.request({
      method: "POST",
      url: this.data.api + "student_info",
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      data: this.data.account,
      success: function(res) {
        if (res.statusCode != 200) {
          wx.showToast({
            title: "服务器响应错误",
            icon: "none"
          })
          return
        }
        if (res.data.statusCode != 200) {
          wx.showToast({
            title: "账号或密码错误",
            icon: "none"
          })
          return
        }
        // 缓存账户信息
        wx.setStorage({
          key: 'account',
          data: that.data.account,
        })
        // 缓存学生信息
        wx.setStorage({
          key: 'student_info',
          data: res.data.data,
        })
        app.globalData.bindStatus = true
        app.globalData.account = that.data.account
        // 同步课表
        that.syncData("course")
      },
      fail: function(err) {
        console.log("err:", err)
        wx.showToast({
          title: "访问超时",
          icon: "none"
        })
      },
      complete: function() {
        wx.hideLoading()
        that.setData({
          hideLoginBtn1: false,
          hideLoginBtn2: true,
        })
      }
    })
  },

  // 同步数据
  syncData(type) {
    wx.showLoading({
      title: '同步课表中...',
    })
    wx.request({
      method: "POST",
      url: this.data.api + type,
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      data: this.data.account,
      success: function(res) {
        console.log(res)
        if (res.statusCode != 200) {
          wx.showToast({
            title: "服务器响应错误",
            icon: "none"
          })
          return
        }
        if (res.data.statusCode != 200) {
          wx.showToast({
            title: "账号或密码错误",
            icon: "none"
          })
          return
        }

        // 缓存信息
        wx.setStorage({
          key: type,
          data: res.data.data,
          success: function() {
            wx.reLaunch({
              url: '/pages/Campus/home/home',
            })
          }
        })
        wx.showToast({
          title: "同步完成",
        })
      },
      fail: function(err) {
        console.log("err:", err)
        wx.showToast({
          title: "访问超时",
          icon: "none"
        })
      },
      complete: function(res) {
        wx.hideLoading()
      }
    })
  },

  // 清除本地缓存
  cleanStorage() {
    wx.showModal({
      title: '警告',
      content: '确认操作将会清除课表、成绩等所有缓存信息!',
      success: function(res) {
        if (res.confirm) {
          wx.clearStorage({
            success: function() {
              app.globalData.bindStatus = false
              wx.showToast({
                title: '清除完成',
                duration: 1500,
                success: function() {
                  setTimeout(function() {
                    wx.reLaunch({
                      url: "/pages/Campus/home/home"
                    })
                  }, 1500)
                }
              })
            }
          })

        }
      }
    })
  },

})