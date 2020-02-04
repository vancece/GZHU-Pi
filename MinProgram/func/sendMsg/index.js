const cloud = require('wx-server-sdk')
cloud.init()
exports.main = async (event, context) => {
  try {
    const result = await cloud.openapi.subscribeMessage.send({
      touser: 'off364mCDvz_Ck5XYsMV0BN80img',
      page: 'pages/Life/wall/post',
      data: {
        name1: {
          value: '哈哈'
        },
        thing2: {
          value: '2015年01月05日'
        },
        time3: {
          value: '2019年10月1日'
        },
        number5: {
          value: '2'
        },
        thing7: {
          value: '广州市新港中路397号'
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