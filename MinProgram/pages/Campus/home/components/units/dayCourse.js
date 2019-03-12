var utils = require("../../../../../utils/utils.js")
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
    todayCourse: utils.getTodayCourse(),
    week: utils.getSchoolWeek(), //周数
    schoolWeek: utils.getSchoolWeek(), //校历周
    weekDays: ['日', '一', '二', '三', '四', '五', '六', ],
    weekday: new Date().getDay()
  },

  /**
   * 组件的方法列表
   */
  methods: {
    nav() {
      wx.navigateTo({
        url: '/pages/Setting/login/bindStudent',
      })
    },


  },

  lifetimes: {
    attached: function() {
      let that = this
      wx.getStorage({
        key: 'account',
        success: function(res) {
          that.setData({
            account: true
          })
        },
        fail: function(res) {
          that.setData({
            account: false
          })
        }
      })
    },

  }
})