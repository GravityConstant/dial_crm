package models

import (
    // "fmt"

    "github.com/astaxie/beego/orm"
)

// TableName 设置表名
func (a *Task) TableName() string {
    return TaskTBName()
}

// TaskQueryParam 用于搜索的类
type TaskQueryParam struct {
    BaseQueryParam
    NameLike string
    UIds     []int
}

// Task 用户角色 实体类
type Task struct {
    Id              int
    Name            string
    State           int
    Created         string
    BackendUserId   int
    UserCompanyId int
    Desc            string
}

// TaskPageList 获取分页数据
func TaskPageList(params *TaskQueryParam) ([]*Task, int64) {
    query := orm.NewOrm().QueryTable(TaskTBName())
    data := make([]*Task, 0)
    cond := orm.NewCondition()
    //默认排序
    sortorder := "Id"
    switch params.Sort {
    case "Id":
        sortorder = "Id"
    case "Name":
        sortorder = "Name"
    case "Created":
        sortorder = "Created"
    }
    if params.Order == "desc" {
        sortorder = "-" + sortorder
    }
    if len(params.NameLike) > 0 {
        cond = cond.And("name__icontains", params.NameLike)
    }
    if len(params.UIds) > 0 {
        cond = cond.And("backend_user_id__in", params.UIds)
    }
    query = query.SetCond(cond)
    total, _ := query.Count()
    query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
    return data, total
}

// TaskDataList 获取角色列表
func TaskDataList(params *TaskQueryParam) []*Task {
    params.Limit = -1
    params.Sort = "Id"
    params.Order = "asc"
    data, _ := TaskPageList(params)
    return data
}

// TaskBatchDelete 批量删除
func TaskBatchDelete(ids []int) (int64, error) {
    query := orm.NewOrm().QueryTable(TaskTBName())
    num, err := query.Filter("id__in", ids).Delete()
    return num, err
}

// TaskOne 获取单条
func TaskOne(id int) (*Task, error) {
    o := orm.NewOrm()
    m := Task{Id: id}
    err := o.Read(&m)
    if err != nil {
        return nil, err
    }
    return &m, nil
}
