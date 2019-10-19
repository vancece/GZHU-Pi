package gzhu_jw

import (
	"fmt"
	"net/http"
)

var (
	LoginError = fmt.Errorf("认证失败，账号或密码错误")
	AuthError  = fmt.Errorf("认证失败，可能是缓存过期，请重试")

	SemCode = []string{"3", "12"} //3是第一学期，12是第二学期
	Year    = "2018"

	jsonHeader       = http.Header{"Content-Type": []string{"application/json"}}
	urlencodedHeader = http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}
)

const baseUrl=""

var Urls = map[string]string{

	"jw-login":     "https://cas.gzhu.edu.cn/cas_server/login?service=http://jwxt.gzhu.edu.cn/sso/lyiotlogin",
	"info":         "http://jwxt.gzhu.edu.cn/jwglxt/xsxxxggl/xsgrxxwh_cxXsgrxx.html?gnmkdm=N100801&layout=default",
	"course":       "http://jwxt.gzhu.edu.cn/jwglxt/kbcx/xskbcx_cxXsKb.html?gnmkdm=N2151",
	"grade":        "http://jwxt.gzhu.edu.cn/jwglxt/cjcx/cjcx_cxDgXscj.html?doType=query&gnmkdm=N100801",
	"exam":         "http://jwxt.gzhu.edu.cn/jwglxt/kwgl/kscx_cxXsksxxIndex.html?doType=query&gnmkdm=N358105",
	"id-credit":    "http://jwxt.gzhu.edu.cn/jwglxt/xsxxxggl/xsxxwh_cxXsxkxx.html?gnmkdm=N100801",
	"empty-room":   "http://jwxt.gzhu.edu.cn/jwglxt/cdjy/cdjy_cxKxcdlb.html?doType=query&gnmkdm=N2155",
	"all-course":   "http://jwxt.gzhu.edu.cn/jwglxt/design/funcData_cxFuncDataList.html?func_widget_guid=DA1B5BB30E1F4CB99D1F6F526537777B&gnmkdm=N219904",
	"second-login": "https://cas.gzhu.edu.cn/cas_server/login?service=http://172.17.1.123/Login.aspx",
	"second-my":    "http://172.17.1.123/XS/XMSB.aspx",
	"second-all":    "http://172.17.1.123/JWC/view.aspx",
}

var WeekdayMatcher = map[string]int{
	"星期一": 1,
	"星期二": 2,
	"星期三": 3,
	"星期四": 4,
	"星期五": 5,
	"星期六": 6,
	"星期日": 7,
	"星期天": 7,
}
