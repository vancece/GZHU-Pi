var config = {
  configData: "20190316",
  schedule_mode: "day", //启动模式 day/week
  showExp: false, //是否展示实验课
  version: "1.0.6.20190930", //版本号
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
function get(name = "") {
  let config = wx.getStorageSync("config")
  if (config == "") init()
  if (name == "") return config
  return config[name]
}


// 修改/设置缓存
function set(name, value) {
  let config = wx.getStorageSync("config")
  if (config == "") init()
  config[name] = value
  wx.setStorageSync("config", config)
}


module.exports = {
  config: config,
  init: init,
  reInit: reInit,
  get: get,
  set: set
}