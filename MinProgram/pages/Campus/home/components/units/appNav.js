
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
          
          break
        case "2":
          wx.navigateTo({
            url: '/pages/Campus/grade/grade',
          })
          break
        case "3":
          
          break
        case "4":
         
          break
        case "5":
          wx.navigateTo({
            url: '/pages/Campus/evaluation/evaluation',
          })
          break
        case "6":
          // wx.navigateTo({
          //   url: '/pages/Campus/welfare/welfare',
          // })
          this.triggerEvent('touchEvent')
          break
        case "7":
          wx.navigateTo({
            url: '/pages/Setting/login/bindStudent',
          })
          break
      }
    },
  }
})
