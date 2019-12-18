const Page = require('../../../utils/sdk/ald-stat.js').Page;
let app = getApp()
Page({

  data: {
    success: false, //用于发布超时检测
    loading: false,
    imgList: [], //临时图片地址
    imgUrls: [], //上传后的图片地址
    file_ids: [], //图片文件id，方便删除
    label: [],
    category: ["图书文具", "生活用品", "电子产品", "化妆用品", "服装鞋包", "其它"],
    categoryIndex: 5,
    isBuy: false, //是否求购
    hasPhone: false, //是否绑定手机号码
    debounce: false, //防抖
  },

  onLoad: function(options) {
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
    wx.BaaS.auth.getCurrentUser().then(user => {
      console.log(user)
      if (user._phone) {
        this.setData({
          hasPhone: true
        })
      }
    })
  },

  // 保存用户手机号码
  getPhoneNumber(e) {
    wx.BaaS.wxDecryptData(e.detail.encryptedData, e.detail.iv, 'phone-number').then(res => {
      wx.BaaS.auth.getCurrentUser().then(user => {
        user.setMobilePhone(res.phoneNumber)
      })
    })
  },

  // 不同意，返回上一页
  navBack() {
    wx.navigateBack({})
  },

  // 获取用户信息
  userInfoHandler(data) {
    let that = this

    wx.BaaS.auth.loginWithWechat(data, {
      createUser: true,
      syncUserProfile: "overwrite"
    }).then(user => {
      console.log(user)
    })
  },

  // 添加标签
  labelAdd: function() {
    if (this.data.inputValue == "" || this.data.inputValue == undefined) return
    if (this.data.label.length == 3) {
      wx.showToast({
        title: '最多加3个标签',
        icon: "none"
      })
      return
    }
    this.setData({
      label: this.data.label.concat(this.data.inputValue),
      inputValue: ""
    })
  },
  // 删除标签
  labelDel(e) {
    let id = Number(e.target.id)
    this.data.label.splice(id, 1)
    this.setData({
      label: this.data.label
    })
  },
  // 读取标签内容
  labelInput: function(e) {
    this.data.inputValue = e.detail.value
  },

  // 求购开关
  switchBuy(e) {
    this.setData({
      isBuy: e.detail.value
    })
    if (e.detail.value) {
      wx.showToast({
        title: '发布求购',
        icon: "none"
      })
      this.setData({
        label: this.data.label.concat("求购")
      })
    } else {
      wx.showToast({
        title: '发布二手',
        icon: "none"
      })

      // 删除求购标签
      Array.prototype.remove = function(val) {
        var index = this.indexOf(val);
        if (index > -1) {
          this.splice(index, 1);
        }
      };
      this.data.label.remove("求购")
      this.setData({
        label: this.data.label
      })
    }
  },

  // 选择分类
  selectCategory() {
    let that = this
    wx.showActionSheet({
      itemList: this.data.category,
      success(res) {
        that.setData({
          categoryIndex: res.tapIndex
        })
      }
    })
  },
  // 违规检测
  formSubmit(e) {
    let v = e.detail.value
    if (v.title == "" || v.content == "" || v.price == "" || v.phone == "" || v.wechat == "") {
      wx.showToast({
        title: '信息不完整',
        icon: "none"
      })
      return
    }
    if (this.data.imgList.length == 0) {
      wx.showToast({
        title: '至少上传一张图片',
        icon: "none"
      })
      return
    }

    if (v.phone != "" && !(/^1\d{10}$/.test(v.phone))) {
      wx.showToast({
        title: '手机号格式不正确，请检查！',
        icon: 'none',
        duration: 1500
      })
      return
    }
    let text = v.title + v.content + v.price + v.phone + v.wechat
    wx.BaaS.wxCensorText(text).then(res => {
      console.log(res.data.risky)
      if (res.data.risky) {
        wx.showModal({
          title: '警告',
          content: '您的发布内容包含违规词语',
        })
        return
      }
      // 通过校验
      this.submit(e)
    }, err => {
      console.log(err)
    })
  },

  // 表单校验提交
  submit: function(e) {
    wx.BaaS.wxReportTicket(e.detail.formId)
    if (!app.globalData.bindStatus) {
      wx.showToast({
        title: '未绑定学号',
        icon: "none"
      })
      return
    }

    var that = this
    let v = e.detail.value

    wx.showModal({
      title: '提示',
      content: '确认立即发布',
      success(res) {
        if (res.confirm) {

          that.setData({
            debounce: true
          })
          //记录发布者学号，如有
          let stu = wx.getStorageSync("account")
          let stu_id = stu ? stu.username : ""

          let recordData = {
            title: e.detail.value.title,
            content: e.detail.value.content,
            price: Number(e.detail.value.price),
            info: {
              name: e.detail.value.name,
              phone: e.detail.value.phone,
              wechat: e.detail.value.wechat,
              file_ids: that.data.file_ids,
              isBuy: that.data.isBuy,
              stu_id: stu_id
            },
            label: that.data.label,
            image: [],
            category: that.data.category[that.data.categoryIndex],
            // 刷新时间，秒级时间戳
            refresh_time: Date.parse(new Date()) / 1000
          }
          // debugger;
          that.saveRecord(recordData)
        }
      }
    })
  },

  // 上传图片并保存记录到数据库
  saveRecord: function(data) {
    let that = this
    this.setData({
      loading: true
    })
    // 批量同步上传图片
    for (let i = 0; i < this.data.imgList.length; i++) {
      this.uploadFile('flea_market', this.data.imgList[i]).then(res => {
        console.log("文件上传", res)

        this.checkTimeout(res.file.id)
        this.data.imgUrls.push(res.path)
        this.data.file_ids.push(res.file.id)
        console.log(i + 1, this.data.imgList.length)

        // 图片上传完
        if (i + 1 == this.data.imgList.length) {
          console.log(i, this.data.imgList.length)
          data["image"] = this.data.imgUrls
          console.log("表单数据", data)

          let Table = new wx.BaaS.TableObject('flea_market')
          let record = Table.create()
          record.set(data).save().then(res => {
            console.log("数据保存", res.data)
            this.setData({
              loading: false,
              success: true
            })
            wx.showToast({
              title: '发布成功',
            })
            wx.redirectTo({
              url: '/pages/Life/oldthings/detail?id=' + res.data.id,
            })
          }, err => {
            this.setData({
              loading: false,
              debounce: false
            })
            wx.showToast({
              title: '发布失败' + err,
            })
          })
        }
      })
    }
  },

  // 异步上传单个文件
  uploadFile: function(categoryName, filePath) {
    let MyFile = new wx.BaaS.File()
    let metaData = {
      categoryName: categoryName
    }
    //返回上传文件后的信息
    return new Promise(function(callback) {
      let fileParams = {
        filePath: filePath
      }
      MyFile.upload(fileParams, metaData).then(res => {
        callback(res.data)
      }, err => {
        wx.showToast({
          title: '上传失败',
          icon: "none"
        })
      })
    })
  },

  chooseImage() {
    wx.chooseImage({
      count: 3 - this.data.imgList.length, //默认9
      sizeType: ['compressed'],
      sourceType: ['album', "camera"],
      success: (res) => {
        this.checkImage(res.tempFilePaths)
        this.setData({
          imgList: this.data.imgList.concat(res.tempFilePaths)
        })
      }
    });
  },

  checkImage(tempFilePaths = []) {
    for (let i = 0; i < tempFilePaths.length; i++) {
      wx.BaaS.wxCensorImage(tempFilePaths[i]).then(res => {
        console.log("图片检测",res.risky)
        if (res.risky) {
          wx.showModal({
            title: '警告',
            content: '您的发布图片包含违规内容',
          })
          this.setData({
            imgList: []
          })
        }
      }, err => {
        // HError 对象
      })
    }

  },

  viewImage(e) {
    wx.previewImage({
      urls: this.data.imgList,
      current: e.currentTarget.dataset.url
    });
  },

  deleteImg(e) {
    wx.showModal({
      title: '提示',
      content: '确定要删除这张照片吗？',
      cancelText: '留着',
      confirmText: '再见',
      success: res => {
        if (res.confirm) {

          this.data.imgList.splice(e.currentTarget.dataset.index, 1);
          this.setData({
            imgList: this.data.imgList
          })
        }
      }
    })
  },

  checkTimeout(files) {
    let that = this
    setTimeout(function() {
      if (!that.data.success) {
        that.delFile(files)
        that.setData({
          loading: false,
          debounce: false
        })
        wx.showToast({
          title: '响应超时，请检查网络',
          icon: "none"
        })
      } else {
        console.log("发布未超时")
      }
    }, 60 * 1000)
  },

  delFile(fileIDs = []) {
    if (!fileIDs) return
    let MyFile = new wx.BaaS.File()
    MyFile.delete(fileIDs).then()
  },

})