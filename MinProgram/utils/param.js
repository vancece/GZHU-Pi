var param = {
  //启动模式 debug verify prod
  "mode": "verify",
  // 服务器地址列表
  "server": {
    "aliyun_go": "https://pi.ifeel.vip/api/v1",
    "aliyun_py": "https://1171058535813521.cn-shanghai.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/Spider/",
    "cst": "https://cst.gzhu.edu.cn/gzhupi/api/v1",
    "default": "https://pi.ifeel.vip/api/v1",
    "localhost": "http://localhost:9000/api/v1",
    "prest": "https://pi.ifeel.vip/api/v1",
    "scheme": "/gzhupi/public",
    "tx": "https://ifeel.vip/api/v1"
  },
  // 教务系统位置
  "school": {
    "year_sem": "2019-2020-2", //当前学期
    "first_monday": "2020-03-02", //学期第一周周一
    "sem_list": ["2020-2021-2", "2020-2021-1", "2019-2020-2", "2019-2020-1", "2018-2019-2", "2018-2019-1"]
  },
  // 首页功能导航
  "nav": [{
    "show": true,
    "name": "校历",
    "icon": "https://shaw-1256261760.cos.ap-guangzhou.myqcloud.com/gzhu-pi/images/icon/xiaoli-red.png",
    "url": "/pages/Campus/tools/calendar"
  }, {
    "show": true,
    "name": "空教室",
    "icon": "https://shaw-1256261760.cos.ap-guangzhou.myqcloud.com/gzhu-pi/images/icon/building1.png",
    "url": "/pages/Campus/tools/emptyRoom"
  }, {
    "show": true,
    "name": "查成绩",
    "icon": "https://shaw-1256261760.cos.ap-guangzhou.myqcloud.com/gzhu-pi/images/icon/grade-a%2B.png",
    "url": "/pages/Campus/grade/grade"
  }, {
    "show": true,
    "name": "图书馆",
    "icon": "https://shaw-1256261760.cos.ap-guangzhou.myqcloud.com/gzhu-pi/images/icon/library.png",
    "url": "/pages/Campus/library/search"
  }, {
    "show": true,
    "name": "考试",
    "icon": "https://shaw-1256261760.cos.ap-guangzhou.myqcloud.com/gzhu-pi/images/icon/exam-green.png",
    "url": "/pages/Campus/tools/exam"
  }, {
    "show": true,
    "name": "蹭课",
    "icon": "https://shaw-1256261760.cos.ap-guangzhou.myqcloud.com/gzhu-pi/images/icon/rili.png",
    "url": "/pages/Campus/course/search"
  }, {
    "show": true,
    "name": "实验课",
    "icon": "/assets/exp.svg",
    "url": "/pages/Campus/exp/exp"
  }, {
    "show": false,
    "name": "校园二手",
    "icon": "https://shaw-1256261760.cos.ap-guangzhou.myqcloud.com/gzhu-pi/images/icon/ershou.svg",
    "url": "/pages/Life/oldthings/index?vaild=true"
  }, {
    "icon": "https://cos.ifeel.vip/gzhu-pi/images/pic/rank.png",
    "name": "成绩排名",
    "show": true,
    "url": "/pages/Campus/grade/rank"
  }]
}


module.exports = {
  param: param
}