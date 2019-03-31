var utils = require("../../../../../utils/utils.js")
var showTimes = 0
Component({

  properties: {

  },

  data: {
    todayCourse: utils.getTodayCourse(),
    week: utils.getSchoolWeek(), //周数
    schoolWeek: utils.getSchoolWeek(), //校历周
    weekDays: ['日', '一', '二', '三', '四', '五', '六', ],
    weekday: new Date().getDay()
  },

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
  },
  pageLifetimes: {
    show() {
      // 初次onshow不执行
      if (showTimes) {
        this.setData({
          todayCourse: utils.getTodayCourse()
        })
      }
      showTimes++
    }
  }
})