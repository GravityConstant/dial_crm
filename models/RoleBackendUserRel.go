package models

import (
    "time"
    "fmt"
    "strconv"

    "github.com/astaxie/beego/orm"
)

// RoleBackendUserRel 角色与用户关系
type RoleBackendUserRel struct {
	Id          int
	Role        *Role        `orm:"rel(fk)"`  //外键
	BackendUser *BackendUser `orm:"rel(fk)" ` // 外键
	Created     time.Time    `orm:"auto_now_add;type(datetime)"`
}

// TableName 设置表名
func (a *RoleBackendUserRel) TableName() string {
	return RoleBackendUserRelTBName()
}


// 根据uid获取roleids
func GetRoleIdsByUids(uid int) []int {
    roleIds := make([]int, 0)
    o := orm.NewOrm()

    sql := `select role_id from %s where backend_user_id=?`
    sql = fmt.Sprintf(sql, RoleBackendUserRelTBName())
    
    var list orm.ParamsList
    num, err := o.Raw(sql, uid).ValuesFlat(&list)
    if err == nil && num > 0 {
        for _, val := range list {
            if v, ok := val.(string); ok {
                if vi, err := strconv.Atoi(v); err == nil {
                    roleIds = append(roleIds, vi)
                }
            }
        }
    }

    return roleIds
}