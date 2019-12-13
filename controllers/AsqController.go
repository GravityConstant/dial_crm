package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
	"zq/callout_crm/enums"
	"zq/callout_crm/models"
)

type AsqController struct {
	BaseController
}

func (c *AsqController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

}
func (c *AsqController) Index() {
	//是否显示更多查询条件的按钮
	c.Data["showMoreQuery"] = true
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	//页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "asq/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "asq/index_footerjs.html"
	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("AsqController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("AsqController", "Delete")
}

func (c *AsqController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值（要求配置文件里 copyrequestbody=true）
	var params models.AsqQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	data, total := models.AsqPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

// Edit 添加 编辑 页面
func (c *AsqController) Edit() {
	//如果是Post请求，则由Save处理
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := &models.Asq{}
	var err error
	if Id > 0 {
		m, err = models.AsqOne(Id)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["m"] = m

	c.setTpl("asq/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "asq/edit_footerjs.html"
}

func (c *AsqController) Save() {
	var err error
	m := models.Asq{}
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	o := orm.NewOrm()
	if m.Id == 0 {
		m.BackendUserId = c.curUser.Id
		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}

	} else {
		m.Updated = time.Now()
		if _, err = o.Update(&m,"Name","Description","Updated"); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}
}


func (c *AsqController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	query := orm.NewOrm().QueryTable(models.AsqTBName())
	if num, err := query.Filter("id__in", ids).Delete(); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}


// Edit 添加 编辑 页面
func (c *AsqController) Set() {
	//如果是Post请求，则由Save处理
	if c.Ctx.Request.Method == "POST" {
		//c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := &models.Asq{}
	var err error
	if Id > 0 {
		m, err = models.AsqOne(Id)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["m"] = m

	c.setTpl("asq/set.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "asq/set_footerjs.html"
}