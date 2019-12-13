package controllers

import (
    "fmt"
    "strconv"
    "strings"
    "encoding/json"

    "zq/callout_crm/enums"
    "zq/callout_crm/models"  

    "github.com/astaxie/beego/orm"
)

//TaskDetailController
type TaskDetailController struct {
    BaseController
}

//Prepare 参考beego官方文档说明
func (c *TaskDetailController) Prepare() {
    //先执行
    c.BaseController.Prepare()
    //如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
    c.checkAuthor("Assign", "DataGrid", "AverageAssign")
    //如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
    //权限控制里会进行登录验证，因此这里不用再作登录验证
    //c.checkLogin()
}

func (c *TaskDetailController) Assign() {
    strs := c.GetString("ids")
    taskId, _ := c.GetInt("task_id", 0)
    userId, _ := c.GetInt("user_id", 0)

    ids := make([]int, 0, len(strs))
    for _, str := range strings.Split(strs, ",") {
        if id, err := strconv.Atoi(str); err == nil {
            ids = append(ids, id)
        }
    }

    if len(ids) == 0 || taskId == 0 || userId == 0 {
        c.jsonResult(enums.JRCodeFailed, "获取数据失败，请刷新重试", 0)
    }

    // prepare sql
    sqls := make([]string, 0)
    sqlTmpl := `insert into ` + models.TaskDetailTBName() + 
               ` (task_id, user_client_id, belong_user_id) values (%d, %d, %d)`
    for _, id := range ids {
        sqls = append(sqls, fmt.Sprintf(sqlTmpl, taskId, id, userId))
    }

    // execute sql
    o := orm.NewOrm()
    _, err := o.Raw(strings.Join(sqls, ";")).Exec()
    if err == nil {
        // num, _ := res.RowsAffected()
        c.jsonResult(enums.JRCodeSucc, "分配成功", 0)
    } else {
        c.jsonResult(enums.JRCodeFailed, "分配失败，请刷新重试", 0)
    }
}

func (c *TaskDetailController) DataGrid() {
    //直接反序化获取json格式的requestbody里的值
    var params models.TaskDetailQueryParam
    json.Unmarshal(c.Ctx.Input.RequestBody, &params)
    // 根据uid获取记录
    if !c.curUser.IsSuper {
        params.CompanyId = c.curUser.UserCompanyId
        if len(params.BackendUserIds) > 0 {
            if params.BackendUserIds[0] == -1 {
                params.BackendUserIds = c.curUser.UIds
            }
        } else {
            params.BackendUserIds = []int{c.curUser.Id}
        }
    }

    //获取数据列表和总数
    data, total := models.TaskDetailPageList(&params)
    //定义返回的数据结构
    result := make(map[string]interface{})
    result["total"] = total
    result["rows"] = data
    c.Data["json"] = result
    c.ServeJSON()
}

func (c *TaskDetailController) AverageAssign() {
    taskId, _ := c.GetInt("task_id", 0)
    strs := c.GetStrings("uids[]")

    uids := make([]int, 0, len(strs))
    for _, str := range strs {
        if id, err := strconv.Atoi(str); err == nil {
            uids = append(uids, id)
        }
    }

    if len(uids) == 0 || taskId == 0 {
        c.jsonResult(enums.JRCodeFailed, "获取数据失败，请刷新重试", 0)
    }
    // 根据uid获取记录
    var params models.TaskDetailQueryParam
    if !c.curUser.IsSuper {
        params.CompanyId = c.curUser.UserCompanyId
        if len(params.BackendUserIds) > 0 {
            if params.BackendUserIds[0] == -1 {
                params.BackendUserIds = c.curUser.UIds
            }
        } else {
            params.BackendUserIds = []int{c.curUser.Id}
        }
    }
    ucs := models.UnAssignUserClientList(&params)
    // 确定平均分配的值
    pivots := make(map[int][]int)
    ucsLen := len(ucs)
    uidsLen := len(uids)
    if ucsLen == 0 {
        c.jsonResult(enums.JRCodeFailed, "没有需要分配的项", 0)
    } else if ucsLen < uidsLen {
        pivots[uids[0]] = []int{0, ucsLen}
    } else {
        cnt := ucsLen/uidsLen
        for i:=0; i<uidsLen; i++ {
            if cnt*(i+1) < ucsLen && cnt*(i+2) > ucsLen {
                pivots[uids[i]] = []int{cnt*i, ucsLen}
            } else {
                pivots[uids[i]] = []int{cnt*i, cnt*(i+1)}
            }
        }
    }

    // prepare sql
    sqls := make([]string, 0)
    sqlTmpl := `insert into ` + models.TaskDetailTBName() + 
               ` (task_id, user_client_id, belong_user_id) values (%d, %d, %d)`
    for k, v := range pivots {
        for index, item := range ucs {
            if index >= v[0] && index < v[1] {
                sqls = append(sqls, fmt.Sprintf(sqlTmpl, taskId, item.Id, k))
            }
        }
    }

    // execute sql
    o := orm.NewOrm()
    _, err := o.Raw(strings.Join(sqls, ";")).Exec()
    if err == nil {
        // num, _ := res.RowsAffected()
        c.jsonResult(enums.JRCodeSucc, "分配成功", 0)
    } else {
        c.jsonResult(enums.JRCodeFailed, "分配失败，请刷新重试", 0)
    }
}