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

//UserCompanyController 角色管理
type UserCompanyController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *UserCompanyController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid", "DataList")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

//Index 角色管理首页
func (c *UserCompanyController) Index() {
	//是否显示更多查询条件的按钮
	c.Data["showMoreQuery"] = false
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "usercompany/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "usercompany/index_footerjs.html"
	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("UserCompanyController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("UserCompanyController", "Delete")
}

// DataGrid 角色管理首页 表格获取数据
func (c *UserCompanyController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.UserCompanyQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	data, total := models.UserCompanyPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

// DataList 角色列表
func (c *UserCompanyController) DataList() {
	var params = models.UserCompanyQueryParam{}
	//获取数据列表和总数
	data := models.UserCompanyDataList(&params)
	//定义返回的数据结构
	c.jsonResult(enums.JRCodeSucc, "", data)
}

// Edit 添加、编辑角色界面
func (c *UserCompanyController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := models.UserCompany{Id: Id}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["m"] = m
	c.setTpl("usercompany/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "usercompany/edit_footerjs.html"
}

//Save 添加、编辑页面 保存
func (c *UserCompanyController) Save() {
	var err error
	m := models.UserCompany{}
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	o := orm.NewOrm()
	if m.Id == 0 {
		// insert添加created时间
		m.Created = time.Now().Format(enums.BaseTimeFormat)
		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}

	} else {
		if _, err = o.Update(&m, "Name", "NameAbbr", "GatewayGroupId", "LimitDial", "LimitCaller"); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}

}

//Delete 批量删除
func (c *UserCompanyController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if num, err := models.UserCompanyBatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}
