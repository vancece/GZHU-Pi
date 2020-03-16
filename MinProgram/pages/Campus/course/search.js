const Page = require('../../../utils/sdk/ald-stat.js').Page;

Page({
  data: {
    weekday: new Date().getDay(), //星期几
    showTarget: 0, //默认展开第一条
    page: 0,
    courses: [],
    queryType: "onload", //查询类型，默认同启动时

    zoneList: ["大学城", "桂花岗"],
    zoneIndex: 0, //默认校区
    semList: wx.$param.school["sem_list"],
    semIndex: wx.$param.school["sem_list"].indexOf(wx.$param.school["year_sem"]), //默认显示的学期索引
    typeList: ["专业必修课程", "专业选修课程", "通识类必修课程", "通识类选修课程", "学科基础课程", "教师教育类必修课程"],
    typeIndex: 0, //默认显示的课程类型
    //条件筛选内容
    fliter_content: {
      zone: "大学城",
      sem: wx.$param.school["year_sem"],
      type: "专业必修课程"
    },
  },

  onLoad: function(options) {
    this.getCourse("通识类选修课", true)
  },

  onShareAppMessage: function() {

  },
  catchtap() {},

  nav(e) {
    let id = e.currentTarget.id
    wx.navigateTo({
      url: '/pages/Campus/course/tools?id=' + id,
    })
  },

  // 加载更多
  onReachBottom: function() {
    this.data.page = this.data.page + 1
    let offset = this.data.page * 20
    if (offset >= this.data.total) {
      wx.showToast({
        title: '加载完啦！',
        icon: "none"
      })
      this.setData({
        loadAll: true
      })
      return
    }else{
      this.setData({
        loadAll: false
      })
    }
    if (this.data.queryType == "vague") this.getCourse(this.data.query, false, offset)
    if (this.data.queryType == "fliter") this.getCourseFliter(offset)
    if (this.data.queryType == "onload") this.getCourse("通识类选修课", true, offset)
  },

  //点击展开折叠
  flod(e) {
    let showTarget = e.currentTarget.id
    if (showTarget == this.data.showTarget) showTarget = -1
    this.setData({
      showTarget: showTarget
    })
  },

  // 输入框搜索
  formSubmit(e) {
    wx.BaaS.wxReportTicket(e.detail.formId)
    if (e.detail.value.query == "") {
      wx.showToast({
        title: '请输入查询内容',
        icon: "none"
      })
      return
    }
    this.data.page = 0
    this.data.queryType = "vague"
    this.data.query = e.detail.value.query
    this.data.courses = []

    this.getCourse(e.detail.value.query)
  },

  // 打开弹窗
  openFilter() {
    this.setData({
      showFilter: true
    })
  },

  // 弹窗确定
  confirm() {
    this.data.page = 0
    this.data.queryType = "fliter"
    this.data.courses = []

    this.getCourseFliter()

    this.setData({
      query: ""
    })

  },

  // 切换选项
  actionSheet(e) {
    let id = e.currentTarget.id
    let that = this
    switch (id) {
      case "zone":
        wx.showActionSheet({
          itemList: this.data.zoneList,
          success(res) {
            that.setData({
              zoneIndex: res.tapIndex
            })
            that.data.fliter_content[id] = that.data.zoneList[res.tapIndex]
          }
        })
        break
      case "sem":
        wx.showActionSheet({
          itemList: this.data.semList,
          success(res) {
            that.setData({
              semIndex: res.tapIndex
            })
            that.data.fliter_content[id] = that.data.semList[res.tapIndex]
          }
        })
        break
      case "type":
        wx.showActionSheet({
          itemList: this.data.typeList,
          success(res) {
            that.setData({
              typeIndex: res.tapIndex
            })
            that.data.fliter_content[id] = that.data.typeList[res.tapIndex]
          }
        })
        break
    }
  },

  getInput(e) {
    let id = e.currentTarget.id
    this.data.fliter_content[id] = e.detail.value
  },


  // 模糊查询（或）
  getCourse(queryStr = "通识类选修课", dayMode = false, offset = 0) {
    this.setData({
      loading: true
    })

    let that = this
    let Obj = new wx.BaaS.TableObject("all_course")

    let fliter = this.data.fliter_content
    let sp = fliter["sem"].split("-")
    
    let query = new wx.BaaS.Query()
    query.compare('xn', '=', sp[0] + "-" + sp[1]) //学年
    query.compare('xq', '=', sp[2]) //学期
    if (dayMode) {
      query.compare('xqj', '=', new Date().getDay()) //星期几
    }

    // 课程名称
    let query1 = new wx.BaaS.Query()
    query1.contains('kcmc', queryStr)
    // 班级组成
    let query2 = new wx.BaaS.Query()
    query2.contains('jxbzc', queryStr)
    // 专业组成
    let query3 = new wx.BaaS.Query()
    query3.contains('zyzc', queryStr)
    // 教师姓名
    let query4 = new wx.BaaS.Query()
    query4.contains('xm', queryStr)
    // 开课学院
    let query5 = new wx.BaaS.Query()
    query5.contains('kkxy', queryStr)
    // 上课地点
    let query6 = new wx.BaaS.Query()
    query6.contains('jxdd', queryStr)
    // 课程性质
    let query7 = new wx.BaaS.Query()
    query7.contains('kcxzmc', queryStr)

    let orQuery = wx.BaaS.Query.or(query1, query2, query3, query4, query5, query6, query7)
    let andQuery = wx.BaaS.Query.and(query, orQuery)

    Obj.setQuery(andQuery).offset(offset).find().then(res => {
      that.setData({
        courses: that.data.courses.concat(res.data.objects)
      })
    })
    // 获取总条数
    Obj.setQuery(andQuery).count().then(num => {
      that.setData({
        total: num,
        loading: false
      })
    }, err => {
      that.setData({
        loading: false
      })
    })
  },

  // 条件查询（与）
  getCourseFliter(offset = 0) {
    this.setData({
      loading: true
    })

    let that = this
    let fliter = this.data.fliter_content
    let sp = fliter["sem"].split("-")

    let Obj = new wx.BaaS.TableObject("all_course")
    let query = new wx.BaaS.Query()

    query.compare('xqmc', '=', fliter["zone"]) //校区
    query.compare('xn', '=', sp[0] + "-" + sp[1]) //学年
    query.compare('xq', '=', sp[2]) //学期
    query.compare('kcxzmc', '=', fliter["type"]) //课程类型
    if (fliter["weekday"]) query.compare('xqj', '=', Number(fliter["weekday"])) //星期
    if (fliter["college"]) query.contains('kkxy', fliter["college"]) //学院
    if (fliter["class"]) query.contains('jxbzc', fliter["class"]) //教学班
    if (fliter["weeks"]) query.contains('qsjsz', fliter["weeks"]) //周段
    if (fliter["section"]) query.contains('skjc', fliter["section"]) //节次

    Obj.setQuery(query).offset(offset).find().then(res => {
      that.setData({
        courses: that.data.courses.concat(res.data.objects)
      })
    })
    // 获取总条数
    Obj.setQuery(query).count().then(num => {
      that.setData({
        total: num,
        loading: false
      })
    }, err => {
      that.setData({
        loading: false
      })
    })
  }

})