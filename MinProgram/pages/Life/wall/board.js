// pages/Life/wall/board.js
Page({

  /**
   * 页面的初始数据
   */
  data: {

  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    this.picture()
  },

  picture: function () {  //生成图片
    let that = this;
    let p1 = new Promise(function (resolve, reject) {
      wx.getImageInfo({
        src: "https://media.ifanrusercontent.com/tavatar/10/a8/10a85d9763bde48259e0a7d4a77506012aa5a294.jpg",
        success(res) {
          resolve(res);
        }
      })
    }).then(res => {
      const ctx = wx.createCanvasContext('shareCanvas');
      ctx.drawImage(res.path, 0, 0, 150, 150);    //绘制背景图
      //ctx.setTextAlign('center')    // 文字居中
      // ctx.setFillStyle('#000000')  // 文字颜色：黑色
      // ctx.setFontSize(20)         // 文字字号：22px
      // ctx.fillText("文本内容", 20, 70) //开始绘制文本的 x/y 坐标位置（相对于画布） 
      // ctx.stroke();//stroke() 方法会实际地绘制出通过 moveTo() 和 lineTo() 方法定义的路径。默认颜色是黑色
      ctx.draw(false, that.drawPicture());//draw()的回调函数 
      console.log(res.path);
    })
  },
  drawPicture: function () { //生成图片
    var that = this;
    setTimeout(function () {
      wx.canvasToTempFilePath({ //把当前画布指定区域的内容导出生成指定大小的图片，并返回文件路径
        x: 0,
        y: 0,
        width: 250,
        height: 250,
        destWidth: 300,  //输出的图片的宽度（写成width的两倍，生成的图片则更清晰）
        destHeight: 300,
        canvasId: 'shareCanvas',
        success: function (res) {
          console.log(res);
          // that.draw_uploadFile(res);
        },
      })
    }, 300)
  },
 
  onShareAppMessage: function () {

  }
})