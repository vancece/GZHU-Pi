const Page = require('../../../utils/sdk/ald-stat.js').Page;
Page({

  data: {
    list: [],
    title: "修读列表"
  },

  onLoad: function(options) {
    if (JSON.stringify(options) == "{}") return
    if (options.type == undefined) return

    let achieve = wx.getStorageSync("achieve")
    if (achieve == "") return
    // 根据传入的类型从缓存中提取对象
    let obj = achieve.find(function(obj) {
      if (obj.type == options.type)
        return obj
    })
    if (obj==undefined || obj.items == undefined) return

    this.setData({
      list: obj.items,
      title: obj.type
    })
  },

  // 点击列表
  tapItem(e) {
    let id = Number(e.currentTarget.id)
    this.setData({
      detail: this.data.list[id],
      showDetail: true
    })
  },

})