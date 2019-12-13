package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func (a *UserServerHistory) TableName() string {
	return UserServerHistoryTBName()
}

type UserServerHistoryQueryParam struct {
	BaseQueryParam
	UserClientId    int
}


//Asq 实体类
type UserServerHistory struct {
	Id              int
	UserClientId    int
	Created         time.Time `orm:"auto_now_add;type(datetime)"`
	BackendUserId   int
	BackendUserName string
	Context         string
}



// UserServerHistory 获取分页数据
func UserServerHistoryList(params *UserServerHistoryQueryParam) ([]*UserServerHistory, int64) {
	query := orm.NewOrm().QueryTable(UserServerHistoryTBName())
	data := make([]*UserServerHistory, 0)
	var total int64

	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}

	if params.UserClientId == 0 {
		return data, total
	}
	query = query.Filter("UserClientId", params.UserClientId)
	total, _ = query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}
