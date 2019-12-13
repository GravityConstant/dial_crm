package models

import (
    "github.com/astaxie/beego/orm"
)

func (a *GatewayGroup) TableName() string {
    return GatewayGroupTBName()
}

type GatewayGroupQueryParam struct {
    BaseQueryParam
    NameLike string
}

type GatewayGroup struct {
    Id               int
    GatewayGroupName string
    GatewayGroupMap  []*GatewayGroupMap `orm:"reverse(many)" json:"-"` // 设置一对多的反向关系
}

func GatewayGroupPageList(params *GatewayGroupQueryParam) ([]*GatewayGroup, int64) {
    query := orm.NewOrm().QueryTable(GatewayGroupTBName())
    data := make([]*GatewayGroup, 0)
    //默认排序
    sortorder := "Id"
    switch params.Sort {
    case "Id":
        sortorder = "Id"
    }
    if params.Order == "desc" {
        sortorder = "-" + sortorder
    }
    query = query.Filter("gateway_group_name__icontains", params.NameLike)
    total, _ := query.Count()
    query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
    return data, total
}

// 根据id获取单条
func GatewayGroupOne(id int) (*GatewayGroup, error) {
    o := orm.NewOrm()
    m := GatewayGroup{Id: id}
    err := o.Read(&m)
    if err != nil {
        return nil, err
    }
    return &m, nil
}

// RoleDataList 获取角色列表
func GatewayGroupDataList(params *GatewayGroupQueryParam) []*GatewayGroup {
    params.Limit = -1
    params.Sort = "Id"
    params.Order = "asc"
    data, _ := GatewayGroupPageList(params)
    return data
}

// BatchDelete 批量删除
func GatewayGroupBatchDelete(ids []int) (int64, error) {
    query := orm.NewOrm().QueryTable(GatewayGroupTBName())
    num, err := query.Filter("id__in", ids).Delete()
    return num, err
}