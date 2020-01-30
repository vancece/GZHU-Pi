var param = {
  //启动模式 debug verify prod
  "mode": "prod",
  // 服务器地址列表
  "server": {
    "prest": "https://1171058535813521.cn-shenzhen.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/trail/api/v1",
    "aliyun_go": "https://1171058535813521.cn-shenzhen.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/go/api/v1",
    "aliyun_py": "https://1171058535813521.cn-shanghai.fc.aliyuncs.com/2016-08-15/proxy/GZHU-API/Spider/",
    "tx": "https://ifeel.vip/api/v1",
    "cst": "https://cst.gzhu.edu.cn/gzhupi/api/v1",
    "localhost": "http://localhost:9000/api/v1"
  },
  // 教务系统位置
  "school": {
    "year_sem": "2019-2020-2", //当前学期
    "first_monday": "2019-08-26", //学期第一周周一
    "sem_list": ["2018-2019-1", "2018-2019-2", "2019-2020-1", "2019-2020-2"]
  },
  // 首页功能导航
  "nav": [{
    "show": true,
    "name": "校历",
    "icon": "https://cloud-minapp-17768.cloud.ifanrusercontent.com/1g7bLteUNRcIFOHE.png",
    "url": "/pages/Campus/tools/calendar"
  }, {
    "show": true,
    "name": "空教室",
    "icon": "https://cloud-minapp-17768.cloud.ifanrusercontent.com/1g7bLtJbWheBlxCk.png",
    "url": "/pages/Campus/tools/emptyRoom"
  }, {
    "show": true,
    "name": "查成绩",
    "icon": "https://cloud-minapp-17768.cloud.ifanrusercontent.com/1g0lGeFliOqPloZT.png",
    "url": "/pages/Campus/grade/grade"
  }, {
    "show": true,
    "name": "图书馆",
    "icon": "https://cloud-minapp-17768.cloud.ifanrusercontent.com/1g0eqSYWkdJOBqfD.png",
    "url": "/pages/Campus/library/search"
  }, {
    "show": true,
    "name": "考试",
    "icon": "https://cloud-minapp-17768.cloud.ifanrusercontent.com/1gEtLYqwE75THYVk.png",
    "url": "/pages/Campus/tools/exam"
  }, {
    "show": true,
    "name": "蹭课",
    "icon": "https://cloud-minapp-17768.cloud.ifanrusercontent.com/1gEtLYHFBZNatsdj.png",
    "url": "/pages/Campus/course/search"
  }, {
    "show": true,
    "name": "实验课",
    "icon": "/assets/exp.svg",
    "url": "/pages/Campus/exp/exp"
  }, {
    "show": true,
    "name": "校园二手",
    "icon": "https://shaw-1256261760.cos.ap-guangzhou.myqcloud.com/gzhu-pi/images/icon/ershou.svg",
    "url": "/pages/Life/oldthings/index?vaild=true"
  }]
}


module.exports = {
  param: param
}