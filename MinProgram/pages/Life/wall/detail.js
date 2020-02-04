Page({

  data: {

  },

  onLoad: function(options) {
    if (!options.id) options.id = 4
    this.data.id = options.id
    this.getDetail(options.id)
  },

  onShareAppMessage: function() {
    return {
      title: '广大墙：' + this.data.detail.title,
      desc: '',
      path: '/pages/Life/wall/detail?id=' + this.data.id,
      imageUrl: this.data.detail.image[0],
      success: function(res) {
        wx.showToast({
          title: '分享成功',
          icon: "none"
        })
      }
    }
  },

  operate(e) {
    let that = this
    switch (e.target.dataset.op) {
      case "star":
        break
      case "delete":
        wx.showModal({
          title: '删除提示',
          content: '确定删除该主题吗？',
          success(res) {
            if (res.confirm) {
              if (!!that.data.detail.addi.file_ids) {
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

  // 获取详情
  getDetail(id) {
    // return
    wx.$ajax({
        url: wx.$param.server["prest"] + "/postgres/public/v_topic?id=$eq." + id,
        method: "get",
        loading: true,
        checkStatus: false,
      })
      .then(res => {
        console.log(res)

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
        this.setData({
          detail: res.data[0]
        })
      }).catch(err => {
        console.log(err)
      })
  },

  delByPk(row_id) {
    if (!row_id) return
    wx.$ajax({
        url: wx.$param.server["prest"] + "/postgres/public/t_topic?id=$eq." + row_id,
        method: "delete",
        loading: true,
        checkStatus: false,
      })
      .then(res => {
        console.log(res)

        if (res.statusCode == 200 && res.data.rows_affected == 1) {
          wx.showToast({
            title: '删除成功！',
          })
          setTimeout(function() {
            wx.$navTo("/pages/Life/wall/wall")
          }, 1000)
        } else {
          wx.showModal({
            title: '失败提示',
            content: JSON.stringify(res.data),
          })
        }

      }).catch(err => {
        console.err(err)
      })
  },

  delFile(fileIDs = []) {
    if (!fileIDs) return
    let MyFile = new wx.BaaS.File()
    MyFile.delete(fileIDs).then()
  },

})