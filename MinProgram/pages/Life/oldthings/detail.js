const Page = require('../../../utils/sdk/ald-stat.js').Page;
import dateUtil from '../../../utils/date.js'
let app = getApp()
let tableName = "flea_market"

Page({

  data: {
    mode: wx.$param["mode"],
    bindStatus: app.globalData.bindStatus,
    loading: true,
    isOwner: false, //物品发布者
    refreshable: false, //是否可以刷新
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

  onLoad: async function (options) {

    if (this.fake()) return
    
    if (options.id == undefined || options.id == "") {
      wx.showModal({
        title: '提示',
        content: '页面不存在...',
        success() {
          wx.$navTo("/pages/Campus/home/home")
        }
      })
      return
    }



    // 获取当前用户
    await wx.BaaS.auth.getCurrentUser().then(user => {
      this.data.uid = user.id
    })
    let id = options.id

    this.getDetail(id)
    this.updateCounter(id)
    this.setData({
      id: id
    })
  },

  onShow(options) {
    let that = this
    setTimeout(function () {
      that.setData({
        bindStatus: app.globalData.bindStatus
      })
    }, 1000)
  },

  // 点击头像
  tapUser: function () {
    if (!this.data.detail.created_by.id) return
    wx.navigateTo({
      url: '/pages/Life/oldthings/mine?id=' + this.data.detail.created_by.id,
    })
  },

  // 获取单个商品全部信息
  getDetail(id) {
    let that = this
    let table = new wx.BaaS.TableObject(tableName)

    table.expand(['created_by']).get(id).then(res => {
      console.log("商品信息：", res.data)
      res.data.created_at = dateUtil.relativeTime(res.data.created_at)
      res.data.updates_at = dateUtil.relativeTime(res.data.updates_at)

      // 判断当前用户是否发布者
      let isOwnwe = false
      if (this.data.uid == res.data.created_by.id) isOwnwe = true

      // 判断是否擦亮，每天可以擦亮一次
      let refreshable1 = Date.parse(new Date()) / 1000 - Date.parse(new Date(res.data.refresh_time * 1000)) / 1000 > 24 * 60 * 60
      let refreshable2 = (Math.abs(new Date().getDate() - new Date(res.data.refresh_time * 1000).getDate()) > 0)
      let refreshable = refreshable1 || refreshable2

      let shareModal = res.data.viewed == 0 ? true : false
      that.setData({
        detail: res.data,
        loading: false,
        shareModal: shareModal,
        isOwner: isOwnwe,
        refreshable: refreshable
      })
    }, err => {
      wx.showModal({
        title: '提示',
        content: '该商品不存在',
        success(res) {
          wx.redirectTo({
            url: '/pages/Life/oldthings/index',
          })
        }
      })
      that.setData({
        loading: false
      })
    })
  },

  // 原子性更新计数器
  updateCounter(id) {
    let table = new wx.BaaS.TableObject(tableName)
    let record = table.getWithoutData(id)
    record.incrementBy('viewed', 1)
    record.update()
  },

  // 复制内容到剪贴板
  onCopy(e) {
    wx.setClipboardData({
      data: e.target.dataset.copy,
    })
  },

  // 页面转跳
  navTo(e) {
    console.log("转跳", e.target.id)
    switch (e.target.id) {
      case "bind":
        wx.navigateTo({
          url: '/pages/Setting/login/bindStudent',
        })
    }
  },

  viewImage(e) {
    wx.previewImage({
      urls: this.data.detail.image,
      current: e.currentTarget.dataset.url
    });
  },

  // 分享页面带上商品id
  onShareAppMessage: function () {
    this.setData({
      shareModal: false
    })
    return {
      title: '校园二手:' + this.data.detail.title,
      desc: '',
      path: '/pages/Life/oldthings/detail?id=' + this.data.id,
      imageUrl: "",
      success: function (res) {
        wx.showToast({
          title: '分享成功',
          icon: "none"
        });
      }
    }
  },

  // ==============管理=============

  manage(e) {
    let that = this
    switch (e.target.dataset.op) {
      case "擦亮":
        if (!this.data.refreshable) {
          wx.showToast({
            title: '明天再来😯',
            icon: "none"
          })
          return
        }
        // 刷新时间，秒级时间戳
        this.update("refresh_time", Date.parse(new Date()) / 1000)
        wx.showToast({
          title: '今日擦亮成功😝',
          icon: "none"
        })
        break
      case "上架":
        this.update("status", 0)
        break
      case "下架":
        this.update("status", -1)
        break
      case "改价":
        this.setData({
          changePrice: true
        })
        break
      case "删除":
        wx.showModal({
          title: '删除提示',
          content: '确定删除该二手物品吗？',
          success(res) {
            if (res.confirm) {
              if (that.data.detail.info.file_ids) {
                that.delFile(that.data.detail.info.file_ids)
              }
              that.delGoods(that.data.id)
            }
          }
        })
        break
      default:
        return
    }
  },

  changePrice() {
    this.update("price", this.data.newPrice)
  },

  priceInput(e) {
    this.data.newPrice = Number(e.detail.value)
  },

  delFile(fileIDs = []) {
    if (!fileIDs) return
    let MyFile = new wx.BaaS.File()
    MyFile.delete(fileIDs).then()
  },

  delGoods(recordID) {
    if (!recordID) return
    let Product = new wx.BaaS.TableObject(tableName)
    Product.delete(recordID).then(res => {
      wx.showToast({
        title: '删除成功！',
      })
      setTimeout(function () {
        wx.redirectTo({
          url: '/pages/Life/oldthings/index',
        })
      }, 1000)
    }, err => {
      wx.showToast({
        title: '删除失败',
        icon: 'none',
      })
    })
  },

  update(key, value) {
    if (!this.data.id || !key) return

    let MyTableObject = new wx.BaaS.TableObject(tableName)
    let product = MyTableObject.getWithoutData(this.data.id)
    product.set(key, value)
    product.update().then(res => {
      // 判断是否擦亮，每天可以擦亮一次
      let refreshable1 = (Date.parse(new Date()) / 1000 - Date.parse(new Date(res.data.refresh_time * 1000)) / 1000) > 24 * 60 * 60
      let refreshable2 = (Math.abs(new Date().getDate() - new Date(res.data.refresh_time * 1000).getDate()) > 0)
      let refreshable = refreshable1 || refreshable2

      this.setData({
        'detail.status': res.data.status,
        'detail.price': res.data.price,
        refreshable: refreshable
      })
    }, err => {
      console.log(err)
    })
  },


  // ===============管理员检查============

  // 判断当前用户是否管理员
  checkAdmin() {
    let MyTableObject = new wx.BaaS.TableObject("config")
    MyTableObject.get("5d712ad6db51692484017b6d").then(res => {
      if (res.data.data.id && res.data.data.id.indexOf(this.data.uid) != -1) {
        this.setData({
          isOwner: true,
          refreshable: true
        })
      }
    })
  }

})