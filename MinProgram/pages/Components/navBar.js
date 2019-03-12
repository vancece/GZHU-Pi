Component({

  properties: {
    // 页面标题
    title: {
      type: String,
      value: "Hello 广大"
    },

    // 标题颜色
    titleColor: {
      type: String,
      value: "black"
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
    }
  },


  data: {
    titleBarHeight: 40 //标题栏高度
  },

  methods: {
    // 返回上一页
    navigateBack() {
      // console.log(getCurrentPages().length)
      if (getCurrentPages().length == 1) {
        wx.reLaunch({
          url: '/pages/Campus/home/home',
        })
      } else(wx.navigateBack({
        delta: 1
      }))
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
            navBarHeight: res.statusBarHeight + 40 // 小程序导航栏高度
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