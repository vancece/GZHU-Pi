var wxCharts = require('../../../utils/sdk/wxcharts.js');
var pieChart = null
var pieChart1 = null
var pieChart2 = null
Page({
  data: {

    canvas1: [{
      name: '大一',
      data: 5650,
    }, {
      name: '大二',
      data: 5610,
    }, {
      name: '大三',
      data: 4388,
    }, {
      name: '大四+',
      data: 3177,
    }, {
      name: '其它',
      data: 4700,
    }],

    canvas2: [{
      name: '男生',
      data: 11500,
    }, {
      name: '女生',
      data: 12500,
    }],

    canvas3: [{
      name: '男生',
      data: 4200,
    }, {
      name: '女生',
      data: 5820,
    }],
  },

  onLoad: function(e) {
    wx.cloud.callFunction({
      // 需调用的云函数名
      name: 'analysis',
      // 传给云函数的参数
      // data: {
      //   a: 12,
      //   b: 19,
      // },
      // 成功回调
      complete: function (res) {

        console.log(res.result)
      }
    })
    this.drawCover()
    this.data.chart1 = this.draw("canvas1", "ring", this.data.canvas1)
    this.data.chart2 = this.draw("canvas2", "ring", this.data.canvas2)
    this.data.chart3 = this.draw("canvas3", "pie", this.data.canvas3)
  },

  draw(canvasId, type, data) {
    var chart = new wxCharts({
      animation: true,
      canvasId: canvasId,
      type: type,
      series: data,
      width: wx.getSystemInfoSync().windowWidth,
      height: 300,
      dataLabel: true,
    });
    return chart
  },

  drawCover() {
    new wxCharts({
      canvasId: 'canvas0',
      type: 'column',
      categories: ['大一', '大二', '大三', '大四+'],
      series: [{
        name: '招生数',
        data: [7619, 7612, 7560, 7200, 0]
      }, {
        name: '用户数',
        data: [5650, 5610, 4388, 3177, 0]
      }],
      yAxis: {
        format: function(val) {
          var ret = Math.ceil((val / 7700) * 100)
          return ret + '%';
        },
      },
      width: wx.getSystemInfoSync().windowWidth * 0.9,
      height: 200,
    });
  },

  touchHandler: function(e) {
    let id = 0
    let data
    switch (e.target.id) {
      case "canvas0":
        return
      case "canvas1":
        id = this.data.chart1.getCurrentDataIndex(e)
        data = this.data.canvas1[id]
        break
      case "canvas2":
        id = this.data.chart2.getCurrentDataIndex(e)
        data = this.data.canvas2[id]
        break
      case "canvas3":
        id = this.data.chart3.getCurrentDataIndex(e)
        data = this.data.canvas3[id]
        break
      default:
        break
    }
    if (data == undefined) return
    let msg = data.name + "：" + data.data
    wx.showToast({
      title: msg,
      icon: "none"
    })
  },

  preview() {
    let imgurl = "https://cos.ifeel.vip/gzhu-pi/images/resource/qrcode.jpg"
    wx.previewImage({
      urls: [imgurl]
    })
  }

});