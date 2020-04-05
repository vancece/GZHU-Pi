var config = {
  config_data: "20190316",
  schedule_mode: "day", //启动模式 day/week
  showExp: false, //是否展示实验课
  version: "2.0.0.20200405", //版本号
  bindStatus: wx.getStorageSync("account") == "" ? false : true, //学号绑定
  schedule_bg: "#a8edea,#fed6e3", //课表背景图片/颜色
  blur: 8, //高斯模糊
  tips: {
    time_line: true, //时间轴
    help: false, //帮助
    star: false, //加入喜欢
    share: false, //分享提示
  }
}


// 生成config加入缓存
function init() {
  let conf = wx.getStorageSync("config")
  if (conf == "") wx.setStorageSync('config', config)
}


// 重新生成config加入缓存
function reInit() {
  wx.setStorage({
    key: 'config',
    data: config,
  })
}


// 获取缓存项
function get(propString = "") {
  let config = wx.getStorageSync("config")
  if (config == "") init()
  return getPropByString(config, propString)
}


// 修改/设置缓存
function set(name, value) {
  let config = wx.getStorageSync("config")
  if (config == "") init()
  config[name] = value
  wx.setStorageSync("config", config)
}

// 获取多级对象属性
function getPropByString(obj, propString) {
  if (!propString)
    return obj;

  var prop, props = propString.split('.');

  for (var i = 0, iLen = props.length - 1; i < iLen; i++) {
    prop = props[i];

    var candidate = obj[prop];
    if (candidate !== undefined) {
      obj = candidate;
    } else {
      break;
    }
  }
  return obj[props[i]];
}


module.exports = {
  config: config,
  init: init,
  reInit: reInit,
  get: get,
  set: set
}