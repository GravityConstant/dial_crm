package models

import (
    "fmt"
    "strings"

    "zq/callout_crm/utils"

    "github.com/astaxie/beego/orm"
)

// TableName 设置表名
func (a *SmsSet) TableName() string {
    return SmsSetTBName()
}

// SmsSetQueryParam 用于搜索的类
type SmsSetQueryParam struct {
    BaseQueryParam
    NameLike      string
    UserCompanyId int
}

// SmsSet 用户角色 实体类
type SmsSet struct {
    Id            int
    Account       string
    Password      string
    Signature     string
    BackendUserId int
    UserCompanyId int
    Created       string
    Updated       string
    Classify      int
}

type SmsSetWithRelated struct {
    Id              int
    Account         string
    Password        string
    Signature       string
    BackendUserId   int
    UserCompanyId   int
    UserCompanyName string
    Created         string
    Updated         string
}

// SmsSetPageList 获取分页数据
func SmsSetPageList(params *SmsSetQueryParam) ([]*SmsSetWithRelated, int64) {
    o := orm.NewOrm()
    data := make([]*SmsSetWithRelated, 0)
    //默认排序
    sortorder := "Id"
    switch params.Sort {
    case "Id":
        sortorder = "Id"
    }

    var sql string

    sqlCol := `SELECT s.id, s.account, s.password, s.signature, s.backend_user_id, s.user_company_id, `+
              `c.name as user_company_name, s.created, s.updated `

    sqlFrom := `FROM %s s INNER JOIN %s c on s.user_company_id=c.id `
    sqlFrom = fmt.Sprintf(sqlFrom, SmsSetTBName(), UserCompanyTBName())
    // where
    var sqlWhere string
    andCond := make([]string, 0)
    if len(params.NameLike) > 0 {
        andCond = append(andCond, fmt.Sprintf(`c.name ilike '%%%s%%'`, params.NameLike))
    }
    if params.UserCompanyId > 0 {
        andCond = append(andCond, fmt.Sprintf(`s.user_company_id = %d`, params.UserCompanyId))
    }
    if len(andCond) > 0 {
        sqlWhere = "where " + strings.Join(andCond, " and ")
    }
    // order by
    var orderBy string
    if len(sortorder) > 0 {
        if params.Order == "desc" {
            orderBy = " order by %s desc "
        } else {
            orderBy = " order by %s "
        }
        orderBy = fmt.Sprintf(orderBy, utils.SnakeString(sortorder))
    }
    // limit offset
    var limitOffset string
    if params.Limit > 0 {
         if params.Offset > 0 {
            limitOffset = "limit %d offset %d"
            limitOffset = fmt.Sprintf(limitOffset, params.Limit, params.Offset)
         } else {
            limitOffset = "limit %d"
            limitOffset = fmt.Sprintf(limitOffset, params.Limit)
         }
    }
    // count
    type Total struct {
        Count int64
    }
    var total Total
    
    sqlCount := `select count(*) `
    sql = sqlCount + sqlFrom + sqlWhere
    if err := o.Raw(sql).QueryRow(&total); err == nil && total.Count > 0 {
        // query
        sql = sqlCol + sqlFrom + sqlWhere + orderBy + limitOffset
        o.Raw(sql).QueryRows(&data)
    }

    return data, total.Count
}

// SmsSetDataList 获取角色列表
func SmsSetDataList(params *SmsSetQueryParam) []*SmsSetWithRelated {
    params.Limit = -1
    params.Sort = "Id"
    params.Order = "asc"
    data, _ := SmsSetPageList(params)
    return data
}

// SmsSetBatchDelete 批量删除
func SmsSetBatchDelete(ids []int) (int64, error) {
    query := orm.NewOrm().QueryTable(SmsSetTBName())
    num, err := query.Filter("id__in", ids).Delete()
    return num, err
}

// SmsSetOne 获取单条
func SmsSetOne(id int) (*SmsSet, error) {
    o := orm.NewOrm()
    m := SmsSet{Id: id}
    err := o.Read(&m)
    if err != nil {
        return nil, err
    }
    return &m, nil
}

// SmsSetOne通过传入查询条件得到
func SmsSetOneByParams(params map[string]interface{}) *SmsSet {
    m := SmsSet{}
    query := orm.NewOrm().QueryTable(SmsSetTBName())

    for k, v := range params {
        query = query.Filter(k, v)
    }
    query.One(&m)

    return  &m
}