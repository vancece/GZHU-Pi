const App = require('../sdk/ald-stat.js').App
import Poster from '../com/wxcanvas/poster/poster';
import PostConfig from './postconfig.js';
import Utils from '../utils.js';
Utils.initSdk();
let that;

Page({

  data: {
    showChoice: false,
    // 屏幕比例
    ratio: wx.getSystemInfoSync().screenHeight / wx.getSystemInfoSync().screenWidth,
    //单位高度
    hw: wx.getSystemInfoSync().screenHeight / 100,
    curIndex: 0,
    isAuth: false, //是否授权成功
    userInfo: {},
    forDrawAvatar: "", //用于画图的图片地址，本地/头像
    avatarPath: '', //合成后的头像路径
    postPath: "", //合成后的海报路径

    config: {}, //全局海报配置
    // 默认相框
    frame: "https://cos.ifeel.vip/gzhu-pi/images/campus/frame01.png",
    // 相框列表
    frames: ["https://cos.ifeel.vip/gzhu-pi/images/campus/frame01.png"],
    // 页面背景
    background: "https://cos.ifeel.vip/gzhu-pi/images/campus/bg16.png",

    showCount: false,
    count: 0, //使用统计
    avatarHeight: 42, //头像margin-top单位高度
    countHeight: 11, //统计文字距离底部的高度
  },

  onLoad: function(options) {
    that = this
    this.getAuthStatus()

    if (options.id == undefined) {
      this.setData({
        showChoice: true,
      })
      return
    }
    wx.showLoading({
      title: '加载中...',
    })
    this.data.id = options.id
    this.getConfigTpl()

    setTimeout(function() {
      wx.hideLoading()
    }, 1000)
  },

  navBack() {
    wx.navigateBack()
    this.setData({
      showChoice: false,
    })
  },

  chooseType(e) {
    this.data.id = e.currentTarget.id
    wx.showLoading({
      title: '加载中...',
    })
    this.setData({
      showChoice: false,
    })
    setTimeout(function() {
      wx.hideLoading()
    }, 1000)
    this.getAuthStatus()
    this.getConfigTpl()
  },

  onShareAppMessage: function() {

    let title = this.data.config.shareTitle
    let imageUrl = this.data.config.shareCover

    return {
      title: title ? title : "",
      path: '/module/gzhu/poster?id=' + this.data.id,
      imageUrl: imageUrl ? imageUrl : "",
      success: function(res) {
        wx.showToast({
          title: '分享成功',
          icon: "none"
        });
      }
    }
  },

  // 按键左右切换相框
  switchFrame(e) {
    let direction = e.currentTarget.dataset.direction
    let index = 0
    if (direction == "right") {
      index = (this.data.curIndex + 1) % this.data.frames.length
    } else {
      if (this.data.curIndex == 0) {
        index = this.data.frames.length - 1
      } else {
        index = (this.data.curIndex - 1) % this.data.frames.length
      }
    }
    this.setData({
      curIndex: index
    })
  },

  // 轮播切换回调
  cardSwiper(e) {
    this.data.frame = this.data.frames[e.detail.current]
    this.setData({
      avatarPath: "",
      forDrawAvatar: this.data.userInfo.avatarUrl
    })
  },

  // 海报预览
  preview(e) {
    if (!e.currentTarget.dataset.src) return
    wx.previewImage({
      urls: [e.currentTarget.dataset.src],
    })
  },

  // 绘图
  drawCanvas(e) {
    let type = e.currentTarget.dataset.type
    let config = {}

    if (type == "avatar") {
      this.setData({
        drawType: "avatar"
      })
      if (!this.data.config.avatar) {
        wx.showToast({
          title: '初始化数据失败',
          icon: "none"
        })
        return
      }
      config = this.data.config.avatar
    }
    if (type == "post") {
      this.setData({
        drawType: "post"
      })
      if (!this.data.config.post) {
        wx.showToast({
          title: '初始化数据失败',
          icon: "none"
        })
        return
      }
      config = this.data.config.post
    }

    if (!this.data.userInfo.avatarUrl) {
      wx.showToast({
        title: '读取头像失败',
        icon: "none"
      })
      return
    }
    console.log("开始绘图", config)
    // 替换头像和相框url
    for (let i = 0; i < config.images.length; i++) {
      if (config.images[i].name == "avatar")
        config.images[i].url = this.data.forDrawAvatar
      if (config.images[i].name == "frame")
        config.images[i].url = this.data.frame
    }

    wx.showLoading({
      title: '合成中...',
      mask: true
    })
    this.setData({
      posterConfig: config
    }, () => {
      Poster.create(true); // 入参：true为抹掉重新生成 
    });
  },

  // 绘图成功回调
  onPosterSuccess(e) {
    console.log("绘图结束", e)
    wx.hideLoading()
    const {
      detail
    } = e;

    let type = e.currentTarget.dataset.type
    if (type == "avatar") {
      that.setData({
        avatarPath: detail
      })
    }
    if (type == "post") {
      that.setData({
        postPath: detail,
        showPost: true
      })
      wx.showToast({
        title: '点击预览保存',
        icon: "none"
      })
    }
    this.updateCounter(this.data.recordId)
  },

  // 绘图失败回调
  onPosterFail(err) {
    wx.hideLoading()
    wx.showToast({
      title: '绘图失败：' + err.errMsg,
      icon: "none",
      duration: 3000
    })
    console.error(err);
  },

  // 保存图片
  eventSave() {
    this.saveToAlbum(this.data.avatarPath)
  },


  // 初始化配置对象
  initConfig(config) {
    if (!config) {
      wx.showToast({
        title: '初始化数据失败',
        icon: "none"
      })
      return
    }
    let that = this
    // 根据屏幕比例选取背景
    if (that.data.ratio > 17 / 9) {
      config.background = config.background18 ? config.background18 : that.data.background
    } else {
      config.background = config.background16 ? config.background16 : that.data.background
    }
    that.setData({
      background: config.background ? config.background : that.data.background, //主背景
      frames: config.frames ? config.frames : that.data.frames, //相框列表
      frame: config.frames[0] ? config.frames[0] : that.data.frame, //相框
      showCount: config.showCount,
      avatarHeight: config.avatarHeight ? config.avatarHeight : this.data.avatarHeight,
      showCount: config.showCount,
      countHeight: config.countHeight ? config.countHeight : this.data.countHeight,
    })

  },

  // 从云端获取基础配置
  getConfigTpl() {
    let that = this
    let Table = new wx.BaaS.TableObject("config")
    let query = new wx.BaaS.Query()

    query.compare('name', '=', this.data.id)
    Table.setQuery(query).limit(1).find().then(res => {
      if (res.data.objects.length == 0) {
        console.log("获取在线配置出错，使用本地配置")
        that.data.config = PostConfig.config
        that.initConfig(PostConfig.config)
        return
      }
      that.data.recordId = res.data.objects[0].id
      let config = res.data.objects[0].data
      that.data.config = config
      that.initConfig(config)
      that.setData({
        count: res.data.objects[0].addition_num
      })
      setTimeout(function() {
        wx.hideLoading()
      }, 1000)
    }, err => {
      setTimeout(function() {
        wx.hideLoading()
      }, 1000)
      console.log("获取在线配置出错，使用本地配置")
      that.data.config = PostConfig.config
      that.initConfig(PostConfig.config)
    })
  },

  // 更新统计
  updateCounter(id) {
    let table = new wx.BaaS.TableObject("config")
    let record = table.getWithoutData(id)
    record.incrementBy('addition_num', 1)
    record.update().then(res => {
      console.log(res)
      this.setData({
        count: res.data.addition_num
      })
    })
  },

  // 获取用户信息
  userInfoHandler(data) {
    wx.showLoading({
      title: '授权中...',
    })
    setTimeout(function() {
      wx.hideLoading()
    }, 3000)

    wx.BaaS.auth.loginWithWechat(data, {
      createUser: true,
      syncUserProfile: "overwrite"
    }).then(user => {
      console.log(user)
      this.getAuthStatus()
      setTimeout(function() {
        wx.hideLoading()
      }, 1000)
    }, err => {
      wx.hideLoading()
      wx.showToast({
        title: '授权失败，无法获取头像',
        icon: "none"
      })
    })
  },

  // 获取认证状态
  getAuthStatus() {
    let that = this
    wx.getSetting({
      success: res => {
        if (res.authSetting['scope.userInfo']) {

          wx.checkSession({
            success() {
              // session_key 未过期，并且在本生命周期一直有效
              wx.getUserInfo({
                success(res) {
                  console.log("已授权微信", res.userInfo)
                  res.userInfo.avatarUrl = Utils.headimgHD(res.userInfo.avatarUrl)
                  that.setData({
                    userInfo: res.userInfo,
                    isAuth: true,
                    forDrawAvatar: res.userInfo.avatarUrl
                  })
                }
              })
            },
            fail() {
              // session_key 已经失效，需要重新执行登录流程
              wx.login({
                success() {
                  wx.getUserInfo({
                    success(res) {
                      res.userInfo.avatarUrl = Utils.headimgHD(res.userInfo.avatarUrl)
                      that.setData({
                        userInfo: res.userInfo,
                        isAuth: true,
                        forDrawAvatar: res.userInfo.avatarUrl
                      })
                    }
                  })
                }
              }) // 重新登录
            }
          })
        } else {
          console.log("微信未授权")
          that.setData({
            isAuth: false
          })
          wx.showToast({
            title: '点击授权合成头像',
            icon: "none"
          })
        }
      }
    })
  },

  chooseImage() {
    wx.chooseImage({
      count: 1,
      sizeType: ['origin'],
      sourceType: ['album', "camera"],
      success: (res) => {
        console.log(res)
        that.setData({
          forDrawAvatar: res.tempFilePaths[0]
        })
        that.drawCanvas({
          currentTarget: {
            dataset: {
              type: "avatar"
            }
          }
        })
      },
    });
  },


  //保存相册
  saveToAlbum: function(filePath) {
    //查看授权状态；
    if (wx.getSetting) { //判断是否存在函数wx.getSetting在版本库1.2以上才能用
      wx.getSetting({
        success(res) {
          if (!res.authSetting['scope.writePhotosAlbum']) {
            wx.authorize({
              scope: 'scope.writePhotosAlbum',
              success(res) {
                wx.saveImageToPhotosAlbum({
                  filePath: filePath,
                  success: function(res) {
                    wx.showToast({
                      title: '图片保存成功',
                    });
                  },
                  fail: function(res) {
                    wx.showToast({
                      title: '图片保存失败',
                      icon: 'none',
                    });
                  }
                })
              },
              fail: function(res) {
                //拒绝授权时会弹出提示框，提醒用户需要授权
                wx.showModal({
                  title: '提示',
                  content: '保存图片需要授权，是否去授权',
                  success: function(res) {
                    if (res.confirm) {
                      wx.openSetting({
                        success: function(res) {}
                      })
                    }
                  }
                })
              }
            })
          } else { //已经授权
            wx.saveImageToPhotosAlbum({
              filePath: filePath,
              success: function(res) {
                wx.showToast({
                  title: '图片保存成功',
                });
              },
              fail: function(res) {
                wx.showToast({
                  title: '图片保存失败',
                  icon: 'none',
                });
              }
            })
          }
        }
      })
    } else {
      wx.showModal({
        title: '提示',
        content: '因当前微信版本过低导致无法保存，请更新至最新版本',
        showCancel: false,
        complete: function() {}
      })
    }
  },

})