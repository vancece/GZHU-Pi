const App = require('../sdk/ald-stat.js').App
import Utils from '../utils.js';
import Config from './config.js';
Utils.initSdk();

Page({

  data: {
    statusBarHeight: wx.getSystemInfoSync().statusBarHeight,
    tmp_list: [],
    selected: [], //选择的三位选手信息
    // picked: "",//线下投票对象
    votedToday: false, //今天是否投票

    key_for: "", //投票对象record
    key: "", //输入的线下投票码

    submit_times: 0, //提交投票次数
    end: true
  },

  // 启动，获取数据，检测投票
  onLoad: function (options) {
    // this.setData({
    //   config: Config.config
    // })
    let belong_to = options.id
    if (belong_to == undefined)
      belong_to = "photo_contest_2019"
    this.data.belong_to = belong_to

    this.getConfig(belong_to)

  },
  onShareAppMessage: function () {
    return {
      title: this.data.config["title"],
      desc: this.data.config["sub_title"],
      path: '/module/vote/vote?id=' + this.data.belong_to,
      imageUrl: "",
      success: function (res) {
        wx.showToast({
          title: '分享成功',
          icon: "none"
        });
      }
    }
  },

  navToRule() {
    let that = this
    wx.navigateTo({
      url: '/module/vote/rule?url=' + that.data.config["rule_pic_url"],
    })
  },

  userInfoHandler(data) {

    wx.BaaS.auth.loginWithWechat(data, {
      createUser: true
    }).then(user => {
      console.log(user)
    })
  },

  // 提交
  formSubmit(e) {
    if (this.data.config["increase_per_vote"] == 0 && e.detail.value.key == "") {
      wx.showModal({
        title: '提示',
        content: '举办方设置了只能使用线下投票码进行投票，请输入投票码！',
      })
      return
    }
    wx.BaaS.wxReportTicket(e.detail.formId)
    if (this.data.selected.length != this.data.config["vote_per_user"]) {
      wx.showToast({
        title: '要选' + this.data.config["vote_per_user"] + '项喔~',
        icon: "none"
      })
      return
    }

    this.data.key = e.detail.value.key

    console.log("密钥：", e.detail.value.key)
    // 没有密钥，线上投票
    if (e.detail.value.key == "") {
      this.setRecord(this.data.selected)
      for (let i = 0; i < this.data.selected.length; i++) {
        this.add(this.data.selected[i].id, this.data.config["increase_per_vote"])
      }
    } else {
      this.setData({
        show: true
      })
    }
  },

  confirm() {
    let that = this
    // 检测、校验、提交
    if (this.data.key_for == "") {
      wx.showToast({
        title: '请选择一项',
        icon: "none"
      })
      return
    }

    for (let i = 0; i < this.data.selected.length; i++) {
      // 线下投票
      if (this.data.selected[i].perform == this.data.key_for) {
        // 校验密钥
        let Obj = new wx.BaaS.TableObject(70508)
        let query = new wx.BaaS.Query()
        query.compare('key', '=', this.data.key)
        Obj.setQuery(query).find().then(res => {
          // 密钥不存在
          if (res.data.objects.length == 0) {
            that.add(that.data.selected[i].id, that.data.config["increase_per_vote"])
            that.setRecord(that.data.selected)
            wx.showToast({
              title: '密钥不存在',
              icon: "none"
            })
          }
          // 密钥已使用
          else if (res.data.objects[0].used) {
            that.add(that.data.selected[i].id, that.data.config["increase_per_vote"])
            that.setRecord(that.data.selected)
            wx.showToast({
              title: '密钥已使用',
              icon: "none"
            })
          } else {
            // 密钥正确
            that.setRecord(that.data.selected, that.data.key_for, that.data.key)
            that.add(that.data.selected[i].id, that.data.config["increase_per_vote"] + that.data.config["vote_per_key"])
            // 更新密钥状态
            let recordID = res.data.objects[0].id
            let MyRecord = Obj.getWithoutData(recordID)
            MyRecord.set("used", true)
            MyRecord.update()
          }
        })

      } else {
        this.add(this.data.selected[i].id, that.data.config["increase_per_vote"])
      }
    }

  },
  // 点击选择线下投票对象，记录节目名称
  keyFor(e) {
    let perform = e.currentTarget.dataset.perform
    this.setData({
      key_for: perform
    })
  },

  // 选择控制，选择三项
  select(e) {
    let id = Number(e.currentTarget.dataset.id)
    let record = this.data.vote_list[id]
    if (this.data.selected.indexOf(record) == -1 && this.data.selected.length >= this.data.config["vote_per_user"]) {
      wx.showToast({
        title: '最多选择' + this.data.config["vote_per_user"] + "项",
        icon: "none"
      })
      return
    }
    if (this.data.tmp_list[id] == true) {
      this.data.tmp_list[id] = false
    } else {
      this.data.tmp_list[id] = true
    }
    this.data.selected = []
    for (let i = 0; i < this.data.tmp_list.length; i++) {
      if (this.data.tmp_list[i]) {
        this.data.selected.push(this.data.vote_list[i])
      }
    }
    this.setData({
      tmp_list: this.data.tmp_list,
      selected: this.data.selected
    })
  },


  // 检测今天是否投票
  checkToday(belong_to) {
    let that = this
    // 获取用户id
    wx.BaaS.auth.loginWithWechat().then(user => {
      let id = user.get("id")
      // 根据用户id获取
      check(id)
    }, err => {
      wx.showToast({
        title: '获取用户id失败',
        icon: "none"
      })
    })

    function check(id = 0) {
      wx.showLoading({
        title: '更新投票状态',
      })
      // 投票记录表
      let Obj = new wx.BaaS.TableObject(70507)
      let query = new wx.BaaS.Query()
      // 云函数获取时间
      wx.BaaS.invokeFunction('getTime').then(e => {
        let serverDate = e.data.timeStr
        console.log("服务器时间", serverDate)
        console.log(e.data)
        query.compare('created_by', '=', id)
        query.compare('belong_to', '=', belong_to)
        Obj.setQuery(query).find().then(res => {
          that.setData({
            end: e.data.end
          })
          res.data.objects.forEach(item => {
            let date = new Date(item.created_at * 1000);
            let voteDate = Utils.formatTime(date)
            // 标记今天已经投票

            if (voteDate == serverDate) {
              console.log("今天已经投票")
              that.setData({
                votedToday: true,
                end: e.data.end
              })
            } else {
              that.setData({
                votedToday: false,
                end: e.data.end
              })
            }
          })
          wx.hideLoading()
        }, err => {
          wx.hideLoading()
          wx.showModal({
            title: '错误提示',
            content: err,
          })
        })
      })
    }
  },

  // 获取投票数据
  get_vote_items(belong_to, offset = 0) {
    wx.showLoading({
      title: '更新数据',
    })
    let that = this
    let query = new wx.BaaS.Query()
    let Product = new wx.BaaS.TableObject(70500)

    query.compare('promotion', '=', true)
    query.compare('belong_to', '=', belong_to)

    Product.setQuery(query).orderBy('No').limit(that.data.config["limit"]).offset(offset).find().then(res => {
      // console.log(res.data.objects)
      that.setData({
        vote_list: res.data.objects
      })
      wx.hideLoading()
    }, err => {
      console.log(err)
    })
  },

  // 图片预览
  viewImg(e) {
    wx.previewImage({
      urls: [e.currentTarget.dataset.img],
    })
  },

  // 记录投票历史
  setRecord(object, key_for = "", key = "") {
    let Obj = new wx.BaaS.TableObject(70507)
    let record = Obj.create()
    let vote_ids = []
    for (let i = 0; i < object.length; i++) {
      vote_ids[i] = object[i].id
    }
    let voteData = {
      vote_ids: vote_ids,
      vote0: object[0].perform,
      // vote1: object[1].perform,
      key_for: key_for,
      key: key,
      belong_to: this.data.config["belong_to"]
    }
    record.set(voteData).save().then(res => {
      console.log(res)
    }, err => {
      console.log(err)
    })
  },


  // 提交投票，原子性增加
  add(recordID, amount = this.data.config["increase_per_vote"]) {
    let that = this
    wx.showLoading({
      title: '提交ing...',
    })
    let Obj = new wx.BaaS.TableObject(70500)
    let record = Obj.getWithoutData(recordID)
    record.incrementBy('count', amount)
    record.update().then(res => {

      // n个选手都提交后，刷新数据
      that.data.submit_times = that.data.submit_times + 1

      if (that.data.submit_times == that.data.config["vote_per_user"]) {
        wx.showToast({
          title: '投票成功',
          icon: "none"
        })
        wx.hideLoading()
        that.setData({
          tmp_list: [],
          selected: [],
          votedToday: true,
          key_for: "",
          key: "",
          submit_times: 0,
          show: false
        })

        that.get_vote_items(this.data.config["belong_to"])
        that.checkToday(that.data.config["belong_to"])

      }

    }, err => {
      wx.showToast({
        title: '连接错误，投票失败',
        icon: "none"
      })
    })

  },

  // 获取线上配置并初始化
  getConfig(object) {
    let that = this
    let Table = new wx.BaaS.TableObject("config")
    let query = new wx.BaaS.Query()
    query.compare('name', '=', object)
    Table.setQuery(query).limit(1).find().then(res => {
      if (res.data.objects.length == 0) {
        console.log("获取在线配置出错，使用本地配置")
        this.data.config = Config.config
        this.get_vote_items(this.data.config["belong_to"])
        this.checkToday(this.data.config["belong_to"])
        that.setData({
          config: Config.config
        })
        return
      }
      let config = res.data.objects[0].data
      this.data.config = config

      // 初始化
      this.get_vote_items(this.data.config["belong_to"])
      this.checkToday(this.data.config["belong_to"])
      that.setData({
        config: config
      })
      // this.create_key()
    }, err => {
      wx.hideLoading()
      console.log("获取在线配置出错，使用本地配置")
      this.data.config = Config.config
      this.get_vote_items(this.data.config["belong_to"])
      this.checkToday(this.data.config["belong_to"])
      this.setData({
        config: Config.config
      })
    })
  },

  // 批量创建投票码
  create_key() {
    for (let i = 0; i < 200; i++) {
      let key = String(Math.random()).slice(2, 8)

      let Obj = new wx.BaaS.TableObject(70508)
      let query = new wx.BaaS.Query()
      query.compare('key', '=', key)
      Obj.setQuery(query).find().then(res => {

        if (res.data.objects.length == 0) {
          let record = Obj.create()
          let voteData = {
            key: key,
            belong_to: this.data.config["belong_to"]
          }
          record.set(voteData).save().then(res => {
            console.log("创建：", res.data.key)
          }, err => {
            console.log(err)
          })
        }

      }, err => {
        console.log(err)
      })
    }
  }
})