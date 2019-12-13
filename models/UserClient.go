package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

// TableName 设置表名
func (a *UserClient) TableName() string {
	return UserClientTBName()
}

// UserClientQueryParam 用于搜索的类
type UserClientQueryParam struct {
	BaseQueryParam
	NameLike       string
	PhoneLike      string
	CreatedStart   string
	CreatedEnd     string
	BackendUserIds []int
	CompanyId      int
	QueryShow      bool
}

// UserClient 用户角色 实体类
type UserClient struct {
	Id                  int
	Name                string
	MobilePhone         string
	ContactPhone        string
	BackendUserId       int    // -1
	Created             string // -1
	Comment             string // 备注
	Address             string
	DialState           string // 接通状态
	State               int    // 1：无意向，2：中等意向，3：有意向
	Feature             int    // 客户特征id,其实应该还有一张对应的特征表；这里作枚举类型
	Complaint           string // 填具体投诉的内容
	LatestCommunicated  string // 最近沟通时间
	ClueFrom            int    // 线索来源。1：百度推广，2：QQ咨询，3：其他
	Email               string // 0
	BelongBackendUserId int    // -1
	UserCompanyId       int    // -1
	Updated             string // -1
	Column1             string
	Column2             string
	Column3             string
	Column4             string
	Column5             string
	Column6             string
	Column7             string
	Column8             string
	Column9             string
	Column10            string
	Column11            string
	Column12            string
	Column13            string
	Column14            string
	Column15            string
	Column16            string
}

// UserClientPageList 获取分页数据
func UserClientPageList(params *UserClientQueryParam) ([]*UserClient, int64) {
	query := orm.NewOrm().QueryTable(UserClientTBName())
	cond := orm.NewCondition()
	data := make([]*UserClient, 0)
	//默认排序
	sortorder := "Updated"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	case "MobilePhone":
		sortorder = "MobilePhone"
	case "ContactPhone":
		sortorder = "ContactPhone"
	case "Updated":
		sortorder = "Updated"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}

	if len(params.PhoneLike) > 0 {
		cond = cond.Or("mobile_phone__contains", params.PhoneLike).Or("contact_phone__contains", params.PhoneLike)
	}
	if len(params.NameLike) > 0 {
		cond = cond.And("name__contains", params.NameLike)
	}
	if len(params.CreatedStart) > 0 {
		cond = cond.And("created__gte", params.CreatedStart)
	}
	if len(params.CreatedEnd) > 0 {
		cond = cond.And("created__lte", params.CreatedEnd)
	}
	if params.CompanyId > 0 {
		cond = cond.And("user_company_id", params.CompanyId)
	}
	// 归属人字段给了belong_backend_user_id
	if len(params.BackendUserIds) > 0 {
		cond = cond.And("belong_backend_user_id__in", params.BackendUserIds, 0)
	}
	query = query.SetCond(cond)
	total, _ := query.Count()

	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)

	return data, total
}

// UserClientDataList 获取角色列表
func UserClientDataList(params *UserClientQueryParam) []*UserClient {
	params.Limit = -1
	params.Sort = "Id"
	params.Order = "asc"
	data, _ := UserClientPageList(params)
	return data
}

// UserClientBatchDelete 批量删除
func UserClientBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(UserClientTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// UserClientOne 获取单条
func UserClientOne(id int) (*UserClient, error) {
	o := orm.NewOrm()
	m := UserClient{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

type MpList struct {
	Id          int
	MobilePhone string
}

// 获取电话号码数组
func GetMobilePhoneList(ucId int) []*MpList {
	o := orm.NewOrm()
	data := make([]*MpList, 0)

	sql := `select id, mobile_phone from %s where user_company_id=?`
	sql = fmt.Sprintf(sql, UserClientTBName())
	o.Raw(sql, ucId).QueryRows(&data)

	return data
}