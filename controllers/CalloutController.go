package controllers

import (
	"strings"
	"time"

	"zq/callout_crm/dial"
	"zq/callout_crm/enums"
	"zq/callout_crm/models"
	"zq/callout_crm/utils"

	// "github.com/astaxie/beego/orm"
	"github.com/vma/esl"
)

type CalloutController struct {
	BaseController
}

func (c *CalloutController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DirectDial")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

func (c *CalloutController) DirectDial() {
	tocallphone := strings.TrimSpace(c.GetString("tocallphone"))
	// 检查号码的合法性
	if !utils.IsPhoneNumber(tocallphone) {
		c.jsonResult(enums.JRCodeFailed, "呼叫失败: ", "外呼号码错误！")
	}
	// 获取对于公司的呼叫限制次数
	cm, _ := models.UserCompanyOne(c.curUser.UserCompanyId)
	// 获取今天该号码被呼叫的次数
	params := models.CallPgCdrQueryParam{
		StartStamp: 	   time.Now().Format(enums.BaseDateFormat),
		DestinationNumber: tocallphone,
	}
	calledCount := models.CalledCountByTime(&params)
	// 呼叫次数大于等于公司限制的次数的，不能呼叫，预防骚扰电话
	if calledCount.Times >= cm.LimitDial {
		c.jsonResult(enums.JRCodeFailed, "呼叫失败: ", "已到达限制呼叫次数！")
	}
	// 获取座席信息
	// m, err := models.AgentGatewayOneByUserId(c.curUser.Id)
	// if err != nil {
	// 	if err == orm.ErrNoRows {
	// 		c.jsonResult(enums.JRCodeFailed, "呼叫失败: ", "该用户未配置座席")
	// 	} else if err == orm.ErrMultiRows {
	// 		c.jsonResult(enums.JRCodeFailed, "呼叫失败: ", "获取到多个座席，配置错误，请联系管理员！")
	// 	} else {
	// 		c.jsonResult(enums.JRCodeFailed, "呼叫失败: ", "获取座席失败，请联系管理员！")
	// 	}

	// }
	
	// // 检查各个参数是否合法
	// if !utils.IsExtensionNumber(m.ExtNo) {
	// 	c.jsonResult(enums.JRCodeFailed, "呼叫失败: ", "座席号码错误！")
	// }
	// if !utils.IsPhoneNumber(m.GatewayPhoneNumber) {
	// 	c.jsonResult(enums.JRCodeFailed, "呼叫失败: ", "中继小号错误！")
	// }
	// if !utils.HasValue(m.GatewayUrl) {
	// 	c.jsonResult(enums.JRCodeFailed, "呼叫失败: ", "网关未配置！")
	// }
	// if !utils.IsPhoneNumber(m.OriginationCallerIdNumber) {
	// 	c.jsonResult(enums.JRCodeFailed, "呼叫失败: ", "透传号码错误！")
	// }
	// 如果是手机号码进行处理，座机号码待处理
	// if len(tocallphone) == 11 && strings.Index(tocallphone, "1") == 0 {
	// 	no := tocallphone[:7]
	// 	districtNo := models.GetGatewayAreaCode(m.GatewayName)

	// 	if !models.MobilePhoneSameAsGateway(no, districtNo) {
	// 		// 手机号和中继号不在同一个区域，加0
	// 		tocallphone = `0` + tocallphone
	// 	}
	// }
	
	// 判断该用户是否分配了坐席
	if len(c.curUser.ExtNo) == 0 {
		c.jsonResult(enums.JRCodeFailed, "无坐席，请联系管理员分配坐席！", 0)
	}
	// 更改了，一个坐席对应4个号码，轮询
	// 根据坐席和网关类型选择号码组
	groupPhone, err := dial.CalloutGetTrunkPhones(c.curUser.ExtNo, c.curUser.DefaultTrunk, false)
	if err != nil && strings.Compare(err.Error(), "try") == 0 {
		groupPhone, _ = dial.CalloutGetTrunkPhones(c.curUser.ExtNo, c.curUser.DefaultTrunk, true)
	}
	if len(groupPhone) == 0 {
		c.jsonResult(enums.JRCodeFailed, "呼叫失败: ", "座席未分配网关号码！")
	} else {
		params = models.CallPgCdrQueryParam{
			StartStamp: 	time.Now().Format(enums.BaseDateFormat),
			CallerIdNumber: c.curUser.ExtNo,
		}
		// 中继号各呼出了几次
		ccs := models.CallerCountByTime(&params)
		// 设置最大呼叫量和今天已呼叫的次数
		for _, gp := range groupPhone {
			gp.SetMaxCount(cm.LimitCaller)
			for _, cc := range ccs {
				if strings.TrimSpace(cc.Caller) == gp.TrunkNo {
					gp.SetCount(cc.Times)
					break
				}
			}
		}
	}

	// 选取外呼网关及号码
	co, err := dial.CalloutPickTrunkNo(groupPhone)
	// 如果是手机号码进行处理，座机号码待处理
	if len(tocallphone) == 11 && strings.Index(tocallphone, "1") == 0 {
		no := tocallphone[:7]
		districtNo := models.GetGatewayAreaCode(co.Gateway)

		if !models.MobilePhoneSameAsGateway(no, districtNo) {
			// 手机号和中继号不在同一个区域，加0
			tocallphone = `0` + tocallphone
		}
	}
	if err == nil {
		// 呼叫实例
		handler := &dial.Handler{
			ExtensionNumber: "user/" + co.ExtNo,
			DialplanNumber:  co.TrunkNo,		
			Gateway:         co.Gateway,
			VirtualPhone:    co.TrunkNo,	
			CalloutPhone:    tocallphone,
			ALegState:       make(map[string]esl.EventName),
			BLegState:       make(map[string]esl.EventName),
			Done:            make(chan bool),
			UserId:          c.curUser.Id,
		}
		go dial.Init(handler)
		dialStatus := <-handler.Done
		utils.Info("dial done: %v", dialStatus)
		co.IncCount()

		if handler.Err != nil {
			c.jsonResult(enums.JRCodeSucc, "呼叫失败: " + handler.Err.Error(), -1)
		} else if !dialStatus {
			c.jsonResult(enums.JRCodeSucc, "呼叫未接通", -1)
		}
		c.jsonResult(enums.JRCodeSucc, "呼叫成功", 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "当前号码不能再外呼了", 0)
	}

	// 呼叫实例
	// handler := &dial.Handler{
	// 	ExtensionNumber: "user/" + m.ExtNo,
	// 	DialplanNumber:  m.GatewayPhoneNumber,		
	// 	Gateway:         m.GatewayUrl,
	// 	VirtualPhone:    m.OriginationCallerIdNumber,	
	// 	CalloutPhone:    tocallphone,
	// 	ALegState:       make(map[string]esl.EventName),
	// 	BLegState:       make(map[string]esl.EventName),
	// 	Done:            make(chan bool),
	// 	UserId:          c.curUser.Id,
	// }

	// go dial.Init(handler)
	// dialStatus := <-handler.Done
	// utils.Info("dial done: %v", dialStatus)
	// models.CalloutSuccess(co)

	// if handler.Err != nil {
	// 	c.jsonResult(enums.JRCodeSucc, "呼叫失败: " + handler.Err.Error(), -1)
	// } else if !dialStatus {
	// 	c.jsonResult(enums.JRCodeSucc, "呼叫未接通", -1)
	// }
	// c.jsonResult(enums.JRCodeSucc, "呼叫成功", 0)
}