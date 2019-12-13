package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"zq/callout_crm/enums"
	"zq/callout_crm/utils"

	"github.com/astaxie/beego/orm"
)

// TableName 设置表名
func (a *UserClientField) TableName() string {
	return UserClientFieldTBName()
}

// UserClientFieldQueryParam 用于搜索的类
type UserClientFieldQueryParam struct {
	BaseQueryParam
	UserCompanyId int
}

// UserClientField 用户角色 实体类
type UserClientField struct {
	Id             int
	ColumnName     string // 列的英文名
	FieldName      string // 列的中文名
	FieldSpecies   int    // 字段的所属。0：自定义，1：系统默认
	FieldType      int    // 字段类型。0：单行文本，1：多行文本，2：单选，3：多选，4：日期，5：数字
	FieldTypeValue string // 单选，多选的值，pg中用数组来做
	ListShow       bool   // 列表显示
	AddShow        bool   // 添加显示
	QueryShow      bool   // 查询显示
	Required       int    // 是否必填。-1为添加数据库默认要填的，0不必填，1必填
	FieldColor     string // 列表颜色
	UserCompanyId  int    // 对应的公司id
}

// UserClientFieldPageList 获取分页数据
func UserClientFieldPageList(params *UserClientFieldQueryParam) ([]*UserClientField, int64) {
	query := orm.NewOrm().QueryTable(UserClientFieldTBName())
	data := make([]*UserClientField, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	case "FieldName":
		sortorder = "FieldName"
	case "FieldSpecies":
		sortorder = "FieldSpecies"
	case "FieldType":
		sortorder = "FieldType"
	case "ListShow":
		sortorder = "ListShow"
	case "AddShow":
		sortorder = "AddShow"
	case "QueryShow":
		sortorder = "QueryShow"
	case "Required":
		sortorder = "Required"
	case "FieldColor":
		sortorder = "FieldColor"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}

	// 设置条件
	if params.UserCompanyId > 0 {
		query = query.Filter("user_company_id__in", params.UserCompanyId, 0)
	}

	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

// UserClientFieldDataList 获取角色列表
func UserClientFieldDataList(params *UserClientFieldQueryParam) []*UserClientField {
	params.Limit = -1
	params.Sort = "Id"
	params.Order = "asc"
	data, _ := UserClientFieldPageList(params)
	return data
}

// UserClientFieldBatchDelete 批量删除
func UserClientFieldBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(UserClientFieldTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// UserClientFieldOne 获取单条
func UserClientFieldOne(id int) (*UserClientField, error) {
	o := orm.NewOrm()
	m := UserClientField{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func UserClientFieldMaxColumnName(comId int) string {
	type MaxColumnName struct {
		Result string
	}
	o := orm.NewOrm()
	sql := `SELECT max(column_name) result from %s where column_name like 'column%%' and user_company_id = ?`
	sql = fmt.Sprintf(sql, UserClientFieldTBName())

	var mcn MaxColumnName
	if err := o.Raw(sql, comId).QueryRow(&mcn); err != nil {
		utils.LogError("UserClientFieldMaxColumnName" + err.Error())
	}
	return mcn.Result
}

func UserClientFieldValidColumnName(comId int) (string, error) {
	type Valid struct {
		ColumnName string
	}
	var err error
	var num int64
	o := orm.NewOrm()
	sql := `SELECT column_name from %s where column_name like 'column%%' and user_company_id = ? order by column_name`
	sql = fmt.Sprintf(sql, UserClientFieldTBName())

	var vs []Valid
	if num, err = o.Raw(sql, comId).QueryRows(&vs); err != nil {
		utils.LogError("UserClientFieldValidColumnName" + err.Error())
		return "", err
	} else if num >= enums.MaxSelfDefineColumn {
		return "", errors.New("Longer than max(16) segment.")
	} else if num == 0 {
		return "column1", nil
	} else {
		var tmpFormat = `column%d`
		var final int
		for i, v := range vs {
			tmp := fmt.Sprintf(tmpFormat, i+1)
			if strings.Compare(v.ColumnName, tmp) == 0 {
				final = i + 1
				continue
			}
			return tmp, nil
		}
		return "column" + strconv.Itoa(final+1), nil
	}
	return "", err
}
