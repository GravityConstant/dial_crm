package models

import (
    "github.com/astaxie/beego/orm"
)

// TableName 设置表名
func (a *SmsTemplate) TableName() string {
    return SmsTemplateTBName()
}

// SmsTemplateQueryParam 用于搜索的类
type SmsTemplateQueryParam struct {
    BaseQueryParam
    TitleLike     string
    UserCompanyId int
}

// SmsTemplate 实体类
type SmsTemplate struct {
    Id            int
    Title         string
    Content       string
    Classify      int
    BackendUserId int
    UserCompanyId int
    Created       string
    Updated       string
    State         int
}

// SmsTemplatePageList 获取分页数据
func SmsTemplatePageList(params *SmsTemplateQueryParam) ([]*SmsTemplate, int64) {
    query := orm.NewOrm().QueryTable(SmsTemplateTBName())
    data := make([]*SmsTemplate, 0)
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
    if params.UserCompanyId > 0 {
        query = query.Filter("user_company_id", params.UserCompanyId)
    }
    if len(params.TitleLike) > 0 {
        query = query.Filter("title__icontains", params.TitleLike)
    }
    
    total, _ := query.Count()
    query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
    return data, total
}

// SmsTemplateDataList 获取角色列表
func SmsTemplateDataList(params *SmsTemplateQueryParam) []*SmsTemplate {
    params.Limit = -1
    params.Sort = "Id"
    params.Order = "asc"
    data, _ := SmsTemplatePageList(params)
    return data
}

// SmsTemplateBatchDelete 批量删除
func SmsTemplateBatchDelete(ids []int) (int64, error) {
    query := orm.NewOrm().QueryTable(SmsTemplateTBName())
    num, err := query.Filter("id__in", ids).Delete()
    return num, err
}

// SmsTemplateOne 获取单条
func SmsTemplateOne(id int) (*SmsTemplate, error) {
    o := orm.NewOrm()
    m := SmsTemplate{Id: id}
    err := o.Read(&m)
    if err != nil {
        return nil, err
    }
    return &m, nil
}
