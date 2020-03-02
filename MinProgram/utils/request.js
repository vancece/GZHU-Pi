/*
 * 统一账户POST数据请求
 * 
 * @username: 用户名
 * @password: 密码
 * @type: 请求API类型
 * @account_key: 存储的账户类型
 * 
 * return 回调函数，返回请求结果
 */
var url = "https://1171058535813521.cn-shanghai.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/Spider/"
// var url = "http://127.0.0.1:5000/"

function sync(username, password, type, account_key = "account",year_sem="2019-2020-1") {

  var time = new Date()
  if (time.getHours() >= 0 && time.getHours() < 7) {
    wx.showToast({
      title: '00:00~07:00不可同步',
      icon: "none"
    })
    return new Promise(function(callback) {
      callback(0)
    })
  }

  wx.showLoading({
    title: '同步中...',
    mask: true
  })

  let account = {
    username: username,
    password: password
  }

  return new Promise(function(callback) {
    let data = account
    data["year_sem"] = year_sem
    wx.request({
      url: url + type,
      method: "POST",
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      data: data,

      success: function(res) {
        if (res.statusCode != 200) {
          callback("请求超时")
          return
        }
        if (res.data.statusCode != 200) {
          callback("账号或密码错误")
          return
        }
        // 缓存账户信息
        wx.setStorageSync(account_key, account)
        // 缓存结果数据
        wx.setStorageSync(type, res.data.data)
        callback("同步完成")
      },

      fail: function(err) {
        callback("服务器响应错误")
      },

      complete: function(res) {
        console.log(res)
        callback(res)
      }
    })
  });
}

function insertToTable(tableName,data){
  let table = tableName
  let Product = new wx.BaaS.TableObject(table)
  let product = Product.create()
  let apple = data
  product.set(apple).save().then(res => {
    // success
    console.log(res)
  }, err => {
    //err 为 HError 对象
  })
}


module.exports = {
  sync: sync,
  insertToTable: insertToTable
}