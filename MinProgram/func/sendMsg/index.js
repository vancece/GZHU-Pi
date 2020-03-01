const cloud = require('wx-server-sdk')
cloud.init()
exports.main = async (event, context) => {

  let tpls = {
    unread: 'mzClt2VmH5tlVqVpbaKaeGZ2XM2GIYrztuGjNIjRaZw', //未读消息提醒 发送人、消息内容、备注
    comment: "qXh2oaTKaNEBF1UJCjYkTqW4vs24yQCfmShdO4SyvXg", //留言通知	文章标题、留言人、留言内容
  }
  let tpl = tpls[event.type]
  if (!tpl) {
    return "unknown type of " + event.type
  }

  let data = {}
  if (event.type == "unread") {
    data = {
      name1: {
        value: event.sender //发送人
      },
      thing2: {
        value: event.content //内容
      },
      thing7: {
        value: event.remark ? event.remark : "无" //备注
      }
    }
  }

  if (event.type == "comment") {
    data = {
      thing4: {
        value: event.title //文章标题
      },
      name1: {
        value: event.sender //留言人
      },
      thing2: {
        value: event.content //留言内容
      }
    }
  }

  try {
    const result = await cloud.openapi.subscribeMessage.send({
      touser: event.touser,
      page: event.page,
      data: data,
      templateId: tpl
    })
    return result
  } catch (err) {
    return err
  }
}