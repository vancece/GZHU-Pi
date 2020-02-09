const cloud = require('wx-server-sdk')
cloud.init()
exports.main = async(event, context) => {

  // 未读消息提醒
  let unreadTpl = 'mzClt2VmH5tlVqVpbaKaeMtPX2UOW2FgbQU2Sxq3ydk'

  try {
    const result = await cloud.openapi.subscribeMessage.send({
      touser: 'off364mCDvz_Ck5XYsMV0BN80img',
      page: 'pages/Life/wall/post',
      data: {
        name1: {
          value: '哈哈' //发送人
        },
        thing2: {
          value: '2015年01月05日' //内容
        },
        time3: {
          value: '2019年10月1日' //发送时间
        },
        number5: {
          value: '2' //数量
        },
        thing7: {
          value: '广州市新港中路397号' //备注
        }
      },
      templateId: 'mzClt2VmH5tlVqVpbaKaeMtPX2UOW2FgbQU2Sxq3ydk'
    })
    console.log(result)
    return result
  } catch (err) {
    console.log(err)
    return err
  }
}