// pages/Components/drawer-modal.js
Component({
  options: {
    multipleSlots: true // 在组件定义时的选项中启用多slot支持
  },
  properties: {
    show: {
      type: Boolean,
      value: false
    },
    mode: {
      type: String,
      value: "left"
    }
  },

  data: {
    hide: true //初始化隐藏
  },
  methods: {

    cancel() {
      this.setData({
        show: false
      })
    },

    confirm() {
      this.setData({
        show: false
      })
      this.triggerEvent('confirm')
    }

  },

  // 生命周期方法
  lifetimes: {
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
})