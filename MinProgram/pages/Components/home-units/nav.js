// pages/Components/home-units/nav.js
Component({

  options: {
    addGlobalClass: true
  },

  properties: {

  },

  /**
   * Component initial data
   */
  data: {
    gridCol: 4,
    iconList: wx.$param["nav"],
  },

  /**
   * Component methods
   */
  methods: {
    navTo(e) {
      wx.$navTo(e)
    },

    getAppParam() {
      let tableName = 'config'
      let recordID = '5d4daa727b9e3c65e7983f54'

      let Product = new wx.BaaS.TableObject(tableName)

      Product.get(recordID).then(res => {
        console.log("在线配置：", res.data.data)
        wx.$param = res.data.data
        this.setData({
          iconList: res.data.data.nav
        })
      }, err => {
        wx.showToast({
          title: '请求出错',
          icon: "none"
        })
      })
    }
  },
  lifetimes: {
    attached: function () {
      this.getAppParam()
    }
  }
})