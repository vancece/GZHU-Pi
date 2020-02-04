Component({

  properties: {
    // 页面标题
    title: {
      type: String,
      value: "广大派"
    },

    // 标题颜色
    titleColor: {
      type: String,
      value: "#333"
    },

    // 顶栏背景颜色，支持rgba
    navColor: {
      type: String,
      value: " "
    },

    // 隐藏标题
    hideTitle: {
      type: Boolean,
      value: false
    },

    // 隐藏返回按键
    hideBackBtn: {
      type: Boolean,
      value: false
    },

    // 导航栏是fixed布局，是否空白占位
    occupy: {
      type: Boolean,
      value: true
    },

    // 如果当前页面是最顶层，可以自定义返回的页面
    redirectTo: {
      type: String,
      value: '/pages/Campus/home/home'
    }
  },

  data: {
    titleBarHeight: 40 //标题栏高度
  },

  methods: {
    // 返回上一页
    navigateBack() {
      if (getCurrentPages().length == 1) {
        wx.navigateTo({
          url: this.data.redirectTo,
          fail: err => {
            wx.switchTab({
              url: this.data.redirectTo,
            })
          }
        })
      } else {
        wx.navigateBack({
          delta: 1
        })
      }
    },
  },

  // 生命周期方法
  lifetimes: {
    created: function() {},

    attached: function() {
      let that = this
      wx.getSystemInfo({
        success: function(res) {
          that.setData({
            statusBarHeight: res.statusBarHeight, //状态栏高度
            navBarHeight: res.statusBarHeight + 45 // 小程序导航栏高度
          })
        },
      })
    }
  },

  // 组件所在页面的生命周期
  pageLifetimes: {
    show() {
      // 页面被展示
    },
    hide() {
      // 页面被隐藏
    },
    resize(size) {
      // 页面尺寸变化
    }
  }
})