package models

import (
    "fmt"
    "strconv"
    "strings"

    "zq/callout_crm/utils"

    "github.com/astaxie/beego/orm"
)

// MyTaskQueryParam 用于搜索的类
type MyTaskQueryParam struct {
    BaseQueryParam
    NameLike  string
    PhoneLike string
    CallState string
    UIds      []int
}

// MyTask 实体类
type MyTask struct {
    Id              int
    UserClientId    int
    Name            string
    CallState       string
    MobilePhone     string
}

// MyTaskPageList 获取分页数据
func MyTaskPageList(params *MyTaskQueryParam) ([]*MyTask, int64) {
    var err error
    o := orm.NewOrm()
    data := make([]*MyTask, 0)
    //默认排序
    sortorder := "Id"
    switch params.Sort {
    case "Id":
        sortorder = "Id"
    case "Name":
        sortorder = "Name"
    case "CallState":
        sortorder = "CallState"
    case "MobilePhone":
        sortorder = "MobilePhone"
    }
    
    var sql string

    sqlCol := `SELECT d.id, d.user_client_id, t.name, d.call_state, c.mobile_phone `

    sqlFrom := `FROM %s d INNER JOIN %s t on d.task_id=t.id INNER JOIN %s c on c.id=d.user_client_id `
    sqlFrom = fmt.Sprintf(sqlFrom, TaskDetailTBName(), TaskTBName(), UserClientTBName())
    // where
    var sqlWhere string
    andCond := make([]string, 0)
    if len(params.NameLike) > 0 {
        andCond = append(andCond, fmt.Sprintf(`t.name ilike '%%%s%%'`, params.NameLike))
    }
    if len(params.PhoneLike) > 0 {
        andCond = append(andCond, fmt.Sprintf(`c.mobile_phone ilike '%%%s%%'`, params.PhoneLike))
    }
    if len(params.CallState) > 0 {
        if params.CallState == "未拨打" {
            andCond = append(andCond, `d.call_state = ''`)
        } else {
            andCond = append(andCond, fmt.Sprintf(`d.call_state = '%s'`, params.CallState))
        }
    }
    if len(params.UIds) > 0 {
        uidstrs := make([]string, 0)
        for _, uid := range params.UIds {
            uidstrs = append(uidstrs, strconv.Itoa(uid))
        }
        andCond = append(andCond, fmt.Sprintf(`d.belong_user_id in (%s)`, strings.Join(uidstrs, ",")))
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
    if err = o.Raw(sql).QueryRow(&total); err == nil && total.Count > 0 {
        // query
        sql = sqlCol + sqlFrom + sqlWhere + orderBy + limitOffset
        o.Raw(sql).QueryRows(&data)
    }

    return data, total.Count
}