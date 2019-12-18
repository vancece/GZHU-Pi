const Page = require('../../../utils/sdk/ald-stat.js').Page;
Page({

  /**
   * 页面的初始数据
   */
  data: {
    title: "我的发布",
    limit: 10,
    offset: 0,
    dataSet: [],
    brick_option: {
      showFullContent: true,
      backgroundColor: "rgb(235, 246, 250)",
      forceRepaint: true,
      defaultExpandStatus: false,
      imageFillMode: 'aspectFill',
      columns: 1,
      icon: {
        fill: 'https://images.ifanr.cn/wp-content/uploads/2018/08/640-90-1024x576.jpeg',
        default: 'https://images.ifanr.cn/wp-content/uploads/2018/08/640-90-1024x576.jpeg',
      },
    }
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    if (wx.$param["mode"] != "prod") {
      this.setData({
        normal: false
      })
      return
    } else {
      this.setData({
        normal: true
      })
    }
    
    let id = options.id
    console.log("用户id:", id)
    if (id == undefined) {
      wx.BaaS.auth.getCurrentUser().then(user => {
        console.log(user)
        this.data.uid = user.id
        this.getGoods()
      }).catch(err => {
        if (err.code === 604) {
          console.log('用户未登录')
        }
      })
    } else {
      this.data.uid = Number(id)
      this.getGoods()

      let MyUser = new wx.BaaS.User()
      MyUser.get(this.data.uid).then(user => {

        let title = "Ta的发布"
        if (user.data.nickname) {
          title = user.data.nickname + "的发布"
        }
        this.setData({
          title: title
        })
      })
    }
  },


  onShareAppMessage: function () {

  },

  // 下拉刷新
  onPullDownRefresh: function () {
    this.setData({
      offset: 0, //恢复偏移量
      loadDone: false, //加载完毕
      queryStr: ""
    })
    this.getGoods()
    setTimeout(function () {
      wx.stopPullDownRefresh()
    }, 3000)
  },

  // 点击卡片，获取商品id，转跳详情页面
  tapCard: function (event) {
    console.log("商品ID：", event.detail.card_id)
    wx.navigateTo({
      url: '/pages/Life/oldthings/detail?id=' + event.detail.card_id,
    })
  },

  // 触底加载更多，需改变offset，判断有无更多
  onReachBottom: function () {
    if (this.data.loadDone) return
    console.log('加载更多')
    this.data.offset = this.data.offset + this.data.limit
    this.getGoods(true)
  },

  // 获取商品
  getGoods(loadMore = false) {
    this.setData({
      loading: true
    })
    let table = new wx.BaaS.TableObject('flea_market')

    let query = new wx.BaaS.Query()
    query.compare('created_by', '=', this.data.uid) //筛选条件，0表示正常


    let limit = this.data.limit
    table.setQuery(query)
      .limit(limit).offset(this.data.offset) //一次获取20条数据
      .orderBy('-created_at') //发布时间倒序
      .expand(['created_by']) //拓展用户信息
      .find().then(res => {
        console.log("商品列表", res.data.objects)

        if (loadMore) {
          // 触底加载更多，把获取到的数据拼接到原数据后面
          this.setData({
            dataSet: this.data.dataSet.concat(res.data.objects),
            loading: false,
            loadDone: res.data.objects.length <= limit ? true : false //加载完毕
          })
        } else {
          // 切换类别，直接渲染
          this.setData({
            dataSet: res.data.objects,
            loading: false,
            loadDone: res.data.objects.length < limit ? true : false //加载完毕
          })
        }
        wx.stopPullDownRefresh()
      }, err => {
        wx.showToast({
          title: '获取商品错误',
          icon: "none"
        })
        that.setData({
          loading: false
        })
      })
  },

})