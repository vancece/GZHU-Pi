// 全局微信变量
wx.$param = require('param').param


// 封装的wx微信全局方法

wx.$ajax = function (option) {
  return new Promise(function (resolve, reject) {
    if (option.method == undefined || typeof option.method !== "string") {
      option.method = "POST"
    }
    if (option.url == undefined) {
      option.url = wx.$param.server["aliyun_go"]
    }
    if (option.header == undefined || typeof option.header != 'object') {
      option.header = {
        "content-type": "application/x-www-form-urlencoded"
      }
    }
    if (typeof option.url === 'string' && option.url.indexOf("http") == -1) {
      option.url = wx.$param.server["aliyun_go"] + option.url
    }
    if (typeof option.loading == "boolean" && option.loading) {
      wx.showLoading({
        title: '加载中',
        duration: 60000,
        mask: true,
      })
    } else if (typeof option.loading == "string") {
      wx.showLoading({
        title: option.loading,
        duration: 60000,
        mask: true,
      })
    }
    // 携带cookie
    option.header["Cookie"] = wx.getStorageSync("gzhupi_cookie")
    wx.request({
      url: option.url,
      data: option.data,
      method: option.method.toUpperCase(),
      header: option.header,
      success: (res) => {

        // http响应错误
        if (res.statusCode >= 400) {
          if (res.statusCode == 401) wx.removeStorageSync('gzhupi_cookie')
          let msg = res.data.error
          msg = msg ? msg : res.errMsg
          reject({
            when: "http_status_error",
            error: msg,
            detail: msg,
          })
          if (option.showErr == false) return
          wx.showModal({
            title: '提示',
            content: msg,
            showCancel: false
          })
          return
        }

        // 缓存cookies
        if (res.header["Set-Cookie"] != undefined) {
          wx.setStorageSync("gzhupi_cookie", res.header["Set-Cookie"]);
        }

        // 自定义响应协议(只返回data)
        if (res.data && res.data.status) {
          if (res.data.status == 401) wx.removeStorageSync('gzhupi_cookie')
          if (res.data.status == 0 || res.data.status == 200) {
            resolve(res.data)
            return
            // } else if (res.data.status == -1) {
          } else {
            let msg = "请求响应失败"
            if (res.data.msg != undefined) msg = res.data.msg
            reject({
              when: "status_error",
              error: msg,
              detail: res.data,
            })
            if (option.showErr == false) return
            wx.showModal({
              title: '提示',
              content: msg,
              showCancel: false
            })
            return
          }
        }
        // 没有使用自定义响应协议
        resolve(res)
        return
      },
      fail: (err) => {
        reject({
          when: "origin_error",
          error: err
        })
        wx.showModal({
          title: '错误提示',
          content: JSON.stringify(err),
          showCancel: false
        })
      },
      complete: (res) => {
        console.log("response :" + option.url, res)
        wx.hideLoading()
      }
    })
  })
}

/**
 * 页面转跳封装
 * @method wx.$navTo
 * @param {object|string}  e    如果是字符串，直接跳转；对象，就解析到e.target.dataset.url
 * @param {object} args         页面参数
 */
wx.$navTo = function (e, args) {
  if (e == undefined && arg == undefined) return
  console.log('fun: navTo', e, args)
  let args_str = []
  if (typeof args === 'object') {
    for (let i in args) {
      args_str.push(i + '=' + encodeURIComponent(args[i]))
    }
    args_str = '?' + args_str.join("&")
  } else {
    args_str = ''
  }
  if (typeof e == 'object') {
    if (e.target.dataset && e.target.dataset.url) {
      wx.navigateTo({
        url: e.target.dataset.url + args_str,
        fail: err => {
          console.warn(err)
          wx.switchTab({
            url: e.target.dataset.url + args_str,
            fail: err => {
              console.err(err)
            }
          })
        }
      })
    } else if (e.currentTarget.dataset && e.currentTarget.dataset.url) {
      wx.navigateTo({
        url: e.currentTarget.dataset.url + args_str,
        fail: err => {
          console.warn(err)
          wx.switchTab({
            url: e.currentTarget.dataset.url + args_str,
            fail: err => {
              console.err(err)
            }
          })
        }
      })
    }
  } else {
    wx.navigateTo({
      url: e + args_str,
      fail: err => {
        console.warn(err)
        wx.switchTab({
          url: e + args_str,
          fail: err => {
            console.err(err)
          }
        })
      }
    })
  }
}


/**
 * 对象转url参数
 * @method wx.$objectToQuery
 * @param {object}  obj
 * @return {string} query
 */
wx.$objectToQuery = function (obj = {}) {

  if (typeof obj != 'object') {
    console.error("not object")
    return
  }
  let args_str = []
  for (let i in obj) {

    if (!obj[i]) continue

    args_str.push(i + '=' + encodeURIComponent(obj[i]))
  }
  let query = '?' + args_str.join("&")
  return query
}

//收集订阅消息，回调成功数量
wx.$subscribe = function () {

  let tpls = {
    unread: 'mzClt2VmH5tlVqVpbaKaeGZ2XM2GIYrztuGjNIjRaZw', //未读消息提醒 发送人、消息内容、备注
    comment: "qXh2oaTKaNEBF1UJCjYkTqW4vs24yQCfmShdO4SyvXg", //留言通知	文章标题、留言人、留言内容
  }
  let tmplIds = [tpls["comment"], tpls["unread"]]

  return new Promise(function (resolve) {
    wx.requestSubscribeMessage({
      tmplIds: tmplIds,
      success: (res) => {
        console.log(res)
        let subscription = []
        for (let i in tmplIds) {
          if (res[tmplIds[i]] === 'accept') {
            subscription.push({
              template_id: tmplIds[i],
              subscription_type: 'once',
            })
          }
        }
        if (subscription.length > 0) {
          resolve(subscription.length)
          wx.BaaS.subscribeMessage({
            subscription
          }).then(res => {
            console.log(res)
          }, err => {
            console.error(err)
          })
        } else {
          resolve(0)
        }
      },
      fail: (err) => {
        console.error(err)
        resolve(0)
      }
    })
  })
}