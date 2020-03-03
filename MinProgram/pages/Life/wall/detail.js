const Page = require('../../../utils/sdk/ald-stat.js').Page
var utils = require("../../../utils/date.js")
Page({

  data: {
    mode: wx.$param["mode"],
    detail: {},
    navTitle: "广大墙",
    claim_list: [],
  },

  fake() {
    let mode = wx.$param["mode"]
    this.setData({
      mode: mode
    })
    if (mode == "prod") {
      return false
    } else return true
  },

  onLoad: function (options) {

    if (this.fake()) return

    if (!options.id) options.id = 4
    this.data.id = options.id
    this.getDetail(options.id)
  },

  onShow: function () {
    let that = this
    setTimeout(function () {
      that.setData({
        wait: true,
        uid: wx.getStorageSync('gzhupi_user').id
      })
    }, 500)
  },

  onShareAppMessage: function () {
    return {
      title: this.data.navTitle + ": " + this.data.detail.title,
      desc: '',
      path: '/pages/Life/wall/detail?id=' + this.data.id,
      imageUrl: this.data.detail.image[0],
      success: function (res) {
        wx.showToast({
          title: '分享成功',
          icon: "none"
        })
      }
    }
  },

  // true 说明在防抖期间，应该停止执行
  isDebounce(timeout = 2000) {
    let that = this
    if (this.data.debounce) {
      console.log("触发防抖")
      return true
    }
    this.data.debounce = true
    setTimeout(() => {
      that.data.debounce = false
    }, timeout)
    return false
  },

  setTitle(type = "广大墙") {
    switch (type) {
      case "情墙":
        type = "广大情墙"
        break
      case "悄悄话":
      case "树洞":
        type = "悄悄话-心情树洞"
        break
      case "日常":
        type = "广大日常"
        break
      default:
        type = "广大墙"
    }
    this.setData({
      navTitle: type
    })
  },

  // 评论成功回调，发送通知
  discussSuccess(e) {
    let cur_uid = wx.getStorageSync('gzhupi_user').id
    if (cur_uid == this.data.detail.created_by) return

    let sender = wx.getStorageSync('gzhupi_user').nickname
    if (e.detail.anonymity) {
      sender = e.detail.anonymity //使用匿名信息
    }
    if (!sender) return

    wx.cloud.callFunction({
      // 需调用的云函数名
      name: 'sendMsg',
      // 传给云函数的参数
      data: {
        touser: this.data.detail.open_id,
        page: '/pages/Life/wall/detail?id=' + this.data.id,
        content: e.detail.content,
        title: this.data.detail.title,
        type: "comment",
        sender: sender,
      },
      complete: function (res) {
        console.log(res.result)
        if (res.result.errCode == 43101) {
          console.log("该用户未订阅通知")
        }
      }
    })
  },


  sendClaimMsg() {
    let sender = wx.getStorageSync('gzhupi_user').nickname
    if (sender == "" || sender == undefined) return

    wx.cloud.callFunction({
      // 需调用的云函数名
      name: 'sendMsg',
      // 传给云函数的参数
      data: {
        touser: this.data.detail.open_id,
        page: '/pages/Life/wall/detail?id=' + this.data.id,
        content: "有人领取了你的表白，点击查看",
        type: "unread",
        sender: sender,
      },
      complete: function (res) {
        console.log(res.result)
      }
    })
  },


  operate(e) {
    if (this.isDebounce()) return

    console.log(e.target)
    let that = this
    switch (e.target.dataset.op) {
      case "star":
        this.setRelation("star")
        break
      case "claim":
        wx.$subscribe()
        this.setRelation("claim")
        break
      case "delete":
        wx.showModal({
          title: '删除提示',
          content: '确定删除该主题吗？',
          success(res) {
            if (res.confirm) {
              if (that.data.detail.addi.file_ids) {
                that.delFile(that.data.detail.addi.file_ids)
              }
              that.delByPk(that.data.id)
            }
          }
        })
        break
      default:
        console.warn("unknown op case")
        return
    }
  },

  viewImage(e) {
    wx.previewImage({
      urls: this.data.detail.image,
      current: e.currentTarget.dataset.url
    });
  },

  // 点赞、认领  type="star/claim"
  setRelation: function (type = "star", object_id = this.data.id, object = "t_topic") {
    if (!type || !object_id || !object) return

    let cur_uid = wx.getStorageSync('gzhupi_user').id
    let list = this.data.detail[type + "_list"]

    for (let i in list) {
      if (list[i].created_by == cur_uid) {
        wx.$ajax({
            url: wx.$param.server["prest"] + wx.$param.server["scheme"] +"/t_relation?id=$eq." + list[i].id,
            method: "delete",
          })
          .then(res => {
            list.splice(i, 1)
            this.setData({
              ["detail." + type + "_list"]: list
            })
          }).catch(err => {
            console.error(err)
          })
        return
      }
    }
    console.log("用户未点赞/认领，创建")
    wx.$ajax({
      url: wx.$param.server["prest"] + wx.$param.server["scheme"] +"/t_relation",
      method: "post",
      data: {
        object_id: Number(object_id),
        object: object,
        type: type
      },
      header: {
        "content-type": "application/json"
      }
    }).then(res => {
      if (typeof list != "object" || list == null) list = []
      if (res.data && res.data.id) {
        res.data.avatar = wx.getStorageSync('gzhupi_user').avatar
        list.push(res.data)
        this.setData({
          ["detail." + type + "_list"]: list
        })
        // 认领成功发送通知
        if (type == "claim" && cur_uid != this.data.detail.created_by) {
          this.sendClaimMsg()
        }
      }
    })
  },

  // 获取详情
  getDetail(id) {
    wx.$ajax({
        url: wx.$param.server["prest"] + wx.$param.server["scheme"] +"/v_topic?id=$eq." + id,
        method: "get",
        loading: true,
      })
      .then(res => {
        if (res.data.length == 0) {
          wx.showModal({
            title: '提示',
            content: '该主题不存在',
            success(res) {
              wx.$navTo("/pages/Life/wall/wall")
            }
          })
          return
        }
        for (let i = 0; i < res.data.length; i++) {
          let time = new Date(res.data[i].created_at)
          res.data[i].created_at = utils.relativeTime(time.getTime() / 1000)
        }
        this.setData({
          detail: res.data[0]
        })
        this.setTitle(res.data[0].type)
      }).catch(err => {
        console.log(err)
      })
  },

  getRelations(id, type) {
    if (!id || !type) return
    wx.$ajax({
        url: wx.$param.server["prest"] + wx.$param.server["scheme"] +"/v_relation?object_id=$eq." + id + "&type=$eq." + type,
        method: "get",
      })
      .then(res => {
        console.log(res)
        if (type == "claim") {
          this.setData({
            claim_list: res.data
          })
        }
        if (type == "star") {
          this.setData({
            star_list: res.data
          })
        }
      }).catch(err => {
        console.log(err)
      })
  },

  delByPk(row_id) {
    if (!row_id) return
    wx.$ajax({
        url: wx.$param.server["prest"] + wx.$param.server["scheme"] +"/t_topic?id=$eq." + row_id,
        method: "delete",
        loading: true,
      })
      .then(res => {
        console.log(res)

        if (res.statusCode == 200 && res.data.rows_affected == 1) {
          wx.showToast({
            title: '删除成功！',
          })
          setTimeout(function () {
            wx.$navTo("/pages/Life/wall/wall")
          }, 1000)
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

  delFile(fileIDs = []) {
    if (!fileIDs) return
    let MyFile = new wx.BaaS.File()
    MyFile.delete(fileIDs).then(res => {
      console.log("delFile ", res)
    })
  },

})