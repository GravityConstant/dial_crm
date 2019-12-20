package sysinit

import (
	"zq/callout_crm/utils"
	"zq/callout_crm/dial"
	"zq/callout_crm/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//启用Session
	beego.BConfig.WebConfig.Session.SessionOn = true
	//初始化日志
	utils.InitLogs()
	//初始化缓存
	utils.InitCache()
	//初始化数据库
	InitDatabase()
	// 设置静态路径
	setStaticPath()
	// 外呼初始化
	dial.CalloutInit()
	// 呼入初始化
	dial.CallinInit()
	// 所有连接的集合
	controllers.CallInHubInit()
}

func setStaticPath() {
	// 录音路径
	beego.SetStaticPath("/recordpath", "/home/voices/records")
}