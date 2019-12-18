const App = require('./utils/sdk/ald-stat.js').App;
var Config = require("/utils/config.js")
var Request = require("/utils/request.js")
var startTime = Date.now(); //启动时间
require("/utils/wx.js")

App({

  globalData: {
    isAuthorized: false, //微信授权
    bindStatus: false //学号绑定
  },

  onLaunch: function (options) {
    Config.init() //初始化配置文件
    this.updata() //更新小程序

    console.log("App启动：", options)
    // 初始化知晓云
    wx.BaaS = requirePlugin('sdkPlugin')
    wx.BaaS.wxExtend(wx.login, wx.getUserInfo, wx.requestPayment)
    let ClientID = 'd5add948fe00fbdd6cdf'
    wx.BaaS.init(ClientID, {
      autoLogin: true
    })
    wx.BaaS.ErrorTracker.enable()

    if (options.scene == 1037 && JSON.stringify(options.referrerInfo) != "{}") {
      this.getAuthStatus(options.referrerInfo.extraData)
    } else {
      this.getAuthStatus()
    }

  },

  onError: function (res) {
    wx.BaaS.ErrorTracker.track(res)

    this.aldstat.sendEvent('小程序启动错误', res)
  },

  onShow: function () {
    this.aldstat.sendEvent('小程序启动时长', {
      time: Date.now() - startTime
    })
  },

  // 获取认证状态
  getAuthStatus(data = {}) {
    let that = this

    wx.getSetting({
      success: res => {
        if (res.authSetting['scope.userInfo']) {
          console.log("已授权微信")
          this.globalData.isAuthorized = true
          wx.checkSession({
            success() {
              // session_key 未过期，并且在本生命周期一直有效
            },
            fail() {
              // session_key 已经失效，需要重新执行登录流程
              wx.login() // 重新登录
            }
          })
        }
      },
      // 检测授权状态后 检测绑定状态
      complete(res) {
        wx.getStorage({
          key: 'account',
          success: function (res) {
            console.log("已绑定学号", res.data)
            that.globalData.bindStatus = true
            that.globalData.account = res.data

            // 本地无信息记录
            if (wx.getStorageSync("student_info") == "")
              Request.sync(res.data.username, res.data.password, "student_info")
          },
          fail: function (res) {
            // 来自迁移
            if (JSON.stringify(data) != "{}") {
              that.migrate(data)
            }
          }
        })
      }
    })
  },

  // 广大课表用户迁移
  migrate(data = {}) {
    wx.navigateTo({
      url: '/pages/Setting/login/bindStudent?username=' + data.username + "&password=" + data.password
    })
  },

  // 检测更新版本
  updata() {
    const updateManager = wx.getUpdateManager()

    updateManager.onCheckForUpdate(function (res) {
      // 请求完新版本信息的回调
      console.log("新版本：", res.hasUpdate)
    })

    updateManager.onUpdateReady(function () {
      wx.showModal({
        title: '更新提示',
        content: '新版本已经准备好，是否重启应用？\n如遇缓存丢失，请重启小程序。',
        success(res) {
          if (res.confirm) {
            // 新的版本已经下载好，调用 applyUpdate 应用新版本并重启
            updateManager.applyUpdate()
          }
        }
      })
    })

    updateManager.onUpdateFailed(function () {
      // 新版本下载失败
    })
  },



})



// // 展示本地存储能力
// var logs = wx.getStorageSync('logs') || []
// logs.unshift(Date.now())
// wx.setStorageSync('logs', logs)