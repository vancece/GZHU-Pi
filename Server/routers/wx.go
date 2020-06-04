package routers

import (
	"GZHU-Pi/env"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/cache"
	"github.com/silenceper/wechat/menu"
	"github.com/silenceper/wechat/message"
	"github.com/silenceper/wechat/util"
	"net/http"
	"strings"
)

var wc *wechat.Wechat
var MinAppID = ""

func GetWC() *wechat.Wechat {
	wxInit()
	return wc
}

//微信公众初始化
func wxInit() (ok bool) {

	if wc != nil {
		return true
	}

	wx := env.Conf.WeiXin
	MinAppID = wx.MinAppID

	if !wx.Enabled || MinAppID == "" {
		logs.Warn("disable weixin")
		return
	}

	//配置微信参数
	config := &wechat.Config{
		AppID:          wx.AppID,
		AppSecret:      wx.Secret,
		Token:          wx.Token,
		EncodingAESKey: wx.AseKey,
		Cache:          cache.NewMemory(),
	}
	wc = wechat.NewWechat(config)

	mu := wc.GetMenu()

	//微信公众菜单
	myMenu := []*menu.Button{
		{
			Type: "click",
			Name: "信息查询",
			Key:  "query",
			SubButtons: []*menu.Button{
				{
					Type:     "miniprogram",
					Name:     "小程序主页",
					Key:      "home",
					URL:      "https://baidu.com",
					AppID:    MinAppID,
					PagePath: "pages/Campus/home/home",
				}, {
					Type:     "miniprogram",
					Name:     "成绩查询",
					Key:      "grade",
					URL:      "https://baidu.com",
					AppID:    MinAppID,
					PagePath: "pages/Campus/grade/grade",
				}, {
					Type:     "miniprogram",
					Name:     "广大校历",
					Key:      "calendar",
					URL:      "https://baidu.com",
					AppID:    MinAppID,
					PagePath: "pages/Campus/tools/calendar",
				}, {
					Type:     "miniprogram",
					Name:     "考试查询",
					Key:      "exam",
					URL:      "https://baidu.com",
					AppID:    MinAppID,
					PagePath: "/pages/Campus/tools/exam",
				}, {
					Type:     "miniprogram",
					Name:     "成绩排行",
					Key:      "rank",
					URL:      "https://baidu.com",
					AppID:    MinAppID,
					PagePath: "pages/Campus/grade/rank",
				},
			},
		}, {
			Type: "click",
			Name: "功能",
			Key:  "function",
			SubButtons: []*menu.Button{
				{
					Type:     "miniprogram",
					Name:     "学业情况",
					Key:      "achieve",
					URL:      "https://baidu.com",
					AppID:    MinAppID,
					PagePath: "pages/Campus/grade/achieve",
				}, {
					Type:     "miniprogram",
					Name:     "图书馆",
					Key:      "library",
					URL:      "https://baidu.com",
					AppID:    MinAppID,
					PagePath: "pages/Campus/library/search",
				}, {
					Type:     "miniprogram",
					Name:     "任意门",
					Key:      "any_door",
					URL:      "https://baidu.com",
					AppID:    MinAppID,
					PagePath: "pages/Campus/course/tools?id=query",
				}, {
					Type:     "miniprogram",
					Name:     "同步中心",
					Key:      "sync",
					URL:      "https://baidu.com",
					AppID:    MinAppID,
					PagePath: "pages/Setting/login/sync",
				}, {
					Type:    "media_id",
					Name:    "联系派派",
					MediaID: "oVb96gPsyuxuaUAhLrub2xqckeMWzoCC5UqwkwGUHLo",
				},
			},
		}, {
			Type: "click",
			Name: "其它",
			Key:  "function",
			SubButtons: []*menu.Button{
				{
					Type: "view",
					Name: "校园全景",
					URL:  "https://720yun.com/t/b8d21qagwni?scene_id=1083548",
				}, {
					Type: "view",
					Name: "失物招领",
					URL:  "http://gzdxzlh3.cn/ssm_wechat/goods/goodsIndex.do",
				}, {
					Type: "view",
					Name: "学号查询",
					URL:  "http://welcome.gzhu.edu.cn/login.portal",
				},
			},
		},
	}

	err := mu.SetMenu(myMenu)
	if err != nil {
		logs.Error(err)
		return
	}
	return true
}

func WxMessage(w http.ResponseWriter, r *http.Request) {

	if !wxInit() {
		return
	}

	server := wc.GetServer(r, w)

	//设置接收消息的处理方法
	server.SetMessageHandler(wxReply)

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		logs.Error(err)
		return
	}
	//发送回复的消息
	err = server.Send()
	if err != nil {
		logs.Error(err)
		return
	}
}

func wxReply(msg message.MixMessage) *message.Reply {

	logs.Info(fmt.Sprintf("收到一条消息：%v", msg))

	switch msg.MsgType {
	//文本消息
	case message.MsgTypeText:

		switch {
		case strings.Contains(msg.Content, "绑定") || strings.Contains(msg.Content, "关联"):
			replyStr := fmt.Sprintf(`<a href="http://www.qq.com" data-miniprogram-appid="%s" data-miniprogram-path="%s?mp_open_id=%s">绑定小程序</a>
回复【上课提醒】试试`,
				env.Conf.WeiXin.MinAppID, mpBindPath, msg.FromUserName)

			return &message.Reply{MsgType: message.MsgTypeText,
				MsgData: message.NewText(replyStr)}

		case strings.Contains(msg.Content, "管理") || strings.Contains(msg.Content, "提醒"):
			replyStr := fmt.Sprintf(`<a href="http://www.qq.com" data-miniprogram-appid="%s" data-miniprogram-path="%s">上课提醒管理</a>
该功能处于内测阶段，绑定公众号后重新同步一次课表会自动添加提醒记录，也可以手动添加提醒，如通知下发失败请联系派派`,
				env.Conf.WeiXin.MinAppID, classNotifyMgrPath)

			return &message.Reply{MsgType: message.MsgTypeText,
				MsgData: message.NewText(replyStr)}
		}

		//图片消息
	case message.MsgTypeImage:
		//do something

		//语音消息
	case message.MsgTypeVoice:
		//do something

		//视频消息
	case message.MsgTypeVideo:
		//do something

		//小视频消息
	case message.MsgTypeShortVideo:
		//do something

		//地理位置消息
	case message.MsgTypeLocation:
		//do something

		//链接消息
	case message.MsgTypeLink:
		//do something

		//事件推送消息
	case message.MsgTypeEvent:
		switch msg.Event {
		//EventSubscribe 订阅
		case message.EventSubscribe:
			//do something

			//取消订阅
		case message.EventUnsubscribe:
			//do something

			//用户已经关注公众号，则微信会将带场景值扫描事件推送给开发者
		case message.EventScan:
			//do something

			// 上报地理位置事件
		case message.EventLocation:
			//do something

			// 点击菜单拉取消息时的事件推送
		case message.EventClick:
			//do something

			// 点击菜单跳转链接时的事件推送
		case message.EventView:
			//do something

			// 扫码推事件的事件推送
		case message.EventScancodePush:
			//do something

			// 扫码推事件且弹出“消息接收中”提示框的事件推送
		case message.EventScancodeWaitmsg:
			//do something

			// 弹出系统拍照发图的事件推送
		case message.EventPicSysphoto:
			//do something

			// 弹出拍照或者相册发图的事件推送
		case message.EventPicPhotoOrAlbum:
			//do something

			// 弹出微信相册发图器的事件推送
		case message.EventPicWeixin:
			//do something

			// 弹出地理位置选择器的事件推送
		case message.EventLocationSelect:
			//do something

		}
	}

	return nil
}

func GetMedia() []byte {
	accessToken, err := wc.GetAccessToken()
	if err != nil {
		logs.Error(err)
	}
	uri := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=%s", accessToken)
	reqMenu := map[string]interface{}{
		"type":   "image",
		"offset": 0,
		"count":  20,
	}
	response, err := util.PostJSON(uri, reqMenu)
	if err != nil {
		logs.Error(err)
	}
	logs.Info(string(response))

	return response
}
