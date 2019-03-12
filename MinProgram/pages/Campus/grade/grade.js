var app = getApp()
Page({

  data: {
    hideSyncTip: true,
    refleshTimes: 0,
  },

  onLoad: function(options) {
    let that = this
    if (!app.globalData.bindStatus) {
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
    if (!app.globalData.bindStatus) {
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

  onShareAppMessage: function() {},

  // 更新成绩
  updateGrade() {
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
    wx.showLoading({
      title: '更新成绩...',
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
          key: "grade",
          data: res.data.data,
        })
        that.setData({
          grade: res.data.data,
          height: 350 + res.data.data.sem_list[0].grade_list.length * 170
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
        console.log(res.data.data)
        clearInterval(that.data.num) // 停止动画
        wx.hideLoading()

        wx.showToast({
          title: '更新完成 ~~ ',
          icon: "success",
          duration: 1500,
        })
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