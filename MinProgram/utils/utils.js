var Config = require("config.js")
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

  return [year, month, day].map(formatNumber).join('-') + ' ' + [hour, minute, second].map(formatNumber).join(':')
}

const formatNumber = n => {
  n = n.toString()
  return n[1] ? n : '0' + n
}


/*
  获取当前校历周
*/
function getSchoolWeek() {
  let schoolWeek
  let startMonday = new Date(2019, 1, 25)
  let today = new Date()

  let interval = today - startMonday
  let intervalDays = interval / (1000 * 60 * 60 * 24)

  if (interval < 0) {
    schoolWeek = Math.ceil(Math.abs(intervalDays)) / 7
    return -(Math.ceil(schoolWeek))
  } else {
    schoolWeek = Math.ceil(intervalDays) / 7
    return Math.ceil(schoolWeek)
  }
}


/*
  设置周对应日期
*/
function setWeekDate(intervalWeeks = 0) {
  const week = [];
  for (let i = 0; i < 7; i++) {
    let Stamp = new Date();
    let weekday = Stamp.getDay() == 0 ? 7 : Stamp.getDay() //周日设置值为7
    let num = intervalWeeks * 7 - weekday + 1 + i;
    Stamp.setDate(Stamp.getDate() + num);
    // week[i] = (Stamp.getMonth() + 1) + '月' + Stamp.getDate() + '日';
    week[i] = Stamp.getDate()
  }
  return week;
}


/*
  选择出当天星期几的课程，包括非本周的
*/
function getTodayCourse() {
  let weekday = new Date().getDay()
  let course = wx.getStorageSync("course")
  let kbList = course == "" ? [] : course.course_list

  let exp = wx.getStorageSync("exp")
  if (Config.get("showExp") && exp != "") {
    kbList = kbList.concat(exp)
  }
  let todayCourse = []
  // if (getSchoolWeek() <= 0) return todayCourse
  if (kbList) {
    kbList.forEach(function(item) {
      if (item.weekday == weekday) {
        todayCourse.push(item)
      }
    })
  }
  return todayCourse
}

module.exports = {
  formatTime: formatTime,
  getSchoolWeek: getSchoolWeek,
  setWeekDate: setWeekDate,
  getTodayCourse: getTodayCourse
}