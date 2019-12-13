package models

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
)

// TableName 设置BackendUser表名
func (a *BackendUser) TableName() string {
	return BackendUserTBName()
}

// BackendUserQueryParam 用于查询的类
type BackendUserQueryParam struct {
	BaseQueryParam
	UserNameLike  string //模糊查询
	RealNameLike  string //模糊查询
	Mobile        string //精确查询
	SearchStatus  string //为空不查询，有值精确查询
	UserCompanyId int
	UIds          []int
}

// BackendUser 实体类
type BackendUser struct {
	Id                 int
	RealName           string `orm:"size(32)"`
	UserName           string `orm:"size(24)"`
	UserPwd            string `json:"-"`
	IsSuper            bool
	Status             int
	Mobile             string                `orm:"size(16)"`
	Email              string                `orm:"size(256)"`
	Avatar             string                `orm:"size(256)"`
	RoleIds            []int                 `orm:"-" form:"RoleIds"`
	RoleBackendUserRel []*RoleBackendUserRel `orm:"reverse(many)"` // 设置一对多的反向关系
	ResourceUrlForList []string              `orm:"-"`
	UIds               []int                 `orm:"-"`
	ExtNo              string                `orm:"-"`
	DefaultTrunk       int                   `orm:"-"`	// 默认使用哪个网关
	LeaderId           int
	UserCompanyId      int
}

// BackendUserPageList 获取分页数据
func BackendUserPageList(params *BackendUserQueryParam) ([]*BackendUser, int64) {
	query := orm.NewOrm().QueryTable(BackendUserTBName())
	data := make([]*BackendUser, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	query = query.Filter("username__istartswith", params.UserNameLike)
	query = query.Filter("realname__istartswith", params.RealNameLike)
	if len(params.Mobile) > 0 {
		query = query.Filter("mobile", params.Mobile)
	}
	if len(params.SearchStatus) > 0 {
		query = query.Filter("status", params.SearchStatus)
	}
	if params.UserCompanyId > 0 {
		query = query.Filter("user_company_id", params.UserCompanyId)
	}
	if len(params.UIds) > 0 {
		query = query.Filter("id__in", params.UIds)
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

//获取用户列表
func BackendUserDataList(params *BackendUserQueryParam) []*BackendUser {
	params.Limit = -1
	params.Sort = "Id"
	params.Order = "asc"
	data, _ := BackendUserPageList(params)
	return data
}

// BackendUserOne 根据id获取单条
func BackendUserOne(id int) (*BackendUser, error) {
	o := orm.NewOrm()
	m := BackendUser{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// BackendUserOneByUserName 根据用户名密码获取单条
func BackendUserOneByUserName(username, userpwd string) (*BackendUser, error) {
	m := BackendUser{}
	err := orm.NewOrm().QueryTable(BackendUserTBName()).Filter("username", username).Filter("userpwd", userpwd).One(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// 根据backend_user的id递归查询其下子项
func BackendUserRecursiveById(id int) []*BackendUser {
	data := make([]*BackendUser, 0)

	o := orm.NewOrm()
	sql := `with RECURSIVE le (id, real_name, user_name) as ( 
            select id, real_name, user_name from %s where id=? 
            UNION all 
            select u1.id, u1.real_name, u1.user_name from %s u1, le u2 where u1.leader_id=u2.id 
        ) 
        select * from le`
	sql = strings.Replace(sql, "\n", "", -1)
	sql = strings.Replace(sql, "    ", "", -1)
	sql = fmt.Sprintf(sql, BackendUserTBName(), BackendUserTBName())

	o.Raw(sql, id).QueryRows(&data)

	return data
}

type BackendUserCompany struct {
	Id          int
	RealName    string
	UserName    string
	CompanyName string
	GatewayId   int
}

// 根据backend_user_id查company
func BackendUserCompanyByUserId(id int) (*BackendUserCompany, error) {
	m := BackendUserCompany{}
	o := orm.NewOrm()
	sql := `select u.id, u.real_name, u.user_name, c.name as company_name, c.gateway_id from %s u left join %s c on u.user_company_id=c.id where u.id=?`
	sql = fmt.Sprintf(sql, BackendUserTBName(), UserCompanyTBName())

	if err := o.Raw(sql, id).QueryRow(&m); err != nil {
		return nil, err
	}
	return &m, nil
}
