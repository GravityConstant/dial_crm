package models

import (
	"github.com/astaxie/beego/orm"
)

// TableName 设置表名
func (a *Role) TableName() string {
	return RoleTBName()
}

// RoleQueryParam 用于搜索的类
type RoleQueryParam struct {
	BaseQueryParam
	NameLike string
	UComId 	 int
	RoleIds  []int
}

// Role 用户角色 实体类
type Role struct {
	Id                 int    `form:"Id"`
	Name               string `form:"Name"`
	Seq                int
	UserCompanyId      int
	RoleResourceRel    []*RoleResourceRel    `orm:"reverse(many)" json:"-"` // 设置一对多的反向关系
	RoleBackendUserRel []*RoleBackendUserRel `orm:"reverse(many)" json:"-"` // 设置一对多的反向关系
}

// RolePageList 获取分页数据
func RolePageList(params *RoleQueryParam) ([]*Role, int64) {
	query := orm.NewOrm().QueryTable(RoleTBName())
	data := make([]*Role, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	case "Seq":
		sortorder = "Seq"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	if params.UComId > 0 {
		query = query.Filter("user_company_id", params.UComId)
	}
    if len(params.RoleIds) > 0 {
        query = query.Exclude("id__in", params.RoleIds)
    }
	query = query.Filter("name__istartswith", params.NameLike)
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

// RoleDataList 获取角色列表
func RoleDataList(params *RoleQueryParam) []*Role {
	params.Limit = -1
	params.Sort = "Seq"
	params.Order = "asc"
	data, _ := RolePageList(params)
	return data
}

// RoleBatchDelete 批量删除
func RoleBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(RoleTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// RoleOne 获取单条
func RoleOne(id int) (*Role, error) {
	o := orm.NewOrm()
	m := Role{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
