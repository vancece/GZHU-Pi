// 公共数据

// 星期栏
const weekDays = ['一', '二', '三', '四', '五', '六', '日']

// 课表时间轴
const timeLine = ["08:30-09:15", "09:20-10:05", "10:25-11:10", "11:15-12:00", "13:50-14:35", "14:40-15:25", "15:45-16:30", "16:35-17:20", "18:20-19:05", "19:10-19:55", "20:00-20:45"]

// 课表颜色
const colors = ["#86b0fe", "#71eb55", "#f7c156", "#76e9eb", "#ff9dd8", "#80f8e6", "#eaa5f7", "#86b3a5", "#85B8CF", "#90C652", "#D8AA5A", "#FC9F9D", "#29ab97", "#61BC69", "#12AEF3", "#E29AAD", "#AFD7A4", "#F1BBB9", "#A0A8AE", "#AD918C", "#ff5f67"]

// 课表数据样例
const course_sample = [{
  "color": 0,
  "weekday": 1,
  "start": 1,
  "last": 3,
  "weeks": "1-16周",
  "course_name": "↖点击左上方周数回到校历本周"
},
{
  "color": 8,
  "weekday": 1,
  "start": 1,
  "last": 3,
  "weeks": "1-16周",
  "course_name": "↖点击左上方周数回到校历本周"
},
{
  "color": 11,
  "weekday": 1,
  "start": 5,
  "last": 3,
  "weeks": "1-16周",
  "course_name": "←点击侧边栏查看课程时间轴←"
},
{
  "color": 2,
  "weekday": 3,
  "start": 1,
  "last": 2,
  "weeks": "1-16周",
  "course_name": "点开小格子添加课程"
},
{
  "color": 7,
  "weekday": 4,
  "start": 4,
  "last": 3,
  "weeks": "1-16周",
  "course_name": "有问题记得联系开发者！",
},
{
  "color": 9,
  "weekday": 6,
  "start": 3,
  "last": 2,
  "weeks": "1-16周",
  "course_name": "点击图标查看更多",
},
{
  "color": 10,
  "weekday": 3,
  "start": 8,
  "last": 4,
  "weeks": "1-16周",
  "course_name": "觉得好用的话，请分享喔！",
  "class_place": "点击小飞机图标"
},
{
  "color": 3,
  "weekday": 2,
  "start": 3,
  "last": 3,
  "weeks": "1-16周",
  "course_name": "←左右滑动切换周数 →",
  "class_place": ""
},
{
  "color": 4,
  "weekday": 5,
  "start": 5,
  "last": 3,
  "weeks": "1-16周",
  "course_name": "点击图标切换模式",
  "class_place": ""
}
]

module.exports = {
  weekDays: weekDays,
  timeLine: timeLine,
  colors: colors,
  course_sample: course_sample
}