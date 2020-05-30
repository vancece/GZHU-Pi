import UserService from "../../../services/user.js"
var userService = new UserService()
Page({

  /**
   * 页面的初始数据
   */
  data: {

    v_user: {},
    bind_success: false,

  },

  onLoad: function (options) {
    console.log(options)
    let mp_open_id = options.mp_open_id


    // 绑定公众号
    if (!!mp_open_id) {
      this.bind_mp(mp_open_id)
      return
    }

    this.setData({
      v_user: wx.getStorageSync('gzhupi_user')
    })

  },
  onShareAppMessage: function () {

  },

  // 关联公众号请求
  bind_mp(mp_open_id) {
    console.log("绑定公众号：", mp_open_id)
    wx.$ajax({
        url: wx.$param.server["prest"] + "/auth?type=bind_mp&mp_open_id=" + mp_open_id,
        // url: "http://localhost:9000/api/v1/auth?type=bind_mp&mp_open_id=" + mp_open_id,
        method: "post",
        data: {},
      })
      .then(res => {
        console.log(res)

        if (res.data.id > 0) {
          this.setData({
            v_user: res.data,
            bind_success: true
          })
          wx.setStorage({
            key: 'gzhupi_user',
            data: res.data,
          })
        } else {
          wx.showToast({
            title: '未知错误',
            icon: "none"
          })
        }
      })
  },

  userInfoHandler(data) {
    let that = this
    wx.showLoading({
      title: '授权中...',
    })
    wx.BaaS.auth.loginWithWechat(data, {
      createUser: true,
      syncUserProfile: "overwrite"
    }).then(user => {
      console.log("minapp", user)
      wx.hideLoading()

      this.auth(user)

    }).catch(err => {
      console.log("拒绝授权", err)
      wx.hideLoading()
      wx.showToast({
        title: '授权失败，可退出重试',
        icon: "none"
      })
    })

  },

  auth(user) {

    let form = {
      minapp_id: user.id,
      open_id: user.openid,
      union_id: user.unionid,
      avatar: user.avatar,
      nickname: user.nickname,
      city: user.city,
      province: user.province,
      country: user.country,
      gender: user.gender,
      language: user.language,
      phone: user._phone,
    }
    wx.$ajax({
        url: wx.$param.server["prest"] + "/auth",
        // url: "http://localhost:9000/api/v1/auth",
        method: "post",
        data: form,
        header: {
          "content-type": "application/json"
        }
      })
      .then(res => {
        console.log("auth", res)
        if (res.data.open_id == user.openid) {
          wx.setStorage({
            key: 'gzhupi_user',
            data: res.data,
          })
          this.setData({
            v_user: res.data
          })
        }
      })
  },


  navTo(e){
    wx.$navTo(e)
  },

  onClick(){
    wx.showModal({
      title:"提示",
      content:"在广大派公众号发送【绑定】"
    })
  }

})