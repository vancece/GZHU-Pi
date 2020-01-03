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
        wx.request({
            url: option.url,
            data: option.data,
            method: option.method.toUpperCase(),
            header: {
                "content-type": "application/x-www-form-urlencoded"
            },
            success: (res) => {
                if (res && (res.data.status == 0 || res.data.status == 200)) {
                    resolve(res.data)
                    // } else if (res.data.status == -1) {
                } else {
                    let msg = "请求响应失败"
                    if (res.data.msg != undefined) msg = res.data.msg
                    reject({
                        when: "status_error",
                        error: msg,
                        detail: res.data,
                    })
                    wx.showModal({
                        title: '提示',
                        content: msg,
                        showCancel: false
                    })
                }
            },
            fail: (err) => {
                reject({
                    when: "origin_error",
                    error: err
                })
                wx.showModal({
                    title: '错误提示',
                    content: err,
                    showCancel: false
                })
            },
            complete: (res) => {
                console.log("response:", res)
                wx.hideLoading()
            }
        })
    })
}

/**
 * 页面转跳封装
 * @method wx.$navTo
 * @param {object|string}  e    如果写字符，就跳转对象的话，就解析到e.target.dataset.url
 * @param {object} args         页面参数
 */
wx.$navTo = function (e, args) {
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
                url: e.target.dataset.url + args_str
            })
        } else if (e.currentTarget.dataset && e.currentTarget.dataset.url) {
            wx.navigateTo({
                url: e.currentTarget.dataset.url + args_str
            })
        }
    } else {
        wx.navigateTo({
            url: e + args_str
        })
    }
}