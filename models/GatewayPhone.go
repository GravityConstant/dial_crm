package models

import (
    "fmt"
    "strings"
    "strconv"

    "zq/callout_crm/utils"

    "github.com/astaxie/beego/orm"
)

// TableName 设置表名
func (a *GatewayPhone) TableName() string {
    return GatewayPhoneTBName()
}

// GatewayPhoneQueryParam 用于搜索的类
type GatewayPhoneQueryParam struct {
    BaseQueryParam
    PhoneLike       string
    GatewayNameLike string
    ExtNoLike       string
    UserId          int
    Ids             []int
}

// GatewayPhone 用户角色 实体类
type GatewayPhone struct {
    Id          int
    Phone       string
    GatewayId   string
    AgentId     int
    Created     string
}

type GatewayPhoneWithRelated struct {
    Id          int
    Phone       string
    GatewayId   int
    GatewayName string 
    GatewayUrl  string
    AgentId     int
    ExtNo       string
    Created     string
    GatewayType int    // 0电话线，1插手机卡
}

// GatewayPhonePageList 获取分页数据
func GatewayPhonePageList(params *GatewayPhoneQueryParam) ([]*GatewayPhoneWithRelated, int64) {
    var err error
    o := orm.NewOrm()
    data := make([]*GatewayPhoneWithRelated, 0)
    //默认排序
    sortorder := "Id"
    switch params.Sort {
    case "Id":
        sortorder = "Id"
    case "Phone":
        sortorder = "Phone"
    case "GatewayName":
        sortorder = "GatewayName"
    case "ExtNo":
        sortorder = "ExtNo"
    }
    
    var sql string

    sqlCol := 
    `SELECT gp.id, gp.phone, gp.gateway_id, g.gateway_name, g.gateway_url, g.gateway_type,
    gp.agent_id, a.ext_no, gp.created `
    
    sqlFrom := 
    `FROM %s gp 
    LEFT JOIN %s a on gp.agent_id=a.id 
    LEFT JOIN %s g on gp.gateway_id=g.id `
    sqlFrom = fmt.Sprintf(sqlFrom, GatewayPhoneTBName(), AgentTBName(), GatewayTBName())
    // where
    var sqlWhere string
    andCond := make([]string, 0)
    if len(params.PhoneLike) > 0 {
        andCond = append(andCond, fmt.Sprintf(`gp.phone like '%%%s%%'`, params.PhoneLike))
    }
    if len(params.GatewayNameLike) > 0 {
        andCond = append(andCond, fmt.Sprintf(`g.gateway_name ilike '%%%s%%'`, params.GatewayNameLike))
    }
    if len(params.ExtNoLike) > 0 {
        andCond = append(andCond, fmt.Sprintf(`a.ext_no like '%%%s%%'`, params.ExtNoLike))
    }
    if params.UserId > 0 {
        andCond = append(andCond, fmt.Sprintf(`a.backend_user_id = %d`, params.UserId))
    }
    if len(params.Ids) > 0 {
        idstrs := make([]string, 0)
        for _, id := range params.Ids {
            idstrs = append(idstrs, strconv.Itoa(id))
        }
        andCond = append(andCond, fmt.Sprintf(`gp.id in (%s)`, strings.Join(idstrs, ",")))
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

// GatewayPhoneDataList 获取角色列表
func GatewayPhoneDataList(params *GatewayPhoneQueryParam) []*GatewayPhoneWithRelated {
    params.Limit = -1
    params.Sort = "Id"
    params.Order = "asc"
    data, _ := GatewayPhonePageList(params)
    return data
}

// GatewayPhoneBatchDelete 批量删除
func GatewayPhoneBatchDelete(ids []int) (int64, error) {
    query := orm.NewOrm().QueryTable(GatewayPhoneTBName())
    num, err := query.Filter("id__in", ids).Delete()
    return num, err
}

// GatewayPhoneOne 获取单条
func GatewayPhoneOne(id int) (*GatewayPhone, error) {
    o := orm.NewOrm()
    m := GatewayPhone{Id: id}
    err := o.Read(&m)
    if err != nil {
        return nil, err
    }
    return &m, nil
}

// GatewayPhoneBatchDelete 批量分配坐席
func GatewayPhoneBatchAllocateAgent(ids []int, agentId int) (int64, error) {
    query := orm.NewOrm().QueryTable(GatewayPhoneTBName())
    num, err := query.Filter("id__in", ids).Update(orm.Params{
        "agent_id": agentId,
    })
    return num, err
}

// 通过中继获取ext_no
func GetExtNoByTrunkNo(trunkNo string) string {
    sql := `select a.ext_no from %s gp inner join %s a on gp.agent_id=a.id where phone='%s'`
    sql = fmt.Sprintf(sql, GatewayPhoneTBName(), AgentTBName(), trunkNo)

    type Rst struct {
        ExtNo string
    }
    var rst Rst
    o := orm.NewOrm()
    o.Raw(sql).QueryRow(&rst)

    return rst.ExtNo
}