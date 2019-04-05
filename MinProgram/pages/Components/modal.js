Component({

  properties: {
    //是否显示modal
    show: {
      type: Boolean,
      value: false
    },
    //自定义模式
    custom: {
      type: Boolean,
      value: false
    },
    // 启用确认按键
    confirm: {
      type: Boolean,
      value: true
    },
    // 启用取消按键
    cancel: {
      type: Boolean,
      value: true
    },
    confirmText: {
      type: String,
      value: "确定"
    },
    cancelText: {
      type: String,
      value: "取消"
    },
    title: {
      type: String,
      value: "标题"
    },
    // 禁止关闭
    hideCancel: {
      type: Boolean,
      value: false
    },
    // 宽度
    width: {
      type: String,
      value: "70%"
    }
  },

  data: {

  },

  methods: {
    clickMask() {
      // this.setData({show: false})
    },

    cancel() {
      this.setData({
        show: false
      })
      this.triggerEvent('cancel')
    },

    confirm() {
      this.setData({
        show: false
      })
      this.triggerEvent('confirm')
    }
  }
})