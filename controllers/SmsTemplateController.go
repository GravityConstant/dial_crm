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

//SmsTemplateController 角色管理
type SmsTemplateController struct {
    BaseController
}

//Prepare 参考beego官方文档说明
func (c *SmsTemplateController) Prepare() {
    //先执行
    c.BaseController.Prepare()
    //如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
    c.checkAuthor("DataGrid", "DataList", "Edit", "Delete")
    //如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
    //权限控制里会进行登录验证，因此这里不用再作登录验证
    //c.checkLogin()
}

// DataGrid 角色管理首页 表格获取数据
func (c *SmsTemplateController) DataGrid() {
    //直接反序化获取json格式的requestbody里的值
    var params models.SmsTemplateQueryParam
    json.Unmarshal(c.Ctx.Input.RequestBody, &params)
    //获取数据列表和总数
    if !c.curUser.IsSuper {
        params.UserCompanyId = c.curUser.UserCompanyId
    }
    data, total := models.SmsTemplatePageList(&params)
    //定义返回的数据结构
    result := make(map[string]interface{})
    result["total"] = total
    result["rows"] = data
    c.Data["json"] = result
    c.ServeJSON()
}

// DataList 角色列表
func (c *SmsTemplateController) DataList() {
    var params = models.SmsTemplateQueryParam{}
    //获取数据列表和总数
    data := models.SmsTemplateDataList(&params)
    //定义返回的数据结构
    c.jsonResult(enums.JRCodeSucc, "", data)
}

// Edit 添加、编辑角色界面
func (c *SmsTemplateController) Edit() {
    if c.Ctx.Request.Method == "POST" {
        c.Save()
    }
    Id, _ := c.GetInt(":id", 0)
    m := models.SmsTemplate{Id: Id}
    if Id > 0 {
        o := orm.NewOrm()
        err := o.Read(&m)
        if err != nil {
            c.pageError("数据无效，请刷新后重试")
        }
    }
    c.Data["m"] = m
    c.setTpl("smstemplate/edit.html", "shared/layout_pullbox.html")
    c.LayoutSections = make(map[string]string)
    c.LayoutSections["footerjs"] = "smstemplate/edit_footerjs.html"
}

//Save 添加、编辑页面 保存
func (c *SmsTemplateController) Save() {
    var err error
    m := models.SmsTemplate{}
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
        m.BackendUserId = c.curUser.UserCompanyId
        if _, err = o.Insert(&m); err == nil {
            c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
        } else {
            c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
        }

    } else {
        m.Updated = time.Now().Format(enums.BaseTimeFormat)
        // password为空不更新，不为空更新
        fields := []string{"Title", "Content", "Classify", "Updated"}
        if _, err = o.Update(&m, fields...); err == nil {
            c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
        } else {
            c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
        }
    }

}

//Delete 批量删除
func (c *SmsTemplateController) Delete() {
    strs := c.GetString("ids")
    ids := make([]int, 0, len(strs))
    for _, str := range strings.Split(strs, ",") {
        if id, err := strconv.Atoi(str); err == nil {
            ids = append(ids, id)
        }
    }
    if num, err := models.SmsTemplateBatchDelete(ids); err == nil {
        c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
    } else {
        c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
    }
}
