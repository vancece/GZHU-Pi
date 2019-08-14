Page({

  onLoad: function(options) {
    let tableName = 'config'
    let recordID = '5d261b2c0ec42259a6dbd077'

    let Product = new wx.BaaS.TableObject(tableName)

    Product.get(recordID).then(res => {
      this.setData({
        section: res.data.data.section
      })
    }, err => {
      wx.showToast({
        title: '请求出错',
        icon: "none"
      })
    })
  },

})