package models

import (
    "zq/callout_crm/utils"

    "github.com/astaxie/beego/orm"
)

// TableName 设置表名
func (a *SmsRecord) TableName() string {
    return SmsRecordTBName()
}

// SmsRecordQueryParam 用于搜索的类
type SmsRecordQueryParam struct {
    BaseQueryParam
    MobileContent string
    Classify      string
    UIds          []int
    UserCompanyId int
}

// SmsRecord 用户角色 实体类
type SmsRecord struct {
    Id            int
    Mobile        string
    Content       string
    Result        int
    SendTime      string
    Classify      int
    BackendUserId int
    UserCompanyId int
}

// SmsRecordPageList 获取分页数据
func SmsRecordPageList(params *SmsRecordQueryParam) ([]*SmsRecord, int64) {
    query := orm.NewOrm().QueryTable(SmsRecordTBName())
    data := make([]*SmsRecord, 0)
    //默认排序
    sortorder := "Id"
    switch params.Sort {
    case "Id":
        sortorder = "Id"
    }
    if params.Order == "desc" {
        sortorder = "-" + sortorder
    }
    if len(params.UIds) > 0 {
        query = query.Filter("backend_user_id__in", params.UIds)
    }
    if params.UserCompanyId > 0 {
        query = query.Filter("user_company_id", params.UserCompanyId)
    }
    if len(params.MobileContent) > 0 {
        if utils.IsPhoneNumber(params.MobileContent) {
            query = query.Filter("mobile__contains", params.MobileContent)
        } else {
            query = query.Filter("content__icontains", params.MobileContent)
        }
    }
    if len(params.Classify) > 0 {
        query = query.Filter("classify", params.Classify)
    }
    total, _ := query.Count()
    query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
    return data, total
}

// SmsRecordDataList 获取角色列表
func SmsRecordDataList(params *SmsRecordQueryParam) []*SmsRecord {
    params.Limit = -1
    params.Sort = "Id"
    params.Order = "asc"
    data, _ := SmsRecordPageList(params)
    return data
}

// SmsRecordBatchDelete 批量删除
func SmsRecordBatchDelete(ids []int) (int64, error) {
    query := orm.NewOrm().QueryTable(SmsRecordTBName())
    num, err := query.Filter("id__in", ids).Delete()
    return num, err
}

// SmsRecordOne 获取单条
func SmsRecordOne(id int) (*SmsRecord, error) {
    o := orm.NewOrm()
    m := SmsRecord{Id: id}
    err := o.Read(&m)
    if err != nil {
        return nil, err
    }
    return &m, nil
}
