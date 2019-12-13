package models

import (
	"strconv"

	"github.com/astaxie/beego/orm"
)

func (a *Gateway) TableName() string {
	return GatewayTBName()
}

type GatewayQueryParam struct {
	BaseQueryParam
	NameLike string
}

type Gateway struct {
	Id          	int
	GatewayName 	string
	GatewayUrl  	string
	MaxCall     	int
	GatewayGroupIds []int              `orm:"-" form:"GatewayGroupIds"`
	GatewayGroupMap []*GatewayGroupMap `orm:"reverse(many)" json:"-"`
	GatewayType     int
}

func GatewayPageList(params *GatewayQueryParam) ([]*Gateway, int64) {
	query := orm.NewOrm().QueryTable(GatewayTBName())
	data := make([]*Gateway, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	query = query.Filter("gateway_name__icontains", params.NameLike)
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

// 根据id获取单条
func GatewayOne(id int) (*Gateway, error) {
	o := orm.NewOrm()
	m := Gateway{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// DataList 获取角色列表
func GatewayDataList(params *GatewayQueryParam) []*Gateway {
	params.Limit = -1
	params.Sort = "Id"
	params.Order = "asc"
	data, _ := GatewayPageList(params)
	return data
}

// DataList 获取角色列表
func DataListByUserCompany(uId int) []*Gateway {
	data := make([]*Gateway, 0)
	o := orm.NewOrm()
	var sql string

	if uId > 0 {
		sql = "SELECT id, gateway_name, gateway_url from call_gateway where id in (" + 
			    "SELECT gateway_id from call_gateway_group_map where gateway_group_id = (" + 
				  "SELECT gateway_group_id from crm_user_company where id=(" + 
				    "SELECT user_company_id from crm_backend_user where id=" + strconv.Itoa(uId) + 
				  ")" + 
				")" + 
			  ") order by id asc"
	} else {
		sql = "SELECT id, gateway_name, gateway_url from call_gateway order by id asc"
	}
	
	o.Raw(sql).QueryRows(&data)

	return data
}

// BatchDelete 批量删除
func GatewayBatchDelete(ids []int) (int64, error) {
    query := orm.NewOrm().QueryTable(GatewayTBName())
    num, err := query.Filter("id__in", ids).Delete()
    return num, err
}

// 获取GatewayName对应的区号
func GetGatewayAreaCode(gatewayName string) string {
	switch gatewayName {
	case "fuzhou":
		return "0591"
	case "xiamen":
		return "0592"
	case "ningde":
		return "0593"
	case "putian":
		return "0594"
	case "quanzhou":
		return "0595"
	case "zhangzhou":
		return "0596"
	case "longyan":
		return "0597"
	case "sanming":
		return "0598"
	case "nanping":
		return "0599"
	default:
		return "0591"
	}
}

// 座席号码取自freeswitch的配置文件/usr/local/freeswitch/conf/sip_profiles/external下的文件名
type FsRegisterGateway struct {
	Name string
}