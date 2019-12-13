package models

import (
	"github.com/astaxie/beego/orm"
)

// TableName 设置表名
func (a *UserCompany) TableName() string {
	return UserCompanyTBName()
}

// UserCompanyQueryParam 用于搜索的类
type UserCompanyQueryParam struct {
	BaseQueryParam
	NameLike string
}

// UserCompany 用户角色 实体类
type UserCompany struct {
	Id        	   int
	Name      	   string
	NameAbbr  	   string
	Created   	   string
	GatewayGroupId int
	LimitDial      int
	LimitCaller    int
}

// UserCompanyPageList 获取分页数据
func UserCompanyPageList(params *UserCompanyQueryParam) ([]*UserCompany, int64) {
	query := orm.NewOrm().QueryTable(UserCompanyTBName())
	data := make([]*UserCompany, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	query = query.Filter("name__istartswith", params.NameLike)
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

// UserCompanyDataList 获取角色列表
func UserCompanyDataList(params *UserCompanyQueryParam) []*UserCompany {
	params.Limit = -1
	params.Sort = "Id"
	params.Order = "asc"
	data, _ := UserCompanyPageList(params)
	return data
}

// UserCompanyBatchDelete 批量删除
func UserCompanyBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(UserCompanyTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// UserCompanyOne 获取单条
func UserCompanyOne(id int) (*UserCompany, error) {
	o := orm.NewOrm()
	m := UserCompany{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
