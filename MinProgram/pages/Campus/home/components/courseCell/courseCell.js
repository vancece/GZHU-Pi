var publicData = require("../../../../../utils/public_data.js")
Component({

  properties: {
    course: {
      type: Object,
      value: {
        "color": 0,
        "weekday": 2,
        "start": 1,
        "last": 3,
        "weeks": "1-16周",
        "course_name": "传入course数据"
      }
    }
  },

  data: {
    colors: publicData.colors
  },

  methods: {
  },

  lifetimes: {
    created: function() {},
    attached: function() {},
  }

})