var utils = require("../../../utils/date.js")
Page({


  // 对应模板id：aFpe_zN27IOKa3I_WhATW4-CxxcsOhwlFJbLJpz1zuk

  data: {

    date: utils.formatTime(new Date(), "y-m-d"),
    time: "12:00",

    pageSize: 10, //每页数量
    page: 1, //页数
    list: [],
    loading: true,

    tplData: {
      keyword1: {
        value: "" //课程名称
      },
      keyword2: {
        value: "" //上课时间
      },
      keyword3: {
        value: "" //上课地点
      }
    },
    sentTime: "",
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {

    let that = this
    setTimeout(() => {
      that.setData({
        loading: false
      })
    }, 500);

    this.getList()
  },

  onShareAppMessage: function () {

  },

  dateChange(e) {
    this.setData({
      date: e.detail.value
    })
  },
  timeChange(e) {
    this.setData({
      time: e.detail.value
    })
  },

  /*
   * 伪双向绑定
   * wxml input 定义属性：data-field="field1.field2" value="{{field1.field2}}"
   * 输入内容将绑定到：this.data.field1.field2 = e.detail.value
   */
  inputBind(e) {
    wx.$bindInput.call(this, e)
  },

  // 触底加载更多，需改变offset，判断有无更多
  onReachBottom: function () {
    if (this.data.loadDone) return
    console.log('加载更多')
    this.data.page = this.data.page + 1
    this.getList(true)
  },

  getList(loadMore = false) {

    let id = wx.getStorageSync('gzhupi_user').id
    if (id == 0 || id == undefined || !id) {
      return
    }

    let query = {
      _page: this.data.page,
      _page_size: this.data.pageSize,
      created_by: id,
      _order: "status,sent_time",
    }
    query = wx.$objectToQuery(query)

    let url = wx.$param.server["prest"] + wx.$param.server["scheme"] + "/t_notify" + query
    // url = "http://192.168.2.214:9000/api/v1/postgres/public/" + "/t_notify" + query
    wx.$ajax({
        url: url,
        method: "get",
      })
      .then(res => {
        console.log("通知列表", res)
        // 格式化时间
        for (let i = 0; i < res.data.length; i++) {
          let time = new Date(res.data[i].sent_time)
          res.data[i].sent_time = utils.formatTime(time, "full")
        }
        if (loadMore) {
          this.data.list = this.data.list.concat(res.data)
        } else {
          this.data.list = res.data
        }
        this.setData({
          list: this.data.list,
          loadDone: res.data.length < this.data.pageSize ? true : false //加载完毕
        })

      }).catch(err => {
        console.log(err)
      })
  },

  operate(e) {

    console.log(e.target)
    let that = this
    switch (e.target.dataset.op) {

      case "delete":
        wx.showModal({
          title: '删除提示',
          content: '确定删除该提醒吗？',
          success(res) {
            if (res.confirm) {
              that.delByPk(e.target.dataset.id)
            }
          }
        })
        break
      case "add":
        this.setData({
          showAdd: true
        })

        break
      default:
        console.warn("unknown op case")
        return
    }
  },

  delByPk(row_id) {
    if (!row_id) return
    wx.$ajax({
        url: wx.$param.server["prest"] + wx.$param.server["scheme"] + "/t_notify?id=$eq." + row_id,
        // url: "http://localhost:9000/api/v1/postgres/public/" + "/t_notify?id=$eq." + row_id,
        method: "delete",
        loading: true,
      })
      .then(res => {
        console.log(res)

        if (res.statusCode == 200 && res.data.rows_affected == 1) {
          wx.showToast({
            title: '删除成功！',
          })

          // 删除数组元素
          for (let i = 0; i < this.data.list.length; i++) {
            if (this.data.list[i].id == row_id) {
              this.data.list.splice(i, 1)
              this.setData({
                list: this.data.list
              })
            }
          }
        } else {
          wx.showModal({
            title: '失败提示',
            content: JSON.stringify(res.data),
          })
        }

      }).catch(err => {
        console.error(err)
      })
  },


  post() {

    let user = wx.getStorageSync('gzhupi_user')
    if (!user || !user.mp_open_id) {
      wx.showModal({
        title: "提示",
        content: "当前用户未绑定公众号，请到广大派公众号回复【绑定】"
      })
      return
    }

    let d = this.data.tplData
    if (d.keyword1.value == "" || d.keyword2.value == "" || d.keyword3 == "") {
      wx.showToast({
        title: '信息不完整',
        icon: "none"
      })
      return
    }

    let sentTime = new Date(this.data.date + " " + this.data.time)

    // if (sentTime.getTime() < new Date().getTime() + 10 * 60 * 1000) {
    //   wx.showToast({
    //     title: '时间无效，至少10分钟后',
    //     icon: "none"
    //   })
    //   return
    // }

    if (sentTime.getTime() < new Date().getTime()) {
      wx.showToast({
        title: '过去的时间无效',
        icon: "none"
      })
      return
    }


    let form = {
      data: this.data.tplData,
      type: "上课提醒",
      sent_time: new Date(this.data.date + " " + this.data.time)
    }

    console.log(form)

    wx.$ajax({
        url: wx.$param.server["prest"] + wx.$param.server["scheme"] + "/t_notify",
        // url: "http://192.168.2.214:9000/api/v1/postgres/public/t_notify",
        data: form,
        loading: true,
        header: {
          "content-type": "application/json"
        }
      })
      .then(res => {
        console.log("添加结果：", res.data)
        if (res.statusCode >= 400) {
          wx.showModal({
            title: '添加失败',
            content: res.errMsg + res.data ? res.data.msg : "",
          })
          return
        }
        if (typeof res.error != "undefined") {
          wx.showModal({
            title: '添加失败',
            content: res.error,
          })
          return
        }
        wx.showToast({
          title: '添加成功',
        })

        this.data.page = 1
        this.getList()

      }).catch(err => {
        console.log(err)
      })

  }

})