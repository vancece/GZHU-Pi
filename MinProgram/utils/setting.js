/*
 *设置课表背景图，存在则先删除
 */
function setBg() {
  return new Promise(function(callback) {
    wx.chooseImage({ //选择图片
      count: 1,
      success(res) {
        wx.compressImage({ //压缩图片
          src: res.tempFilePaths[0],
          quality: 80,
          complete: function(target) {
            wx.getSavedFileList({ //获取原有缓存
              success(res) {
                if (res.fileList.length > 0) {
                  wx.removeSavedFile({ //删除原有缓存
                    filePath: res.fileList[0].filePath,
                    complete(res) {
                      wx.saveFile({ //缓存新图片
                        tempFilePath: target.tempFilePath,
                        success(res) {
                          callback(res.savedFilePath)
                        }
                      })
                    }
                  })
                } else {
                  wx.saveFile({
                    tempFilePath: target.tempFilePath,
                    success(res) {
                      callback(res.savedFilePath)
                    }
                  })
                }
              }
            })
          }
        })
      }
    })
  })
}


module.exports = {
  setBg: setBg
}