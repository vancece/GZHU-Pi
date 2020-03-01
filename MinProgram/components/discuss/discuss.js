var utils = require("../../utils/date.js")
Component({

  properties: {
    object_id: {
      type: String,
      // id会有延迟传入，需要监听变化
      observer: function (newVal) {
        if (!newVal) {
          return
        }
        this.refresh()
      }
    },
  },


  data: {
    content: "",
    authorized: true,
    debounce: false,
  },

  lifetimes: {
    attached: function () {
      this.data.user = wx.getStorageSync("gzhupi_user")
    }
  },

  methods: {

    userInfoHandler(data) {
      let that = this
      wx.showLoading({
        title: '授权中...',
      })
      wx.BaaS.auth.loginWithWechat(data, {
        createUser: true,
        syncUserProfile: "overwrite"
      }).then(user => {
        console.log(user)
        this.setData({
          authorized: true
        })
        wx.hideLoading()
      }).catch(err => {
        wx.hideLoading()
      })
    },

    refresh() {
      this.query(this.data.object_id)
    },
    // 读取输入
    commentInput(e) {
      this.data.content = e.detail.value
    },
    // @某用户
    atSomebody(e) {
      let item = e.currentTarget.dataset.item
      let content = "@" + item.nickname + " " + this.data.content
      this.setData({
        content: content
      })
    },

    // 违规检测
    checkComment() {
      // if (this.data.content == "" || this.data.content == undefined) return
      // wx.BaaS.wxCensorText(this.data.content).then(res => {
      //   console.log(res.data.risky)
      //   if (res.data.risky) {
      //     wx.showModal({
      //       title: '警告',
      //       content: '您的发布内容包含违规词语',
      //     })
      //     return
      //   }
      this.addComment()
      // }, err => {
      //   console.log(err)
      // })
    },
    // 发布评论
    addComment() {
      wx.$subscribe()
      let that = this
      // 防抖处理
      if (this.data.debounce) return
      this.data.debounce = true
      setTimeout(() => {
        that.data.debounce = false
      }, 2000)
      let object_id = this.data.object_id
      let content = this.data.content
      if (object_id == "" || content == "") {
        console.log("illegal argument")
        return
      }
      let data = {
        object_id: Number(this.data.object_id),
        content: this.data.content,
        type: ""
      }
      this.create(data)
    },

    deleteComment(e) {

      let item = e.currentTarget.dataset.item
      let row_id = item.id

      if (this.data.user.id != item.created_by) {
        this.atSomebody(e)
        return
      }

      let that = this
      wx.showModal({
        title: '提示',
        content: '是否删除留言？',
        success(res) {
          if (res.confirm) {
            wx.$ajax({
                url: wx.$param.server["prest"] + "/postgres/public/t_discuss?id=$eq." + row_id,
                method: "delete",
                loading: true,
              })
              .then(res => {
                console.log(res)

                if (res.statusCode == 200 && res.data.rows_affected == 1) {
                  wx.showToast({
                    title: '删除成功！',
                  })
                  that.query(that.data.object_id)
                } else {
                  wx.showModal({
                    title: '失败提示',
                    content: JSON.stringify(res.data),
                  })
                }
              }).catch(err => {
                console.error(err)
              })
          }
        }
      })
    },

    create(data) {
      wx.$ajax({
          url: wx.$param.server["prest"] + "/postgres/public/t_discuss",
          method: "post",
          data: data,
          header: {
            "content-type": "application/json"
          }
        })
        .then(res => {
          this.setData({
            loading: false,
            content: ""
          })
          wx.showToast({
            title: '留言成功',
          })
          this.query(this.data.object_id)
          this.triggerEvent('success', res.data)
        })
    },

    query(object_id) {
      wx.$ajax({
          url: wx.$param.server["prest"] + "/postgres/public/v_discuss?object_id=$eq." + object_id,
          method: "get",
          header: {
            "content-type": "application/json"
          }
        })
        .then(res => {
          console.log("评论列表", res.data)
          for (let i = 0; i < res.data.length; i++) {
            let time = new Date(res.data[i].created_at)
            res.data[i].created_at = utils.relativeTime(time.getTime() / 1000)
          }
          this.setData({
            comments: res.data
          })
        })
    }
  }

})