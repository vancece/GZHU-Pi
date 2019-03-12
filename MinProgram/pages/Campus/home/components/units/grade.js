Component({

  properties: {

  },

  data: {},

  methods: {

  },

  lifetimes: {
    attached: function() {
      let that = this
      wx.getStorage({
        key: 'grade',
        success: function(res) {
          that.setData({
            grade: res.data
          })
        }
      })
    }
  },

  pageLifetimes: {
    show() {
      let that = this
      if (this.data.grade == undefined) {
        wx.getStorage({
          key: 'grade',
          success: function(res) {
            that.setData({
              grade: res.data
            })
          }
        })
      }
    },
    hide() {
      // 页面被隐藏
    },
    resize(size) {
      // 页面尺寸变化
    }
  }

})