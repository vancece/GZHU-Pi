Component({
  /**
   * 组件的属性列表
   */
  properties: {

  },

  /**
   * 组件的初始数据
   */
  data: {

  },

  /**
   * 组件的方法列表
   */
  methods: {
    navigate(e) {

      switch (e.currentTarget.id) {
        case "0":
          wx.navigateTo({
            url: '/pages/Campus/tools/calendar',
          })
          break
        case "1":
          wx.navigateTo({
            url: '/pages/Campus/tools/emptyRoom',
          })
          break
        case "2":
          wx.navigateTo({
            url: '/pages/Campus/grade/grade',
          })
          break
        case "3":
          wx.navigateTo({
            url: '/pages/Campus/library/search',
          })
          break
        case "4":
          wx.navigateTo({
            url: '/pages/Campus/tools/exam',
          })
          break
        case "5":
          wx.navigateTo({
            url: '/pages/Campus/course/search',
          })
          break
        case "6":
          wx.navigateTo({
            url: '/pages/Campus/exp/exp',
          })
          break
        case "7":
          wx.navigateTo({
            url: '/module/vote/vote',
          })
          break
      }
    },
  }
})