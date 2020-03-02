// components/fake/fake.js
Component({

  options: {
    addGlobalClass: true
  },

  properties: {

    type: {
      type: String,
      value: "default"
    }

  },


  data: {

    article: [{
      src: "https://mp.weixin.qq.com/s/NbVEBpvlPgpfbAf_tRh09w",
      cover: "http://mmbiz.qpic.cn/mmbiz_jpg/9fnicKncw0JPibAPOAA1f5a7kDJ0B5GC6KeurCE4F1WZ32jhS7G2FfqKsfvy586chnSfL8pPNVTD6SZXlSIas3Ow/0?wx_fmt=jpeg",
      desc: "新功能上线！知道你与毕业的距离有多远吗？",
      title: "学业情况",
      date: "2019年12月18日"
    }, {
      src: "https://mp.weixin.qq.com/s/IAh_JOhb3AC8Jfkmxb8rZQ",
      cover: "http://mmbiz.qpic.cn/mmbiz_jpg/hk9FAIGXk22OjSFIjzTibibP2lQkpspqxIibxKfNcSOBbXFWtuXCLNcj2Et4LmFgYc8nNJsF5WQON79GiaicEiaiaOLmA/0?wx_fmt=jpeg",
      desc: "十佳校媒巡展 | 广州大学青年传媒中心",
      title: "小青媒",
      date: "2019年12月12日"
    }, {
      src: "https://mp.weixin.qq.com/s/s01gMimoUzbk9q8bZ8gNDQ",
      cover: "http://mmbiz.qpic.cn/mmbiz_jpg/9fnicKncw0JPgRbe2VzHadlCia0Afkby5BciarZvqjoozGUULwSianMCJ6IicnmluDq6sX7c2l2HmYEJMB8sK25u90A/0?wx_fmt=jpeg",
      desc: "新上功能 | 广大派查询第二课堂学分",
      title: "第二课堂",
      date: "2019年10月20日"
    }, {
      src: "https://mp.weixin.qq.com/s/DwlvT6r3PKUeNRiUDuQXxg",
      cover: "http://mmbiz.qpic.cn/mmbiz_jpg/47d2LKstNMjICJLlCcsZbVYx39VpXtczI5EN0HZtweMkDPhAGOHqk9Pg6et7FTD5evThQ8iaIfuScsLdGrdGkUQ/0?wx_fmt=jpeg",
      desc: "广大学子即将刷屏广州地铁！53321275，31176753！",
      title: "70周年国庆",
      date: "2019年09月30日"
    }]

  },

  methods: {

    onTapAdd() {
      wx.showToast({
        title: '请联系派派！',
        icon: "none"
      })
    },

    onTap(e) {
      console.log(e)
      let url = "/pages/Setting/webview/webview?src=" + e.currentTarget.dataset.src
      wx.$navTo(url)
    }
  },

  // 生命周期方法
  lifetimes: {
    attached: function () {
      switch (this.data.type) {
        case "nav1":
          break
        case "nav1":
          break
        case "nav1":
          break
        default:

      }
    }
  },

  pageLifetimes: {
    show() {
      this.setData({
        fake: wx.$param.mode != "prod" ? true : false //非生产模式启用
      })
    },
  }
})