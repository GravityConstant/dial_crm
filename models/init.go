package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// init 初始化
func init() {
	orm.RegisterModel(new(BackendUser), new(Resource), new(Role), new(RoleResourceRel), new(RoleBackendUserRel), new(UserCompany), new(Gateway), new(Agent), new(UserClient), new(UserClientField), new(Asq), new(UserServerHistory), new(BaseMobileLocation), new(GatewayGroup), new(GatewayGroupMap), new(Task), new(TaskDetail), new(SmsSet), new(SmsTemplate), new(SmsRecord), new(GatewayPhone))
}

// TableName 下面是统一的表名管理
func TableName(name string) string {
	prefix := beego.AppConfig.String("db_dt_prefix")
	return prefix + name
}

// BackendUserTBName 获取 BackendUser 对应的表名称
func BackendUserTBName() string {
	return TableName("crm_backend_user")
}

// ResourceTBName 获取 Resource 对应的表名称
func ResourceTBName() string {
	return TableName("crm_resource")
}

// RoleTBName 获取 Role 对应的表名称
func RoleTBName() string {
	return TableName("crm_role")
}

// RoleResourceRelTBName 角色与资源多对多关系表
func RoleResourceRelTBName() string {
	return TableName("crm_role_resource_rel")
}

// RoleBackendUserRelTBName 角色与用户多对多关系表
func RoleBackendUserRelTBName() string {
	return TableName("crm_role_backenduser_rel")
}

// 用户的公司表
func UserCompanyTBName() string {
	return TableName("crm_user_company")
}

// 网关表
func GatewayTBName() string {
	return TableName("call_gateway")
}

// 网关组表
func GatewayGroupTBName() string {
	return TableName("call_gateway_group")
}

// 网关组，网关映射表
func GatewayGroupMapTBName() string {
	return TableName("call_gateway_group_map")
}


// 座席表
func AgentTBName() string {
	return TableName("crm_agent")
}

// 问卷调查表
func AsqTBName() string {
	return TableName("crm_asq")
}

// 客户资料表
func UserClientTBName() string {
	return TableName("crm_user_client")
}

// 自定义字段
func UserClientFieldTBName() string {
	return TableName("crm_user_client_field")
}

// 客户服务记录
func UserServerHistoryTBName() string {
	return TableName("crm_user_server_history")
}

// 通话记录
func CallPgCdrTBName() string {
	return TableName("call_pg_cdr")
}

// 手机区号
func BaseMobileLocationTBName() string {
	return TableName("base_mobile_location")
}

// 任务表
func TaskTBName() string {
	return TableName("crm_task")
}

// 任务细节表
func TaskDetailTBName() string {
	return TableName("crm_task_detail")
}

// 短信设置
func SmsSetTBName() string {
	return TableName("crm_sms_set")
}

// 短信模板
func SmsTemplateTBName() string {
	return TableName("crm_sms_template")
}

// 发送记录
func SmsRecordTBName() string {
	return TableName("crm_sms_record")
}

func GatewayPhoneTBName() string {
	return TableName("crm_gateway_phone")
}