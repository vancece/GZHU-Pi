const Page = require('../../../utils/sdk/ald-stat.js').Page;
var Utils = require("../../../utils/utils.js")
Page({

  data: {
    total: 0,
    rooms: [],
    showFilter: false,
    schoolWeek: Utils.getSchoolWeek(), //校历周
    today: new Date().getDay() == 0 ? 7 : new Date().getDay(), //星期几
    times: [{
      text: "全 天",
      checked: true
    }, {
      text: "上 午",
      checked: true
    }, {
      text: "下 午",
      checked: true
    }, {
      text: "晚 上",
      checked: true
    }, ],

    postDate: {
      username: "",
      password: "",
      xqh_id: 1, //校区号 大学城1，桂花岗2
      xnm: 2018, //学年 *
      xqm: 12, //学期 *
      cdlb_id: "", //场地类别 *
      qszws: "", //最小座位数 *
      jszws: "", //最大座位数 *
      lh: "", //楼号 *
      cdmc: "", //场地名称 *
      zcd: Utils.getSchoolWeek(), //周次 +
      xqj: new Date().getDay() == 0 ? 7 : new Date().getDay(), //星期几 +
      jcd: "1,2,3,4,5,6,7,8,9,10,11", //节次 *
      "queryModel.currentPage": 1, //前往页数 *
    }

  },

  onLoad: function(options) {
    this.initWeekDay(this.data.schoolWeek)
    this.initWeek()

    this.getRooms()
  },

  catchtap() {},
  onShareAppMessage: function() {},

  onReachBottom: function() {
    let page = this.data.postDate["queryModel.currentPage"]
    this.data.postDate["queryModel.currentPage"] = page + 1
    if (page + 1 > this.data.total / 30) {
      wx.showToast({
        title: '没有更多啦！',
        icon: "none"
      })
      return
    }
    this.getRooms(true)
  },

  checkAccount() {
    let account = wx.getStorageSync("account")
    if (account == "") {
      wx.showToast({
        title: '未绑定学号',
        icon: "none"
      })
      wx.navigateTo({
        url: 'pages/Setting/login/bindStudent',
      })
      return 0
    } else {
      this.data.postDate.username = account.username
      this.data.postDate.password = account.password
      return 1
    }
  },

  formSubmit(e) {
    console.log(e.detail.value.query)
    if (e.detail.value.query == "") {
      wx.showToast({
        title: '默认查询当天全部',
        icon: "none"
      })
    }
    this.data.postDate.cdmc = e.detail.value.query
    this.getRooms()
  },


  getRooms(load = false) {
    var time = new Date()
    if (time.getHours() >= 0 && time.getHours() < 7) {
      wx.showToast({
        title: '当前时间段不可用~',
        icon: "none"
      })
      return
    }
    if (!this.checkAccount()) return

    let that = this
    wx.showLoading({
      title: '加载中...',
      icon: "none"
    })

    wx.request({
      // url:"http://127.0.0.1:5000/room",
      url: 'https://1171058535813521.cn-shanghai.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/Spider/room',
      data: this.data.postDate,
      method: "post",
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      success(res) {
        if (res.statusCode != 200) {
          wx.showToast({
            title: '请求超时',
            icon: "none"
          })
          return
        }
        if (res.data.statusCode != 200) {
          wx.showToast({
            title: '账号或密码错误',
            icon: "none"
          })
          return
        }
        wx.showToast({
          title: '找到 ' + res.data.data.total + " 间教室",
          icon: "none"
        })

        let rooms = res.data.data.rooms
        if (load) { //加载更多，拼接
          rooms = that.data.rooms.concat(res.data.data.rooms)
        }
        that.setData({
          rooms: rooms,
          total: res.data.data.total,
        })
      },
      fail: function(err) {
        wx.showToast({
          title: '服务器响应错误',
          icon: "none"
        })
      },
      complete: function(res) {
        wx.hideLoading()
      }
    })
  },

  // 打开弹窗
  openFilter() {
    this.setData({
      showFilter: true
    })
  },

  // 弹窗确定
  confirm() {
    this.data.postDate["queryModel.currentPage"] = 1
    this.form()
    this.getRooms()
  },

  // 切换校区
  switchZone() {
    let that = this
    wx.showActionSheet({
      itemList: ["大学城校区", "桂花岗校区"],
      success(res) {
        that.setData({
          tapIndex: res.tapIndex
        })
        that.data.postDate.xqh_id = res.tapIndex + 1
      }
    })
  },

  // 生成表单
  form() {
    // 周
    let weeks = []
    for (let i = 0; i < this.data.weeks.length; i++) {
      if (this.data.weeks[i].checked) weeks.push(i)
    }
    // 星期
    let weekdays = []
    for (let i = 0; i < this.data.weekdays.length; i++) {
      if (this.data.weekdays[i].checked) weekdays.push(i)
    }
    // 时间段
    let times = []
    let str = ["1,2,3,4,5,6,7,8,9,10,11", "1,2,3,4", "5,6,7,8", "9,10,11"]
    for (let i = 0; i < this.data.times.length; i++) {
      if (this.data.times[0].checked) {
        times = str[0]
        break
      } else {
        if (this.data.times[i].checked) times.push(str[i])
      }
    }
    // 转换成字符串
    this.data.postDate.zcd = weeks.length == 0 ? Utils.getSchoolWeek() : weeks.toString()
    this.data.postDate.xqj = weekdays.length == 0 ? this.data.today : weekdays.toString()
    this.data.postDate.jcd = times.length == 0 ? str[0] : times.toString()
  },


  // 多选，处理各种逻辑
  check(e) {
    let id = e.currentTarget.id
    let index = Number(e.currentTarget.dataset.index)

    if (id == "week") {
      let weeks = this.data.weeks
      if (weeks[index].forbid) return
      weeks[index].checked = !weeks[index].checked
      if (!weeks[this.data.schoolWeek].checked) {
        this.initWeekDay(this.data.schoolWeek + 1)
      } else {
        this.initWeekDay(this.data.schoolWeek)
      }
      this.setData({
        weeks: weeks
      })
      return
    }

    if (id == "weekday") {
      let weekdays = this.data.weekdays
      if (weekdays[index].forbid) return
      weekdays[index].checked = !weekdays[index].checked
      this.setData({
        weekdays: weekdays
      })
      return
    }

    if (id == "time") {
      let times = this.data.times
      if (index == 0) {
        let flag = times[0].checked
        times.forEach(item => {
          item.checked = !flag
        })
      } else {
        times[index].checked = !times[index].checked
        if (times[1].checked && times[2].checked && times[3].checked) {
          times[0].checked = true
        } else {
          times[0].checked = false
        }
      }
      this.setData({
        times: times
      })
      return
    }
  },

  // 初始周数，禁止过去的时间
  initWeek() {
    let weeks = []
    for (let i = 0; i < 26; i++) {
      let item = {
        checked: false,
        forbid: false
      }
      if (i < this.data.schoolWeek) item.forbid = true
      if (i == this.data.schoolWeek) item.checked = true
      weeks.push(item)
    }
    this.setData({
      weeks: weeks
    })
  },

  // 初始星期，禁止过去的时间
  initWeekDay(week = 0) {
    let str = ["整 周", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六", "星期日"]
    let weekdays = []

    for (let i = 0; i < 8; i++) {
      let item = {
        text: str[i],
        checked: false,
        forbid: false
      }
      if (week == this.data.schoolWeek) {
        item.checked = this.data.today == i ? true : false
        item.forbid = this.data.today > i ? true : false
      }
      weekdays.push(item)
    }
    this.setData({
      weekdays: weekdays
    })
  },

})