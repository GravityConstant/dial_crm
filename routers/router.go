package routers

import (
	"zq/callout_crm/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//用户角色路由
	beego.Router("/role/index", &controllers.RoleController{}, "*:Index")
	beego.Router("/role/datagrid", &controllers.RoleController{}, "Get,Post:DataGrid")
	beego.Router("/role/edit/?:id", &controllers.RoleController{}, "Get,Post:Edit")
	beego.Router("/role/delete", &controllers.RoleController{}, "Post:Delete")
	beego.Router("/role/datalist", &controllers.RoleController{}, "Post:DataList")
	beego.Router("/role/allocate", &controllers.RoleController{}, "Post:Allocate")
	beego.Router("/role/updateseq", &controllers.RoleController{}, "Post:UpdateSeq")

	//资源路由
	beego.Router("/resource/index", &controllers.ResourceController{}, "*:Index")
	beego.Router("/resource/treegrid", &controllers.ResourceController{}, "POST:TreeGrid")
	beego.Router("/resource/edit/?:id", &controllers.ResourceController{}, "Get,Post:Edit")
	beego.Router("/resource/parent", &controllers.ResourceController{}, "Post:ParentTreeGrid")
	beego.Router("/resource/delete", &controllers.ResourceController{}, "Post:Delete")
	//快速修改顺序
	beego.Router("/resource/updateseq", &controllers.ResourceController{}, "Post:UpdateSeq")

	//通用选择面板
	beego.Router("/resource/select", &controllers.ResourceController{}, "Get:Select")
	//用户有权管理的菜单列表（包括区域）
	beego.Router("/resource/usermenutree", &controllers.ResourceController{}, "POST:UserMenuTree")
	beego.Router("/resource/checkurlfor", &controllers.ResourceController{}, "POST:CheckUrlFor")

	//后台用户路由
	beego.Router("/backenduser/index", &controllers.BackendUserController{}, "*:Index")
	beego.Router("/backenduser/datagrid", &controllers.BackendUserController{}, "POST:DataGrid")
	beego.Router("/backenduser/datalist", &controllers.BackendUserController{}, "Post:DataList")
	beego.Router("/backenduser/edit/?:id", &controllers.BackendUserController{}, "Get,Post:Edit")
	beego.Router("/backenduser/delete", &controllers.BackendUserController{}, "Post:Delete")
	//后台用户中心
	beego.Router("/usercenter/profile", &controllers.UserCenterController{}, "Get:Profile")
	beego.Router("/usercenter/basicinfosave", &controllers.UserCenterController{}, "Post:BasicInfoSave")
	beego.Router("/usercenter/uploadimage", &controllers.UserCenterController{}, "Post:UploadImage")
	beego.Router("/usercenter/passwordsave", &controllers.UserCenterController{}, "Post:PasswordSave")

	beego.Router("/home/index", &controllers.HomeController{}, "*:Index")
	beego.Router("/home/login", &controllers.HomeController{}, "*:Login")
	beego.Router("/home/dologin", &controllers.HomeController{}, "Post:DoLogin")
	beego.Router("/home/logout", &controllers.HomeController{}, "*:Logout")

	beego.Router("/home/404", &controllers.HomeController{}, "*:Page404")
	beego.Router("/home/error/?:error", &controllers.HomeController{}, "*:Error")

	beego.Router("/", &controllers.HomeController{}, "*:Index")

	// 快速拨打电话
	beego.Router("/callout/directdial", &controllers.CalloutController{}, "Post:DirectDial")

	// 用户公司
	beego.Router("/usercompany/index", &controllers.UserCompanyController{}, "*:Index")
	beego.Router("/usercompany/datagrid", &controllers.UserCompanyController{}, "Post:DataGrid")
	beego.Router("/usercompany/datalist", &controllers.UserCompanyController{}, "Post:DataList")
	beego.Router("/usercompany/edit/?:id", &controllers.UserCompanyController{}, "Get,Post:Edit")
	beego.Router("/usercompany/delete", &controllers.UserCompanyController{}, "Post:Delete")

	// 座席管理
	beego.Router("/agent/index", &controllers.AgentController{}, "*:Index")
	beego.Router("/agent/datagrid", &controllers.AgentController{}, "Post:DataGrid")
	beego.Router("/agent/datalist", &controllers.AgentController{}, "Post:DataList")
	beego.Router("/agent/edit/?:id", &controllers.AgentController{}, "Get,Post:Edit")
	beego.Router("/agent/delete", &controllers.AgentController{}, "Post:Delete")
	beego.Router("/agent/fsregisteruserlist", &controllers.AgentController{}, "Post:FsRegisterUserList")
	beego.Router("/agent/editv1/?:id", &controllers.AgentController{}, "Get,Post:EditV1")

	// 客户的客户表
	beego.Router("/userclient/index", &controllers.UserClientController{}, "*:Index")
	beego.Router("/userclient/datagrid", &controllers.UserClientController{}, "Post:DataGrid")
	beego.Router("/userclient/datalist", &controllers.UserClientController{}, "Post:DataList")
	beego.Router("/userclient/edit/?:id/?:phone/?:dialSuccess", &controllers.UserClientController{}, "Get,Post:Edit")
	beego.Router("/userclient/delete", &controllers.UserClientController{}, "Post:Delete")
	beego.Router("/userclient/updatefield", &controllers.UserClientController{}, "Post:UpdateFieldByIds")
	beego.Router("/userclient/downloadtmpl", &controllers.UserClientController{}, "Post:DownloadTmpl")
	beego.Router("/userclient/uploadexcel", &controllers.UserClientController{}, "Post:UploadExcel")
	beego.Router("/userclient/download", &controllers.UserClientController{}, "Post:Download")
	beego.Router("/userclient/isexist", &controllers.UserClientController{}, "Post:IsExistedCalledPhone")

	// 客户的客户表对应的字段表
	beego.Router("/userclientfield/index", &controllers.UserClientFieldController{}, "*:Index")
	beego.Router("/userclientfield/datagrid", &controllers.UserClientFieldController{}, "Post:DataGrid")
	beego.Router("/userclientfield/datalist", &controllers.UserClientFieldController{}, "Post:DataList")
	beego.Router("/userclientfield/edit/?:id", &controllers.UserClientFieldController{}, "Get,Post:Edit")
	beego.Router("/userclientfield/delete", &controllers.UserClientFieldController{}, "Post:Delete")
	beego.Router("/userclientfield/updatefield", &controllers.UserClientFieldController{}, "Post:UpdateField")

	//问卷调查
	beego.Router("/asq/index", &controllers.AsqController{}, "*:Index")
	beego.Router("/asq/datagrid", &controllers.AsqController{}, "Post:DataGrid")
	beego.Router("/asq/edit/?:id", &controllers.AsqController{}, "Get,Post:Edit")
	beego.Router("/asq/delete", &controllers.AsqController{}, "Post:Delete")
	beego.Router("/asq/set/?:id", &controllers.AsqController{}, "Get,Post:Set")

	//客服服务记录
	beego.Router("/userserverhistory/addrecord/?:id/?:state", &controllers.UserServerHistoryController{}, "Get,Post:AddRecord")
	beego.Router("/userserverhistory/recordlist/?:id", &controllers.UserServerHistoryController{}, "Get,Post:RecordList")
	beego.Router("/userserverhistory/recordlistdatagrid", &controllers.UserServerHistoryController{}, "Post:RecordListDataGrid")

	// 通话记录
	beego.Router("/callpgcdr/index/", &controllers.CallPgCdrController{}, "*:Index")
	beego.Router("/callpgcdr/datagrid", &controllers.CallPgCdrController{}, "Get,Post:DataGrid")
	beego.Router("/callpgcdr/downloadcdr", &controllers.CallPgCdrController{}, "Post:DownloadCdr")
	beego.Router("/callpgcdr/downloadrecord", &controllers.CallPgCdrController{}, "Post:DownloadRecord")

	// 网关
	beego.Router("/gateway/index", &controllers.GatewayController{}, "*:Index")
	beego.Router("/gateway/datagrid", &controllers.GatewayController{}, "Get,Post:DataGrid")
	beego.Router("/gateway/edit/?:id", &controllers.GatewayController{}, "Get,Post:Edit")
	beego.Router("/gateway/delete", &controllers.GatewayController{}, "Post:Delete")
	beego.Router("/gateway/datalist", &controllers.GatewayController{}, "Post:DataList")
	beego.Router("/gateway/fsgwlst", &controllers.GatewayController{}, "Post:FsRegisterGatewayList")
	beego.Router("/gateway/datalistbyuc", &controllers.GatewayController{}, "Post:DataListByUserCompany")

	// 网关组
	beego.Router("/gatewaygroup/index", &controllers.GatewayGroupController{}, "*:Index")
	beego.Router("/gatewaygroup/datagrid", &controllers.GatewayGroupController{}, "Get,Post:DataGrid")
	beego.Router("/gatewaygroup/edit/?:id", &controllers.GatewayGroupController{}, "Get,Post:Edit")
	beego.Router("/gatewaygroup/delete", &controllers.GatewayGroupController{}, "Post:Delete")
	beego.Router("/gatewaygroup/datalist", &controllers.GatewayGroupController{}, "Post:DataList")

	// 任务表
	beego.Router("/task/index", &controllers.TaskController{}, "*:Index")
	beego.Router("/task/datagrid", &controllers.TaskController{}, "Post:DataGrid")
	beego.Router("/task/datalist", &controllers.TaskController{}, "Post:DataList")
	beego.Router("/task/edit/?:id", &controllers.TaskController{}, "Get,Post:Edit")
	beego.Router("/task/delete", &controllers.TaskController{}, "Post:Delete")
	beego.Router("/task/updatefield", &controllers.TaskController{}, "Post:UpdateFieldByIds")
	beego.Router("/task/importindex/?:id", &controllers.TaskController{}, "Get,Post:ImportUserClientIndex")

	// 任务细节表
	beego.Router("/taskdetail/assign", &controllers.TaskDetailController{}, "Post:Assign")
	beego.Router("/taskdetail/datagrid", &controllers.TaskDetailController{}, "Post:DataGrid")
	beego.Router("/taskdetail/averageassign", &controllers.TaskDetailController{}, "Post:AverageAssign")

	// 我的任务
	beego.Router("/mytask/index", &controllers.MyTaskController{}, "*:Index")
	beego.Router("/mytask/datagrid", &controllers.MyTaskController{}, "Post:DataGrid")
	beego.Router("/mytask/updatecallstate", &controllers.MyTaskController{}, "Post:UpdateCallState")

	// 短信设置
	beego.Router("/smsset/index", &controllers.SmsSetController{}, "*:Index")
	beego.Router("/smsset/datagrid", &controllers.SmsSetController{}, "Post:DataGrid")
	beego.Router("/smsset/datalist", &controllers.SmsSetController{}, "Post:DataList")
	beego.Router("/smsset/edit/?:id", &controllers.SmsSetController{}, "Get,Post:Edit")
	beego.Router("/smsset/delete", &controllers.SmsSetController{}, "Post:Delete")

	// 短信模板
	beego.Router("/smstmpl/datagrid", &controllers.SmsTemplateController{}, "Post:DataGrid")
	beego.Router("/smstmpl/datalist", &controllers.SmsTemplateController{}, "Post:DataList")
	beego.Router("/smstmpl/edit/?:id", &controllers.SmsTemplateController{}, "Get,Post:Edit")
	beego.Router("/smstmpl/delete", &controllers.SmsTemplateController{}, "Post:Delete")

	// 短信模板
	beego.Router("/smsrcd/datagrid", &controllers.SmsRecordController{}, "Post:DataGrid")
	beego.Router("/smsrcd/datalist", &controllers.SmsRecordController{}, "Post:DataList")
	beego.Router("/smsrcd/edit/?:id", &controllers.SmsRecordController{}, "Get,Post:Edit")
	beego.Router("/smsrcd/delete", &controllers.SmsRecordController{}, "Post:Delete")
	beego.Router("/smsrcd/sendmsg", &controllers.SmsRecordController{}, "Post:SendMsg")

	// 网关号码
	beego.Router("/gatewayphone/index", &controllers.GatewayPhoneController{}, "*:Index")
	beego.Router("/gatewayphone/datagrid", &controllers.GatewayPhoneController{}, "Post:DataGrid")
	beego.Router("/gatewayphone/datalist", &controllers.GatewayPhoneController{}, "Post:DataList")
	beego.Router("/gatewayphone/edit/?:id", &controllers.GatewayPhoneController{}, "Get,Post:Edit")
	beego.Router("/gatewayphone/delete", &controllers.GatewayPhoneController{}, "Post:Delete")
	beego.Router("/gatewayphone/allocateagent", &controllers.GatewayPhoneController{}, "Post:AllocateAgent")

	// websocket
	beego.Router("/ws/echo", &controllers.WebSocketController{}, "*:Echo")
}
