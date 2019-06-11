// pages/Components/tips.js
Component({

  properties: {

  },

  data: {
    res: {}
  },

  methods: {
    close() {
      let now = new Date().getTime()
      wx.setStorageSync("last_time", now)

      let res = {
        show: false
      }
      this.setData({
        res: res
      })
    },
    nav() {
      wx.navigateTo({
        url: this.data.res.url,
      })
    }

  },
  lifetimes: {
    attached: function() {
      let now = new Date().getTime()
      let loc = wx.getStorageSync("last_time")
      let intv = 0
      if (loc != "") intv = now - Number(loc)
      // console.log(intv)
      let that = this
      wx.BaaS.invokeFunction('tips').then(res => {
        if (intv > res.data["time"] || intv == 0) {
          that.setData({
            res: res.data
          })
        }
      })
    }
  },
})