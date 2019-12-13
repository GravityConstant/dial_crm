package controllers

import (
    "encoding/json"

    "zq/callout_crm/enums"
    "zq/callout_crm/models"
    "zq/callout_crm/utils"
    "zq/callout_crm/dial"

    "fmt"
    "strconv"
    "strings"
    "time"

    "github.com/astaxie/beego/orm"
)

//GatewayPhoneController 角色管理
type GatewayPhoneController struct {
    BaseController
}

//Prepare 参考beego官方文档说明
func (c *GatewayPhoneController) Prepare() {
    //先执行
    c.BaseController.Prepare()
    //如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
    c.checkAuthor("DataGrid", "DataList")
    //如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
    //权限控制里会进行登录验证，因此这里不用再作登录验证
    //c.checkLogin()
}

//Index 角色管理首页
func (c *GatewayPhoneController) Index() {
    //是否显示更多查询条件的按钮
    c.Data["showMoreQuery"] = false
    //将页面左边菜单的某项激活
    c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
    c.setTpl()
    c.LayoutSections = make(map[string]string)
    c.LayoutSections["headcssjs"] = "gatewayphone/index_headcssjs.html"
    c.LayoutSections["footerjs"] = "gatewayphone/index_footerjs.html"
    //页面里按钮权限控制
    c.Data["canEdit"] = c.checkActionAuthor("GatewayPhoneController", "Edit")
    c.Data["canDelete"] = c.checkActionAuthor("GatewayPhoneController", "Delete")
}

// DataGrid 角色管理首页 表格获取数据
func (c *GatewayPhoneController) DataGrid() {
    //直接反序化获取json格式的requestbody里的值
    var params models.GatewayPhoneQueryParam
    json.Unmarshal(c.Ctx.Input.RequestBody, &params)
    //获取数据列表和总数
    data, total := models.GatewayPhonePageList(&params)
    //定义返回的数据结构
    result := make(map[string]interface{})
    result["total"] = total
    result["rows"] = data
    c.Data["json"] = result
    c.ServeJSON()
}

// DataList 角色列表
func (c *GatewayPhoneController) DataList() {
    var params = models.GatewayPhoneQueryParam{}
    //获取数据列表和总数
    data := models.GatewayPhoneDataList(&params)
    //定义返回的数据结构
    c.jsonResult(enums.JRCodeSucc, "", data)
}

// Edit 添加、编辑角色界面
func (c *GatewayPhoneController) Edit() {
    if c.Ctx.Request.Method == "POST" {
        c.Save()
    }
    Id, _ := c.GetInt(":id", 0)
    m := models.GatewayPhone{Id: Id}
    if Id > 0 {
        o := orm.NewOrm()
        err := o.Read(&m)
        if err != nil {
            c.pageError("数据无效，请刷新后重试")
        }
    }
    c.Data["m"] = m
    c.setTpl("gatewayphone/edit.html", "shared/layout_pullbox.html")
    c.LayoutSections = make(map[string]string)
    c.LayoutSections["footerjs"] = "gatewayphone/edit_footerjs.html"
}

//Save 添加、编辑页面 保存
func (c *GatewayPhoneController) Save() {
    var err error
    m := models.GatewayPhone{}
    //获取form里的值
    if err = c.ParseForm(&m); err != nil {
        c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
    }
    counts, _ := c.GetInt("PhoneNo", 0)

    o := orm.NewOrm()
    if m.Id == 0 {
        // insert添加created时间
        m.Created = time.Now().Format(enums.BaseTimeFormat)
        // 批量添加
        gps := make([]models.GatewayPhone, 0)
        p := m.Phone
        for i:=0; i<counts; i++ {
            tmp := models.GatewayPhone{
                Phone    : m.Phone,
                GatewayId: m.GatewayId,
                AgentId  : m.AgentId,
                Created  : m.Created,
            }
            gps = append(gps, tmp)
            m.Phone = utils.StringAddInt(p, int64(i+1))
        }
        if _, err = o.InsertMulti(len(gps), gps); err == nil {
            c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
        } else {
            c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
        }

    } else {
        if _, err = o.Update(&m, "Phone", "GatewayId", "AgentId"); err == nil {
            c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
        } else {
            c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
        }
    }

}

//Delete 批量删除
func (c *GatewayPhoneController) Delete() {
    strs := c.GetString("ids")
    ids := make([]int, 0, len(strs))
    for _, str := range strings.Split(strs, ",") {
        if id, err := strconv.Atoi(str); err == nil {
            ids = append(ids, id)
        }
    }
    if num, err := models.GatewayPhoneBatchDelete(ids); err == nil {
        c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
    } else {
        c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
    }
}

// Delete 批量分配坐席
func (c *GatewayPhoneController) AllocateAgent() {
    strs := c.GetString("ids")
    agentId, _ := c.GetInt("agent_id", 0)

    ids := make([]int, 0, len(strs))
    for _, str := range strings.Split(strs, ",") {
        if id, err := strconv.Atoi(str); err == nil {
            ids = append(ids, id)
        }
    }
    if num, err := models.GatewayPhoneBatchAllocateAgent(ids, agentId); err == nil {
        // 更新dial库中的AgentCalloutParam
        var params models.GatewayPhoneQueryParam
        params.Ids = ids
        list := models.GatewayPhoneDataList(&params)
        dial.CalloutUpdate(list)
        c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功分配 %d 项", num), 0)
    } else {
        c.jsonResult(enums.JRCodeFailed, "分配失败", 0)
    }
}
