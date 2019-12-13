package controllers

import (
    "encoding/json"
    "fmt"

    "zq/callout_crm/enums"
    "zq/callout_crm/models"

    "github.com/astaxie/beego/orm"
)

//MyTaskController 角色管理
type MyTaskController struct {
    BaseController
}

//Prepare 参考beego官方文档说明
func (c *MyTaskController) Prepare() {
    //先执行
    c.BaseController.Prepare()
    //如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
    c.checkAuthor("DataGrid", "UpdateCallState")
    //如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
    //权限控制里会进行登录验证，因此这里不用再作登录验证
    //c.checkLogin()
}

//Index 首页
func (c *MyTaskController) Index() {
    //是否显示更多查询条件的按钮
    c.Data["showMoreQuery"] = false
    //将页面左边菜单的某项激活
    c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
    c.setTpl()
    c.LayoutSections = make(map[string]string)
    c.LayoutSections["headcssjs"] = "mytask/index_headcssjs.html"
    c.LayoutSections["footerjs"] = "mytask/index_footerjs.html"
}

// DataGrid 表格获取数据
func (c *MyTaskController) DataGrid() {
    //直接反序化获取json格式的requestbody里的值
    var params models.MyTaskQueryParam
    json.Unmarshal(c.Ctx.Input.RequestBody, &params)
    if !c.curUser.IsSuper {
        params.UIds = []int{c.curUser.Id}
    }

    //获取数据列表和总数
    data, total := models.MyTaskPageList(&params)
    //定义返回的数据结构
    result := make(map[string]interface{})
    result["total"] = total
    result["rows"] = data
    c.Data["json"] = result
    c.ServeJSON()
}

// 在呼叫完成之后，偷偷的更新这个状态，其实可以不用返回值的。
func (c *MyTaskController) UpdateCallState() {
    id, _ := c.GetInt("id", 0)
    callstate := c.GetString("callstate")

    if id == 0 || len(callstate) == 0 {
        c.jsonResult(enums.JRCodeFailed, "获取数据失败，请刷新重试", -1)
    }

    // execute sql
    o := orm.NewOrm()

    sql := `update %s set call_state='%s' where id=%d`
    sql = fmt.Sprintf(sql, models.TaskDetailTBName(), callstate, id)
    if _, err := o.Raw(sql).Exec(); err == nil {
        c.jsonResult(enums.JRCodeSucc, "更新呼叫状态成功", 0)
    } else {
        c.jsonResult(enums.JRCodeFailed, "更新呼叫状态失败", -1)
    }
}