const Page = require('../../../utils/sdk/ald-stat.js').Page;
var app = getApp()
Page({

  data: {
    hideSyncTip: true,
    refleshTimes: 0,
    // showAgree: true,
    showTips: false,
    bindStatus: wx.getStorageSync("account") == "" ? false : true
  },

  agree() {
    wx.setStorageSync('agree', true)
    this.setData({
      showTips: true
    })
  },
  refuse() {
    if (getCurrentPages().length == 1) {
      wx.reLaunch({
        url: '/pages/Campus/home/home',
      })
    } else {
      wx.navigateBack({
        delta: 1
      })
    }
  },
  navTo() {
    wx.navigateTo({
      url: '/pages/Setting/about/agreement',
    })
  },

  onLoad: function(options) {

    // let agree = wx.getStorageSync("agree")
    // if (agree == true) {
    //   this.setData({
    //     showAgree: false
    //   })
    // } else {
    //   return
    // }

    let that = this
    if (!this.data.bindStatus) {
      wx.navigateTo({
        url: '/pages/Setting/login/bindStudent',
      })
    } else {
      this.setData({
        account: app.globalData.account
      })
      // 从缓存读取成绩
      wx.getStorage({
        key: 'grade',
        success: function(res) {
          console.log("成绩", res)
          that.setData({
            grade: res.data,
            height: 350 + res.data.sem_list[0].grade_list.length * 170
          })
        },
        fail: function(res) {
          that.updateGrade()
        }
      })
    }
  },

  onShow: function() {
    var time = new Date()
    if (time.getHours() >= 0 && time.getHours() < 7) {
      this.setData({
        hideSyncTip: false
      })
    }
    if (!this.data.bindStatus) {
      wx.showToast({
        title: '请绑定学号',
        icon: "none",
        duration: 1500
      })
    }
  },

  // 切换学期
  swiperChange(e) {
    let current = e.detail.current
    let length = this.data.grade.sem_list[current].grade_list.length
    this.setData({
      height: length * 170 + 350
    })
  },

  // 下拉刷新
  onPullDownRefresh: function() {
    this.updateGrade()
    setTimeout(function() {
      wx.stopPullDownRefresh()
    }, 3000)
  },

  onShareAppMessage: function() {
    return {
      title: '成绩查询',
      desc: '',
      // path: '路径',
      imageUrl: "https://cos.ifeel.vip/gzhu-pi/images/pic/grade.png",
      success: function(res) {
        // 转发成功
        wx.showToast({
          title: '分享成功',
          icon: "none"
        });
      },
      fail: function(res) {
        // 转发失败
        wx.showToast({
          title: '分享失败',
          icon: "none"
        })
      }
    }
  },

  // 更新成绩
  updateGrade() {
    var time = new Date()
    if (time.getHours() >= 0 && time.getHours() < 7) {
      wx.showToast({
        title: '当前时间段不可用~',
        icon: "none"
      })
      return
    }

    if (!app.globalData.bindStatus) {
      wx.showToast({
        title: '尚未绑定学号',
        icon: "none",
        duration: 1500,
        success: function() {
          setTimeout(function() {
            wx.navigateTo({
              url: "/pages/Setting/login/bindStudent"
            })
          }, 1000)
        }
      })
      return
    }
    // 防止频繁刷新
    if (this.data.refleshTimes) {
      setTimeout(function() {
        wx.stopPullDownRefresh()
      }, 2000)

      wx.showToast({
        title: '刚刚更新过啦~',
        icon: "none",
        duration: 1800,
      })
      return
    }
    this.data.refleshTimes++;

    this.syncData()
  },



  // 同步数据
  syncData() {
    let that = this
    this.iconAnimation()
    this.setData({
      loading: true
    })
    wx.request({
      method: "POST",
      url: "https://1171058535813521.cn-shanghai.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/Spider/grade",
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      data: this.data.account,
      success: function(res) {
        if (res.statusCode != 200) {
          wx.showModal({
            title: '错误提示',
            content: '服务器响应错\n' + res.data.errorMessage,
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
        // 缓存信息
        wx.setStorage({
          key: "grade",
          data: res.data.data,
        })
        that.setData({
          grade: res.data.data,
          height: 350 + res.data.data.sem_list[0].grade_list.length * 170
        })
        wx.showToast({
          title: '更新完成 ~~ ',
          icon: "success",
          duration: 1500,
        })
      },
      fail: function(err) {
        console.log("err:", err)
        wx.showToast({
          title: "请求失败",
          icon: "none"
        })
      },
      complete: function(res) {
        that.setData({
          loading: false
        })
        if (res.statusCode == 502) {
          wx.showToast({
            title: "访问超时 " + res.statusCode,
            icon: "none"
          })
        }
        console.log(res)
        clearInterval(that.data.num) // 停止动画
        wx.stopPullDownRefresh()
      }
    })
  },

  // 图标旋转动画
  iconAnimation() {
    let that = this
    let n = 1

    function ami(n) {
      let animation = wx.createAnimation({
        duration: 1500,
        timingFunction: 'linear'
      })
      animation.rotate(18 * n).step()
      that.setData({
        animation: animation.export()
      })
    }
    this.setData({
      num: setInterval(function() {
        ami(n)
        n++
      }, 150)
    })
  }

})