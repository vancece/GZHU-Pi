/*
  时间戳格式化输出
*/
const formatTime = date => {
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hour = date.getHours()
  const minute = date.getMinutes()
  const second = date.getSeconds()

  return [year, month, day].map(formatNumber).join('-')
}

const formatNumber = n => {
  n = n.toString()
  return n[1] ? n : '0' + n
}


// 将userInfo的头像转换为高清地址
function headimgHD(imageUrl) {
  imageUrl = imageUrl.split('/'); //把头像的路径切成数组
  //把大小数值为 46 || 64 || 96 || 132 的转换为0
  if (imageUrl[imageUrl.length - 1] && (imageUrl[imageUrl.length - 1] == 46 ||
      imageUrl[imageUrl.length - 1] == 64 || imageUrl[imageUrl.length - 1] == 96 ||
      imageUrl[imageUrl.length - 1] == 132)) {
    imageUrl[imageUrl.length - 1] = 0;
  }
  imageUrl = imageUrl.join('/');  //重新拼接为字符串
  return imageUrl;
}

// 初始化知晓云SDK
function initSdk() {
  require('./sdk/sdk-wechat.3.6.0')
  let ClientID = 'd5add948fe00fbdd6cdf'
  wx.BaaS.init(ClientID, {
    autoLogin: true
  })
}

module.exports = {
  formatTime: formatTime,
  headimgHD: headimgHD,
  initSdk: initSdk,
}