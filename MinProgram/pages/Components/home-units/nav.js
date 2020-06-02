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
  },

  lifetimes: {
    ready: function () {
      this.setData({
        iconList: wx.$param["nav"],
      })

      let count = 0
      for (let i = 0; i < wx.$param["nav"].length; i++) {
        if (wx.$param["nav"][i].show == true) {
          count++
        }
      }
      if (count % 5 == 0) {
        this.setData({
          gridCol: 5
        })
      }
    }
  }
})