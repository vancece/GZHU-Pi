// 全局微信变量
wx.$param = require('param').param


// 封装的wx微信全局方法

wx.$ajax = function(setting) {
    return new Promise(function(resolve, reject) {
        if (setting.method == undefined || typeof setting.method !== "string") {
            setting.method = "POST";
        }
        if (setting.url == undefined) {
            setting.url = app.globalData.url;
        }
        if (setting.retry == undefined) {
            setting.retry = 2;
        }
        setting.method = setting.method.toUpperCase();
        if (setting.method == "POST") {
            if (setting.data.key === void 0) {
                setting.data.key = wx.getStorageSync("key");
            }
        }
        if (setting.loading) {
            let load_str = '';
            if (typeof setting.loading === 'string') {
                load_str = setting.loading;
            } else {
                load_str = '加载中';
            }
            wx.showLoading({
                title: load_str,
                duration: 200000,
                mask: true,
            })

        }
        wx.request({
            url: app.globalData.url,
            data: setting.data,
            method: setting.method,
            header: {
                "content-type": "application/x-www-form-urlencoded"
            },
            success: (res) => {
                wx.hideLoading();
                if (setting.loading) {
                    wx.hideLoading();
                }
                if (res && res.data.status == 1) {
                    resolve(res.data)
                } else if (res.data.status == 0 && res.data.error) {
                    if (/登录失败.+/.test(res.data.error) && setting.retry) {
                        login().then(e => {
                            setTimeout(e => {
                                setting.retry -= 1;
                                setting.data.key = wx.getStorageSync("key");
                                ajax(setting).then(e => {
                                    resolve(e);
                                }).catch(e => {
                                    reject(e);
                                });
                            }, 100);
                        });
                    } else {
                        reject({
                            when: "create_error",
                            error: res.data.error,
                            detail: res.data,
                        });
                    }
                    // reject({
                    //     when: "create_error",
                    //     error: res.data.error,
                    //     detail: res.data,
                    // });
                } else //系统报错了
                {
                    reject({
                        when: "decode_error",
                        error: res.data
                    });
                }
            },
            fail: (res) => {
                wx.hideLoading();
                reject({
                    when: "origin_error", //原生产生的错误。网络不行或者服务器挂了
                    error: res
                });
            },

        });
    });
}

/**
 * 页面转跳封装
 * @method wx.$navTo
 * @param {object|string}  e    如果写字符，就跳转;对象的话，就解析到e.target.dataset.url
 * @param {object} args         页面参数
 */
wx.$navTo = function (e, args) {
    console.log('fun: navTo', e, args)
    let args_str = [];
    if (typeof args === 'object') {
        for (let i in args) {
            args_str.push(i + '=' + encodeURIComponent(args[i]));
        }
        args_str = '?' + args_str.join("&");
    } else {
        args_str = '';
    }
    if (typeof e == 'object') {
        if (e.target.dataset && e.target.dataset.url) {
            wx.navigateTo({
                url: e.target.dataset.url + args_str
            });
        } else if (e.currentTarget.dataset && e.currentTarget.dataset.url) {
            wx.navigateTo({
                url: e.currentTarget.dataset.url + args_str
            });
        }
    } else {
        wx.navigateTo({
            url: e + args_str
        });
    }
}