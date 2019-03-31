// pages/Components/drawer-modal.js
Component({
  /**
   * 组件的属性列表
   */
  properties: {
    show: {
      type: Boolean,
      value: false
    },
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