const Page = require('../../../utils/sdk/ald-stat.js').Page;
var app = getApp()
Page({

  data: {
    hideSyncTip: true,
    hideLoginBtn1: false,
    hideLoginBtn2: true,
    hideLogin: false,
    hideSuccess: true,
    checked: true,
    api: "https://1171058535813521.cn-shanghai.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/Spider/"
  },

  onLoad: function (options) {

    this.setData({
      show: !app.globalData.isAuthorized,
      hideLogin: app.globalData.bindStatus,
      hideSuccess: !app.globalData.bindStatus,
      account: options
    })

    // 用户迁移绑定
    if (!app.globalData.isAuthorized || JSON.stringify(options) == "{}") return
    console.log(!app.globalData.bindStatus, options.username)
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
  navToAgreement() {
    wx.navigateTo({
      url: '/pages/Setting/about/agreement',
    })
  },
  agree() {
    this.setData({
      checked: !this.data.checked
    })
  },

  userInfoHandler(data) {
    let that = this
    wx.showLoading({
      title: '授权中...',
    })
    wx.BaaS.auth.loginWithWechat(data, {
      createUser: true,
      syncUserProfile: "overwrite"
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
    }
  },

  // 提交登录表单
  formSubmit(e) {
    if (!this.data.checked) {
      wx.showToast({
        title: '请同意用户协议',
        icon: 'none'
      })
      return
    }
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
    this.setData({
      loading: true
    })
    wx.request({
      method: "POST",
      url: this.data.api + "student_info",
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      data: this.data.account,
      success: function (res) {
        console.log("student_info:",res)
        if (res.statusCode != 200) {
          wx.showToast({
            title: "服务器响应错误",
            icon: "none"
          })
          return
        }
        if (res.data.statusCode != 200) {
          that.setData({
            loading: false
          })
          wx.showToast({
            title: "账号或密码错误",
            icon: "none"
          })
          wx.showModal({
            title: '登录提示',
            content: '账号或密码错误，如需修改请前往  http://my.gzhu.edu.cn',
            showCancel: false
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
      fail: function (err) {
        console.log("err:", err)
        that.setData({
          loading: false
        })
        wx.showToast({
          title: "访问超时",
          icon: "none"
        })
      },
      complete: function () {
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
    let that = this
    let form = this.data.account
    form["year_sem"] = wx.$param.school["year_sem"]

    wx.$ajax({
      url: "/jwxt/course",
      data: form
    })
      .then(res => {
        wx.showToast({
          title: "同步完成",
        })
        that.setData({
          loading: false,
        })
        // 缓存账户信息
        delete form["year_sem"]
        wx.setStorageSync("account", form)
        // 缓存结果数据
        res.data["update_time"] = res.update_time
        wx.setStorageSync("course", res.data)

        setTimeout(function () {
          wx.reLaunch({
            url: '/pages/Campus/home/home',
          })
        }, 300)

      }).catch(err => {
        that.setData({
          loading: false
        })
      })
  },

  // 清除本地缓存
  cleanStorage() {
    wx.showModal({
      title: '警告',
      content: '确认操作将会清除课表、成绩等所有缓存信息!',
      success: function (res) {
        if (res.confirm) {
          wx.clearStorage({
            success: function () {
              app.globalData.bindStatus = false
              wx.showToast({
                title: '清除完成',
                duration: 1500,
                success: function () {
                  setTimeout(function () {
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