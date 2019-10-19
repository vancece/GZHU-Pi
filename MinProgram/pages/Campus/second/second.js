const Page = require('../../../utils/sdk/ald-stat.js').Page
import Data from './data.js';
var baseUrl = "https://gzhu.ifeel.vip"

Page({

  data: {
    myApply: wx.getStorageSync("myApply"), //我的申报
    applyList: wx.getStorageSync("myApply"), //申报项目列表
    showFilter: false,
    showDetail: false,
    loading: false,
    type: ["“三创”能力教育类", "美育体育教育类", "思想成长教育类", "实践公益教育类"],

    yearList: ["全部", "2020-2021", "2019-2020", "2018-2019", "2017-2018"],
    yearIndex: 2, //默认显示的学年索引

    gradeList: ["全部", "2020", "2019", "2018", "2017"],
    gradeIndex: 3, //默认显示的学年索引

    statusList: Data.statusList,
    statusIndex: 0, //默认显示的学年索引
    statusMap: Data.statusMap,

    collegeList: Data.collegeList,
    collegeIndex: 12, //默认显示的学年索引
    collegeMap: Data.collegeMap,
    queryStr: "",

    gridCol: 2,
    iconList: [{
      icon: 'https://cos.ifeel.vip/gzhu-pi/images/icon/tongji.svg',
      name: '学分统计'
    },  {
      icon: 'https://cos.ifeel.vip/gzhu-pi/images/icon/shaixuan.svg',
      name: '筛选'
    }],

    filter: {
      username: "",
      password: "",
      year: "", //学年
      status_no: "", //审核状态，匹配statusMap
      grade: "", //年级
      college_no: "", //学院编号，匹配collegeMap
      page: 1, //每次其它参数变动，page都要还原为1
      stu_name: "", //学生姓名
      item_name: "" //项目名称
    }
  },

  catchtap() {},
  onLoad: function(options) {
    let user = wx.getStorageSync("account")
    if (!user) {
      wx.showToast({
        title: '未绑定学号',
        icon: "none"
      })
      return
    }
    this.getMyData()
  },

  onShareAppMessage: function() {},

  // 加载更多
  onReachBottom: function() {

    if (this.data.myApply == this.data.applyList) return
    if (this.data.applyList.length % 10 > 0) return

    let user = wx.getStorageSync("account")
    let form = this.data.filter
    form["username"] = user.username
    form["password"] = user.password
    form["page"] = form["page"] + 1

    this.getSearchData(form, true)
  },

  search(e) {
    if (this.data.queryStr == "") return
    let user = wx.getStorageSync("account")
    if (!user) {
      wx.showToast({
        title: '未绑定学号',
        icon: "none"
      })
      return
    }
    let form = {}
    form["username"] = user.username
    form["password"] = user.password

    form["page"] = 1
    if (e.currentTarget.id == "stu_name") {
      form["stu_name"] = this.data.queryStr
    }
    if (e.currentTarget.id == "item_name") {
      form["item_name"] = this.data.queryStr
    }
    this.data.filter = form
    this.getSearchData(form)

    this.setData({
      queryStr: this.data.queryStr
    })
  },

  searchInput: function(e) {
    this.data.queryStr = e.detail.value
  },

  // 筛选确认
  filterConfirm() {
    let user = wx.getStorageSync("account")
    if (!user) {
      wx.showToast({
        title: '未绑定学号',
        icon: "none"
      })
      return
    }
    let form = this.data.filter
    form["page"] = 1
    form["username"] = user.username
    form["password"] = user.password

    form["year"] = this.data.yearList[this.data.yearIndex]
    if (form["year"] == "全部") form["year"] = ""
    form["grade"] = this.data.gradeList[this.data.gradeIndex]
    if (form["grade"] == "全部") form["grade"] = ""

    let status = this.data.statusList[this.data.statusIndex]
    form["status_no"] = this.data.statusMap[status]

    let college = this.data.collegeList[this.data.collegeIndex]
    form["college_no"] = this.data.collegeMap[college]

    form["item_name"] = this.data.queryStr

    this.getSearchData(form)

    this.setData({
      queryStr: this.data.queryStr
    })
  },

  // 点击工具栏
  tapTools(e) {
    let id = Number(e.currentTarget.id)
    let name = this.data.iconList[id].name
    switch (name) {
      case "学分统计":
        this.count()
        break
      case "筛选":
        this.setData({
          showFilter: true
        })
        break
    }
  },

  // 列表选项
  actionSheet(e) {
    console.log(e)
    let id = e.currentTarget.id
    let that = this
    switch (id) {
      case "year":
        wx.showActionSheet({
          itemList: this.data.yearList,
          success(res) {
            that.setData({
              yearIndex: res.tapIndex
            })
          }
        })
        break
      case "grade":
        wx.showActionSheet({
          itemList: this.data.gradeList,
          success(res) {
            that.setData({
              gradeIndex: res.tapIndex
            })
          }
        })
        break
      case "status":
        that.setData({
          statusIndex: Number(e.detail.value)
        })
        break
      case "college":
        that.setData({
          collegeIndex: Number(e.detail.value)
        })
        break
    }
  },

  viewImage(e) {
    wx.previewImage({
      urls: this.data.detail.images,
      current: e.currentTarget.dataset.url
    });
  },

  // 点击列表
  tapItem(e) {
    let id = Number(e.currentTarget.id)
    this.setData({
      detail: this.data.applyList[id],
      showDetail: true
    })
    this.getImage(this.data.applyList[id].id)
  },

  getSearchData(form, loadMore = false) {
    this.setData({
      loading: true
    })
    let that = this
    wx.request({
      url: baseUrl + "/api/v1/second/search",
      method: "POST",
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      data: form,
      success: function(res) {
        if (res.data.status != 200) {
          wx.showToast({
            title: res.data.msg,
            icon: "none"
          })
          return
        }
        if (res.data.data.length == 0) {
          wx.showToast({
            title: '没有更多啦',
            icon: "none"
          })
          return
        }
        if (loadMore) {
          that.setData({
            applyList: that.data.applyList.concat(res.data.data),
          })
        } else {
          that.setData({
            applyList: res.data.data,
          })
        }

      },
      fail: function(err) {
        wx.showModal({
          title: '请求失败',
          content: "错误信息:" + err.errMsg,
        })
      },
      complete(res) {
        console.log(res.data)
        that.setData({
          loading: false
        })
      }
    })
  },

  getMyData() {
    this.setData({
      loading: true
    })
    let that = this
    wx.request({
      url: baseUrl + "/api/v1/second/my",
      method: "POST",
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      data: wx.getStorageSync("account"),
      success: function(res) {
        if (res.data.status != 200) {
          wx.showToast({
            title: res.data.msg,
            icon: "none"
          })
          return
        }
        wx.setStorageSync("myApply", res.data.data)
        that.setData({
          applyList: res.data.data,
          myApply: res.data.data
        })
      },
      fail: function(err) {
        wx.showModal({
          title: '请求失败',
          content: "错误信息:" + err.errMsg,
        })
      },
      complete(res) {
        console.log(res.data)
        that.setData({
          loading: false
        })
      }
    })
  },

  // 获取证明材料
  getImage(id) {
    let that = this
    let data = wx.getStorageSync("account")
    data["id"] = id
    wx.request({
      url: baseUrl + "/api/v1/second/image",
      method: "POST",
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      data: data,
      success: function(res) {
        let detail = that.data.detail
        detail["images"] = res.data.data
        that.setData({
          detail: detail
        })
      },
      complete(res) {
        console.log(res.data)
      }
    })
  },

  // 学分统计
  count() {
    let countData = {
      all: 0, //全部
      ability: 0, //三创能力类
      arts: 0, //文体艺术类
      thought: 0, //思想成长类
      practice: 0, //实践公益类

      refuse: 0, //不通过
      unaudited: 0, //未审核
    }
    for (let i = 0; i < this.data.myApply.length; i++) {
      let item = this.data.myApply[i]
      if (item.audit_mark.indexOf("不通过") != -1) {
        countData.refuse = countData.refuse + item.apply_credit
        continue
      }
      if (item.audit_mark.indexOf("未审核") != -1) {
        countData.unaudited = countData.unaudited + item.apply_credit
        continue
      }
      if (item.audit_credit == 0) continue
      if (item.audit_mark.indexOf("审核通过") == -1) continue

      countData.all = countData.all + item.apply_credit
      switch (item.type) {
        case '“三创”能力教育类':
          countData.ability = countData.ability + item.audit_credit
          break
        case '美育体育教育类':
          countData.arts = countData.arts + item.audit_credit
          console.log(item.name)
          break
        case '思想成长教育类':
          countData.thought = countData.thought + item.audit_credit
          break
        case '实践公益教育类':
          countData.practice = countData.practice + item.audit_credit
          break
      }
    }
    countData.all = Math.round(countData.all * 100) / 100
    countData.ability = Math.round(countData.ability * 100) / 100
    countData.arts = Math.round(countData.arts * 100) / 100
    countData.thought = Math.round(countData.thought * 100) / 100
    countData.practice = Math.round(countData.practice * 100) / 100
    countData.refuse = Math.round(countData.refuse * 100) / 100
    countData.unaudited = Math.round(countData.unaudited * 100) / 100

    this.setData({
      showCount: true,
      countData: countData
    })
  },
})