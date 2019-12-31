package controllers

import (
	"encoding/json"
    "fmt"
    "strconv"
    "strings"

	"zq/callout_crm/enums"
	"zq/callout_crm/models"
    "zq/callout_crm/utils"

	"github.com/astaxie/beego/orm"
)

//AgentController 角色管理
type AgentController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *AgentController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid", "DataList")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

//Index 角色管理首页
func (c *AgentController) Index() {
	//是否显示更多查询条件的按钮
	c.Data["showMoreQuery"] = false
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "agent/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "agent/index_footerjs.html"
	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("AgentController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("AgentController", "Delete")
}

// DataGrid 角色管理首页 表格获取数据
func (c *AgentController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.AgentQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	if !c.curUser.IsSuper {
		params.UIds = c.curUser.UIds
	}

	//获取数据列表和总数
	data, total := models.AgentPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

//DataList 角色列表
func (c *AgentController) DataList() {
	var params = models.AgentQueryParam{}
	//获取数据列表和总数
	data := models.AgentDataList(&params)
	//定义返回的数据结构
	c.jsonResult(enums.JRCodeSucc, "", data)
}

//Edit 添加、编辑角色界面
func (c *AgentController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := models.Agent{Id: Id}

	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["m"] = m
	c.setTpl("agent/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "agent/edit_footerjs.html"
}

//Save 添加、编辑页面 保存
func (c *AgentController) Save() {
	var err error
	m := models.Agent{}
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	o := orm.NewOrm()
	if m.BackendUserId > 0 {
		if bc, err := models.BackendUserCompanyByUserId(m.BackendUserId); err == nil {
			m.GatewayId = bc.GatewayId
		}
	}
	if m.Id == 0 {
		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}
	} else {
		if _, err = o.Update(&m, "ExtNo", "ExtPwd", "GatewayPhoneNumber", "OriginationCallerIdNumber", "BackendUserId", "GatewayId"); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}

}

//EditV1 新版本：4个号码轮询，插手机卡、中继2个
func (c *AgentController) EditV1() {
	if c.Ctx.Request.Method == "POST" {
		c.SaveV1()
	}
	Id, _ := c.GetInt(":id", 0)
	m := models.Agent{Id: Id}

	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["m"] = m
	c.setTpl("agent/editv1.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "agent/editv1_footerjs.html"
}

//Save 添加、编辑页面 保存
func (c *AgentController) SaveV1() {
	var err error
	m := models.Agent{}
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
    counts, _ := c.GetInt("PhoneNo", 0)

	o := orm.NewOrm()
	if m.Id == 0 {
        // 批量添加
        gps := make([]models.Agent, 0)
        p := m.ExtNo
        for i:=0; i<counts; i++ {
            tmp := models.Agent{
                ExtNo         : m.ExtNo,
                ExtPwd        : m.ExtPwd,
                BackendUserId : m.BackendUserId,
                DefaultTrunk  : m.DefaultTrunk,
            }
            gps = append(gps, tmp)
            m.ExtNo = utils.StringAddInt(p, int64(i+1))
        }
		if _, err = o.InsertMulti(len(gps), gps); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}
	} else {
		if _, err = o.Update(&m, "ExtNo", "ExtPwd", "BackendUserId", "DefaultTrunk", "BindPhone"); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}

}

//Delete 批量删除
func (c *AgentController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if num, err := models.AgentBatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}

// 座席号码取自freeswitch的配置文件conf/directory/default下的文件名
func (c *AgentController) FsRegisterUserList() {
    id, _ := c.GetInt("id", 0)

    data := make([]*models.FsRegisterUser, 0)
    if id == 0 {
        data = models.FsUserListNotInDbList()
    } else {
        data = models.FsRegisterUserList()
    }
	

	c.jsonResult(enums.JRCodeSucc, "", data)
}
