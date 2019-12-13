package models

import (
    "github.com/astaxie/beego/orm"
)

// TableName 设置表名
func (a *BaseMobileLocation) TableName() string {
    return BaseMobileLocationTBName()
}

// BaseMobileLocation 实体类
type BaseMobileLocation struct {
    No         string `orm:"pk"`
    Location   string
    DistrictNo string
}

// 当前手机号和区号是否是在同一个区域
func MobilePhoneSameAsGateway(no, districtNo string) bool {
    query := orm.NewOrm().QueryTable(BaseMobileLocationTBName())
    return query.Filter("no", no).Filter("district_no", districtNo).Exist()
}