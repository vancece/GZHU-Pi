// 云函数入口文件
const cloud = require('wx-server-sdk')

cloud.init()

// 云函数入口函数
exports.main = async (event, context) => {
  try {
    const result = await cloud.openapi.analysis.getDailySummary({
      beginDate: '20200122',
      endDate: '20200122'
    })
    console.log(result)
    return result
  } catch (err) {
    console.log(err)
    return err
  }
}