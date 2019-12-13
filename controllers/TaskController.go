package controllers

import (
    "encoding/json"

    "zq/callout_crm/enums"
    "zq/callout_crm/models"

    "fmt"
    "strconv"
    "strings"
    "time"

    "github.com/astaxie/beego/orm"
)

//TaskController 角色管理
type TaskController struct {
    BaseController
}

//Prepare 参考beego官方文档说明
func (c *TaskController) Prepare() {
    //先执行
    c.BaseController.Prepare()
    //如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
    c.checkAuthor("DataGrid", "DataList", "UpdateFieldByIds", "ImportUserClientIndex")
    //如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
    //权限控制里会进行登录验证，因此这里不用再作登录验证
    //c.checkLogin()
}

//Index 角色管理首页
func (c *TaskController) Index() {
    //是否显示更多查询条件的按钮
    c.Data["showMoreQuery"] = false
    //将页面左边菜单的某项激活
    c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
    c.setTpl()
    c.LayoutSections = make(map[string]string)
    c.LayoutSections["headcssjs"] = "task/index_headcssjs.html"
    c.LayoutSections["footerjs"] = "task/index_footerjs.html"
    //页面里按钮权限控制
    c.Data["canEdit"] = c.checkActionAuthor("TaskController", "Edit")
    c.Data["canDelete"] = c.checkActionAuthor("TaskController", "Delete")
}

// DataGrid 表格获取数据
func (c *TaskController) DataGrid() {
    //直接反序化获取json格式的requestbody里的值
    var params models.TaskQueryParam
    json.Unmarshal(c.Ctx.Input.RequestBody, &params)
    if !c.curUser.IsSuper {
        params.UIds = c.curUser.UIds
    }

    //获取数据列表和总数
    data, total := models.TaskPageList(&params)
    //定义返回的数据结构
    result := make(map[string]interface{})
    result["total"] = total
    result["rows"] = data
    c.Data["json"] = result
    c.ServeJSON()
}

//DataList 角色列表
func (c *TaskController) DataList() {
    var params = models.TaskQueryParam{}
    //获取数据列表和总数
    data := models.TaskDataList(&params)
    //定义返回的数据结构
    c.jsonResult(enums.JRCodeSucc, "", data)
}

//Edit 添加、编辑角色界面
func (c *TaskController) Edit() {
    if c.Ctx.Request.Method == "POST" {
        c.Save()
    }
    Id, _ := c.GetInt(":id", 0)
    m := models.Task{Id: Id}

    if Id > 0 {
        o := orm.NewOrm()
        err := o.Read(&m)
        if err != nil {
            c.pageError("数据无效，请刷新后重试")
        }
    }
    c.Data["m"] = m
    c.setTpl("task/edit.html", "shared/layout_pullbox.html")
    c.LayoutSections = make(map[string]string)
    c.LayoutSections["footerjs"] = "task/edit_footerjs.html"
}

//Save 添加、编辑页面 保存
func (c *TaskController) Save() {
    var err error
    m := models.Task{}
    //获取form里的值
    if err = c.ParseForm(&m); err != nil {
        c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
    }
    o := orm.NewOrm()

    if m.Id == 0 {
        m.Created = time.Now().Format(enums.BaseTimeFormat)
        m.BackendUserId = c.curUser.Id
        m.UserCompanyId = c.curUser.UserCompanyId
        if _, err = o.Insert(&m); err == nil {
            c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
        } else {
            c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
        }
    } else {
        if _, err = o.Update(&m, "Name", "Desc"); err == nil {
            c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
        } else {
            c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
        }
    }

}

//Delete 批量删除
func (c *TaskController) Delete() {
    strs := c.GetString("ids")
    ids := make([]int, 0, len(strs))
    for _, str := range strings.Split(strs, ",") {
        if id, err := strconv.Atoi(str); err == nil {
            ids = append(ids, id)
        }
    }
    if num, err := models.TaskBatchDelete(ids); err == nil {
        c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
    } else {
        c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
    }
}

// 根据ids修改一个字段信息
func (c *TaskController) UpdateFieldByIds() {
    strs := c.GetString("ids")
    fieldName := c.GetString("fieldName")
    fieldValue := c.GetString("fieldValue")
    extras := c.GetStrings("extras[]")

    var updates = orm.Params{fieldName: fieldValue}
    var lenExtras = len(extras)
    var validLen int
    if lenExtras > 0 && lenExtras/2 > 0 {
        if lenExtras % 2 > 0 {
            validLen = lenExtras - 1
        } else {
            validLen = lenExtras
        }

        for i:=0; i<validLen; i+=2 {
            updates[extras[i]] = extras[i+1]
        }
    }


    ids := make([]int, 0, len(strs))
    for _, str := range strings.Split(strs, ",") {
        if id, err := strconv.Atoi(str); err == nil {
            ids = append(ids, id)
        }
    }
    query := orm.NewOrm().QueryTable(models.TaskTBName())
    if num, err := query.Filter("id__in", ids).Update(updates); err == nil {
        c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功更新 %d 项", num), 0)
    } else {
        c.jsonResult(enums.JRCodeFailed, "更新字段失败", 0)
    }
}

// 导入任务的详细信息
func (c *TaskController) ImportUserClientIndex() {
    if c.Ctx.Request.Method == "POST" {
        c.Save()
    }
    Id, _ := c.GetInt(":id", 0)
    if Id == 0 {
        c.pageError("无法获取该任务id，请重试！")
    }

    m := models.Task{Id: Id}

    if Id > 0 {
        o := orm.NewOrm()
        err := o.Read(&m)
        if err != nil {
            c.pageError("数据无效，请刷新后重试")
        }
    }
    c.Data["m"] = m
    c.setTpl("task/import.html", "shared/layout_pullbox.html")
    c.LayoutSections = make(map[string]string)
    c.LayoutSections["footerjs"] = "task/import_footerjs.html"
}