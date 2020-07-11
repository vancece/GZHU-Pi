const Page = require('../../../utils/sdk/ald-stat.js').Page;
let interstitialAd = null

Page({

  data: {
    achieve: wx.getStorageSync("achieve"),
    colors: ['cyan', 'blue', 'purple', 'mauve', 'pink', 'brown', 'red', 'orange', 'olive', 'green'],

  },

  onLoad: function (options) {
    this.initAD()
    if (this.data.achieve == "") {
      let form = wx.getStorageSync("account")
      if (form == "") {
        wx.showToast({
          title: '未绑定学号',
          icon: "none"
        })
        return
      }
      this.getData(form)
    }
    this.count()
  },

  // 下拉刷新
  onPullDownRefresh: function () {
    let form = wx.getStorageSync("account")
    if (form == "") {
      wx.showToast({
        title: '未绑定学号',
        icon: "none"
      })
      return
    }
    this.getData(form)
    setTimeout(function () {
      wx.stopPullDownRefresh()
    }, 3000)
  },
  onShareAppMessage: function () {},

  initAD() {
    if (wx.createInterstitialAd) {
      interstitialAd = wx.createInterstitialAd({
        adUnitId: 'adunit-5336c97b15452cc2'
      })
      interstitialAd.onLoad(() => {
        console.log('onLoad event emit')
      })
      interstitialAd.onError((err) => {
        console.log('onError event emit', err)
      })
      interstitialAd.onClose((res) => {
        console.log('onClose event emit', res)
      })
    }
  },

  insertAD() {
    if (interstitialAd) {
      interstitialAd.show().catch((err) => {
        console.error(err)
      })
    }
  },

  onShow: function () {
    let that = this
    setTimeout(function () {
      that.insertAD()
    }, 1500)
  },

  navTo(e) {
    let type = e.currentTarget.dataset.type
    if (type == "") {
      wx.showToast({
        title: '该类别为父级节点',
        icon: "none",
        duration: 1500
      })
      return
    }
    wx.navigateTo({
      url: '/pages/Campus/grade/list?type=' + type,
    })
  },

  getData(form) {
    var time = new Date()
    if (time.getHours() >= 0 && time.getHours() < 7) {
      wx.showToast({
        title: '00:00~07:00不可同步',
        icon: "none"
      })
      return
    }
    let that = this
    wx.$ajax({
        url: "/jwxt/achieve",
        data: form,
        loading: true
      })
      .then(res => {
        wx.setStorageSync("achieve", res.data)
        that.setData({
          achieve: res.data
        })

        that.count()
      })
      .catch((e) => {})
  },

  count() {
    let achieve = wx.getStorageSync("achieve")

    if (!achieve || achieve.length == 0) {
      return
    }

    let required = achieve[0].required
    let acquired = 0

    for (let i = 0; i < achieve.length; i++) {
      if (achieve[i].type == "必修类" || achieve[i].type == "选修类") {
        acquired = acquired + Number(achieve[i].acquired)
      }
    }

    this.setData({
      acquired: acquired,
      required: required
    })

  }

})