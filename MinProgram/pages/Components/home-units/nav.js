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
    ready: function() {
      this.setData({
        iconList: wx.$param["nav"],
      })
    }
  }
})