package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"zq/callout_crm/enums"
	"zq/callout_crm/models"
)

type UserServerHistoryController struct {
	BaseController
}

func (c *UserServerHistoryController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("AddRecord","RecordList","RecordListDataGrid")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

}


func (c *UserServerHistoryController) AddRecord() {
	//如果是Post请求，则由Save处理
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	State, _ := c.GetInt(":state", 0)
	c.Data["Id"] = Id
	c.Data["BackendUserId"] = c.curUser.Id
	c.Data["BackendUserName"] = c.curUser.RealName

	c.Data["State"] = State

	c.setTpl("userserverhistory/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "userserverhistory/edit_footerjs.html"
}


func (c *UserServerHistoryController) Save() {
	var err error
	m := models.UserServerHistory{}
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}

	uc := models.UserClient{}
	if err = c.ParseForm(&uc); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", uc.State)
	}
	uc.Id=m.UserClientId


	o := orm.NewOrm()

	o.Begin()

	if _, err = o.Insert(&m); err != nil  {
		o.Rollback()
		c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
	}

	if _, err = o.Update(&uc,"State"); err != nil {
		o.Rollback()
		c.jsonResult(enums.JRCodeFailed, "添加失败", uc.Id)
	}
	o.Commit()
	c.jsonResult(enums.JRCodeSucc, "添加成功", uc.Id)

}



func (c *UserServerHistoryController) RecordList() {
	//如果是Post请求，则由Save处理

	Id, _ := c.GetInt(":id", 0)

	c.Data["UserClientId"] = Id


	c.setTpl("userserverhistory/rcordlist.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "userserverhistory/rcordlist_headcssjs.html"
	c.LayoutSections["footerjs"] = "userserverhistory/rcordlist_footerjs.html"
}

func (c *UserServerHistoryController) RecordListDataGrid() {
	//直接反序化获取json格式的requestbody里的值（要求配置文件里 copyrequestbody=true）
	var params models.UserServerHistoryQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	data, total := models.UserServerHistoryList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}