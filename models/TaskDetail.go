package models

import (
    "fmt"
    "strconv"
    "strings"

    "zq/callout_crm/utils"

    "github.com/astaxie/beego/orm"
)

// TableName 设置表名
func (a *TaskDetail) TableName() string {
    return TaskDetailTBName()
}

// TaskDetailQueryParam 用于搜索的类
type TaskDetailQueryParam struct {
    BaseQueryParam
    BackendUserIds []int
    CompanyId      int
}

// TaskDetail 用户角色 实体类
type TaskDetail struct {
    Id           int
    TaskId       string
    UserClientId int
    CallState    string
    BelongUserId int
}

type TaskDetailUserClient struct {
    Id          int
    Name        string
    MobilePhone string
    Checked     int
}

// TaskDetailPageList 获取分页数据
func TaskDetailPageList(params *TaskDetailQueryParam) ([]*TaskDetailUserClient, int64) {
    o := orm.NewOrm()
    data := make([]*TaskDetailUserClient, 0)
    paramsWhere := make([]string, 0)
    var sql string
    var where string
    //默认排序
    sortorder := "Updated"
    switch params.Sort {
    case "Id":
        sortorder = "Id"
    case "MobilePhone":
        sortorder = "MobilePhone"
    case "Name":
        sortorder = "Name"
    case "Updated":
        sortorder = "Updated"
    }

    // sql
    sqlCol := `SELECT id, name, mobile_phone, (` +
                `SELECT id from %s d WHERE c.id=d.user_client_id ` +
              `) as checked `
    sqlCol = fmt.Sprintf(sqlCol, TaskDetailTBName())
    sqlFrom := `FROM %s c `
    sqlFrom = fmt.Sprintf(sqlFrom, UserClientTBName())
    var sqlOrderBy string
    var sqlLimit string

    // where
    if params.CompanyId > 0 {
        paramsWhere = append(paramsWhere, fmt.Sprintf(`user_company_id=%d`, params.CompanyId))
    }
    // 归属人字段给了belong_backend_user_id
    if len(params.BackendUserIds) > 0 {
        ins := []string{"0"}    // 公共池
        for _, val := range params.BackendUserIds {
            ins = append(ins, strconv.Itoa(val))
        }
        paramsWhere = append(paramsWhere, fmt.Sprintf("belong_backend_user_id in (%s)", strings.Join(ins, ",")))
    }
    where = strings.Join(paramsWhere, " and ")
    if len(where) > 0 {
        sqlFrom += "where " + where
    }
    // order by
    if len(sortorder) > 0 {
        if params.Order == "desc" {
            sqlOrderBy = fmt.Sprintf(` order by %s desc `, utils.SnakeString(sortorder))
        }
    }
    // limit
    if params.Limit > 0 {
        limit := `limit %d offset %d`
        sqlLimit = fmt.Sprintf(limit, params.Limit, params.Offset)
    }
    // query
    sql = sqlCol + sqlFrom + sqlOrderBy + sqlLimit
    o.Raw(sql).QueryRows(&data)
    // count
    type Total struct {
        Count int64
    }
    sql = `select count(*) as count ` + sqlFrom
    var total Total
    o.Raw(sql).QueryRow(&total)

    return data, total.Count
}

// 获取未分配的user_client
func UnAssignUserClientList(params *TaskDetailQueryParam) []*TaskDetailUserClient {
    o := orm.NewOrm()
    data := make([]*TaskDetailUserClient, 0)

    sql := `SELECT id, name, mobile_phone, checked FROM ( `+
                `SELECT id, name, mobile_phone, (` +
                    `SELECT id from %s d WHERE d.user_client_id=c.id) as checked `+
                    `FROM %s c %s `+
            `) tmp WHERE checked is null order by id`
    // 准备填充的数据
    var where string
    paramsWhere := make([]string, 0)
    if params.CompanyId > 0 {
        paramsWhere = append(paramsWhere, fmt.Sprintf(`user_company_id=%d`, params.CompanyId))
    }
    // 归属人字段给了belong_backend_user_id
    if len(params.BackendUserIds) > 0 {
        ins := []string{"0"}    // 公共池
        for _, val := range params.BackendUserIds {
            ins = append(ins, strconv.Itoa(val))
        }
        paramsWhere = append(paramsWhere, fmt.Sprintf("belong_backend_user_id in (%s)", strings.Join(ins, ",")))
    }
    where = strings.Join(paramsWhere, " and ")
    if len(where) > 0 {
        where = "where " + where
    }
    // fill
    sql = fmt.Sprintf(sql, TaskDetailTBName(), UserClientTBName(), where)

    o.Raw(sql).QueryRows(&data)

    return data
}