package models

import (
	"fmt"
	"path/filepath"
	"os"
	"strings"

	"github.com/astaxie/beego/orm"
)

// TableName 设置表名
func (a *Agent) TableName() string {
	return AgentTBName()
}

// AgentQueryParam 用于搜索的类
type AgentQueryParam struct {
	BaseQueryParam
	ExtNoLike string
	UIds      []int
}

// Agent 用户角色 实体类
type Agent struct {
	Id                        int
	ExtNo                     string
	ExtPwd                    string
	GatewayPhoneNumber        string
	OriginationCallerIdNumber string
	BackendUserId             int
	CallWay                   int
	GatewayId                 int
	Param                     string
	DefaultTrunk              int
	BindPhone                 string
}

// AgentPageList 获取分页数据
func AgentPageList(params *AgentQueryParam) ([]*Agent, int64) {
	query := orm.NewOrm().QueryTable(AgentTBName())
	data := make([]*Agent, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	case "ExtNo":
		sortorder = "ExtNo"
	case "GatewayPhoneNumber":
		sortorder = "GatewayPhoneNumber"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	if len(params.ExtNoLike) > 0 {
		query = query.Filter("ExtNo__istartswith", params.ExtNoLike)
	}
	if len(params.UIds) > 0 {
		query = query.Filter("backend_user_id__in", params.UIds)
	}

	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

// AgentDataList 获取角色列表
func AgentDataList(params *AgentQueryParam) []*Agent {
	params.Limit = -1
	params.Sort = "Id"
	params.Order = "asc"
	data, _ := AgentPageList(params)
	return data
}

// AgentBatchDelete 批量删除
func AgentBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(AgentTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// AgentOne 获取单条
func AgentOne(id int) (*Agent, error) {
	o := orm.NewOrm()
	m := Agent{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// 根据uid获取单条
func AgentOneByUid(uid int) (*Agent, error) {
	query := orm.NewOrm().QueryTable(AgentTBName())
	var m Agent
	err := query.Filter("backend_user_id", uid).One(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

type AgentGateway struct {
	ExtNo                     string
	GatewayPhoneNumber        string
	OriginationCallerIdNumber string
	GatewayName               string
	GatewayUrl                string
	CallWay                   int
	Param                     string
	DefaultTrunk              int
}

// 根据当前用户的id获取座席
func AgentGatewayOneByUserId(userId int) (*AgentGateway, error) {
	o := orm.NewOrm()
	m := AgentGateway{}
	sql := `select a.ext_no, a.gateway_phone_number, a.origination_caller_id_number, g.gateway_name, g.gateway_url, a.default_trunk,  a.call_way, a.param from %s a left join %s g on a.gateway_id=g.id where a.backend_user_id=?`
	sql = fmt.Sprintf(sql, AgentTBName(), GatewayTBName())
	err := o.Raw(sql, userId).QueryRow(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// 座席号码取自freeswitch的配置文件conf/directory/default下的文件名
type FsRegisterUser struct {
	Name string
}

// 座席号码取自freeswitch的配置文件conf/directory/default下的文件名
func FsRegisterUserList() []*FsRegisterUser {
	data := make([]*FsRegisterUser, 0)
	filepath.Walk("/usr/local/freeswitch/conf/directory/default", func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			filename := strings.Split(info.Name(), ".")
			fsFile := FsRegisterUser{
				Name: strings.TrimSpace(filename[0]),
			}
			data = append(data, &fsFile)
		}
		return nil
	})

    return data
}

// freeswitch中用户没有在数据库中的列表
func FsUserListNotInDbList() []*FsRegisterUser {
    res := make([]*FsRegisterUser, 0)
	data := FsRegisterUserList()

	// 排除已加入的座席号码
    extNos := make([]*FsRegisterUser, 0) 
    o := orm.NewOrm()
    sql := `select ext_no as name from %s order by ext_no`
    sql = fmt.Sprintf(sql, AgentTBName())
    if num, err := o.Raw(sql).QueryRows(&extNos); err == nil && num > 0 {
        found := false
        for _, item := range data {
            for _, ext := range extNos {
                if strings.Compare(ext.Name, item.Name) == 0 {
                    found = true
                    break
                }
            }
            if !found {
                res = append(res, item)
            }
            found = false
        }
    }


    return res
}