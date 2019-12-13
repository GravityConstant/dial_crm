package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func (a *Asq) TableName() string {
	return AsqTBName()
}

//Asq 实体类
type Asq struct {
	Id            int
	Name          string
	Description   string
	BackendUserId int
	Created       time.Time `orm:"auto_now_add;type(datetime)"`
	Updated       time.Time `orm:"auto_now_add;type(datetime)"`
}

type AsqQueryParam struct {
	BaseQueryParam
}

// AsqPageList 获取分页数据
func AsqPageList(params *AsqQueryParam) ([]*Asq, int64) {
	query := orm.NewOrm().QueryTable(AsqTBName())
	data := make([]*Asq, 0)

	total, _ := query.Count()
	query.Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

func AsqOne(id int) (*Asq, error) {
	o := orm.NewOrm()
	m := Asq{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}