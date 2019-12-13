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

//SmsSetController 角色管理
type SmsSetController struct {
    BaseController
}

//Prepare 参考beego官方文档说明
func (c *SmsSetController) Prepare() {
    //先执行
    c.BaseController.Prepare()
    //如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
    c.checkAuthor("DataGrid", "DataList", "Edit")
    //如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
    //权限控制里会进行登录验证，因此这里不用再作登录验证
    //c.checkLogin()
}

//Index 角色管理首页
func (c *SmsSetController) Index() {
    //是否显示更多查询条件的按钮
    c.Data["showMoreQuery"] = false
    //将页面左边菜单的某项激活
    c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
    // 非超级管理员直接显示edit
    if !c.curUser.IsSuper {
        o := orm.NewOrm()
        m := models.SmsSet{UserCompanyId: c.curUser.UserCompanyId}
        o.Read(&m, "UserCompanyId")
        c.Data["m"] = m
        c.setTpl("smsset/edit.html", "shared/layout_page.html")
        c.LayoutSections = make(map[string]string)
        c.LayoutSections["footerjs"] = "smsset/edit_footerjs.html"
    } else {
        c.setTpl()
        c.LayoutSections = make(map[string]string)
        c.LayoutSections["headcssjs"] = "smsset/index_headcssjs.html"
        c.LayoutSections["footerjs"] = "smsset/index_footerjs.html"
        //页面里按钮权限控制
        c.Data["canEdit"] = c.checkActionAuthor("SmsSetController", "Edit")
        c.Data["canDelete"] = c.checkActionAuthor("SmsSetController", "Delete")
    }
    
}

// DataGrid 角色管理首页 表格获取数据
func (c *SmsSetController) DataGrid() {
    //直接反序化获取json格式的requestbody里的值
    var params models.SmsSetQueryParam
    json.Unmarshal(c.Ctx.Input.RequestBody, &params)
    //获取数据列表和总数
    data, total := models.SmsSetPageList(&params)
    //定义返回的数据结构
    result := make(map[string]interface{})
    result["total"] = total
    result["rows"] = data
    c.Data["json"] = result
    c.ServeJSON()
}

// DataList 角色列表
func (c *SmsSetController) DataList() {
    var params = models.SmsSetQueryParam{}
    //获取数据列表和总数
    data := models.SmsSetDataList(&params)
    //定义返回的数据结构
    c.jsonResult(enums.JRCodeSucc, "", data)
}

// Edit 添加、编辑角色界面
func (c *SmsSetController) Edit() {
    if c.Ctx.Request.Method == "POST" {
        c.Save()
    }
    Id, _ := c.GetInt(":id", 0)
    m := models.SmsSet{Id: Id}
    if Id > 0 {
        o := orm.NewOrm()
        err := o.Read(&m)
        if err != nil {
            c.pageError("数据无效，请刷新后重试")
        }
    }
    c.Data["m"] = m
    c.setTpl("smsset/edit.html", "shared/layout_pullbox.html")
    c.LayoutSections = make(map[string]string)
    c.LayoutSections["footerjs"] = "smsset/edit_footerjs.html"
}

//Save 添加、编辑页面 保存
func (c *SmsSetController) Save() {
    var err error
    m := models.SmsSet{}
    //获取form里的值
    if err = c.ParseForm(&m); err != nil {
        c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
    }
    o := orm.NewOrm()
    if m.Id == 0 {
        // insert添加created时间
        m.Created = time.Now().Format(enums.BaseTimeFormat)
        m.Updated = time.Now().Format(enums.BaseTimeFormat)
        m.BackendUserId = c.curUser.Id
        m.UserCompanyId = c.curUser.UserCompanyId
        if _, err = o.Insert(&m); err == nil {
            c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
        } else {
            c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
        }

    } else {
        m.Updated = time.Now().Format(enums.BaseTimeFormat)
        // password为空不更新，不为空更新
        fields := []string{"Account", "Signature", "Updated"}
        if len(m.Password) > 0 {
            fields = append(fields, "Password")
        }
        if _, err = o.Update(&m, fields...); err == nil {
            c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
        } else {
            c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
        }
    }

}

//Delete 批量删除
func (c *SmsSetController) Delete() {
    strs := c.GetString("ids")
    ids := make([]int, 0, len(strs))
    for _, str := range strings.Split(strs, ",") {
        if id, err := strconv.Atoi(str); err == nil {
            ids = append(ids, id)
        }
    }
    if num, err := models.SmsSetBatchDelete(ids); err == nil {
        c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
    } else {
        c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
    }
}
