package enums

type JsonResultCode int

const (
	JRCodeSucc JsonResultCode = iota
	JRCodeFailed
	JRCode302 = 302 //跳转至地址
	JRCode401 = 401 //未授权访问
)

const (
	Deleted = iota - 1
	Disabled
	Enabled
)

const BaseTimeFormat = "2006-01-02 15:04:05"
const BaseDateFormat = "2006-01-02"

const MaxSelfDefineColumn = 16 // crm_user_client表最大只能定义16个字段
