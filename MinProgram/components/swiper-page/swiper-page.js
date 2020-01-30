// components/swiper/swiper.js
Component({
  options: {
    addGlobalClass: true,
    multipleSlots: true
  },
  properties: {

  },

  /**
   * Component initial data
   */
  data: {
    isCard: false,
    list: [{
      created_by: {
        avatar: "https://ossweb-img.qq.com/images/lol/web201310/skin/big10006.jpg",
        nickname: "Shaw"
      },
      created_at: "2019年12月3日",
      image: [],
      content: "哈喽，这是内容"
    }, {
      created_by: {
        avatar: "https://ossweb-img.qq.com/images/lol/web201310/skin/big10006.jpg",
        nickname: "Shaw"
      },
      created_at: "2019年12月3日",
      image: ["https://ossweb-img.qq.com/images/lol/web201310/skin/big10006.jpg"],
      content: "哈喽，这是内容"
    }]
  },

  /**
   * Component methods
   */
  methods: {

  },
  // 生命周期方法
  lifetimes: {
    attached: function() {
      let that = this
      wx.getSystemInfo({
        success: function(res) {
          that.setData({
            bodyHeight: res.windowHeight - res.statusBarHeight - 45
          })
        },
      })
    }
  },
})