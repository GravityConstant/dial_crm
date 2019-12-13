package controllers

import (
	"encoding/json"

	"zq/callout_crm/enums"
	"zq/callout_crm/models"
	"zq/callout_crm/utils"

	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

//UserClientFieldController 管理
type UserClientFieldController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *UserClientFieldController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid", "DataList", "UpdateField")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

//Index 角色管理首页
func (c *UserClientFieldController) Index() {
	//是否显示更多查询条件的按钮
	c.Data["showMoreQuery"] = false
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "userclientfield/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "userclientfield/index_footerjs.html"
	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("UserClientFieldController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("UserClientFieldController", "Delete")
}

// DataGrid 角色管理首页 表格获取数据
func (c *UserClientFieldController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.UserClientFieldQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	if !(c.curUser.Id == 1 || c.curUser.UserName == "admin") {
		params.UserCompanyId = c.curUser.UserCompanyId
	}

	data, total := models.UserClientFieldPageList(&params)
	// utils.Info("user client field: ")
	// for _, item := range data {
	// 	utils.Info("%v\n", item)
	// }
	// utils.Info("%v\n", c.curUser)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

//DataList 角色列表
func (c *UserClientFieldController) DataList() {
	var params = models.UserClientFieldQueryParam{}
	params.UserCompanyId = c.curUser.UserCompanyId
	//获取数据列表和总数
	data := models.UserClientFieldDataList(&params)
	//定义返回的数据结构
	c.jsonResult(enums.JRCodeSucc, "", data)
}

//Edit 添加、编辑角色界面
func (c *UserClientFieldController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := models.UserClientField{Id: Id}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["m"] = m
	c.setTpl("userclientfield/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "userclientfield/edit_footerjs.html"
}

//Save 添加、编辑页面 保存
func (c *UserClientFieldController) Save() {
	var err error
	m := models.UserClientField{}
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	o := orm.NewOrm()
	// 对字段类型做特殊处理
	switch m.FieldType {
	case 2, 3:
		ftvs := strings.Split(strings.TrimSpace(m.FieldTypeValue), "\r\n")
		ftvsTrim := make([]string, 0)
		for _, v := range ftvs {
			if len(v) > 0 {
				ftvsTrim = append(ftvsTrim, strings.TrimSpace(v))
			}
		}
		m.FieldTypeValue = strings.Join(ftvsTrim, ",")
	default:
		m.FieldTypeValue = ""
	}
	if m.Id == 0 {
		// 所属公司
		if !(m.UserCompanyId > 0) {
			m.UserCompanyId = c.curUser.UserCompanyId
		}
		// 真实对应的列名
		if m.ColumnName, err = models.UserClientFieldValidColumnName(m.UserCompanyId); err != nil {
			utils.LogError("编辑用户自定义字段，获取有效列名时发生错误。" + err.Error())
			c.jsonResult(enums.JRCodeFailed, "添加失败,"+err.Error(), m.Id)
		}

		// 自定义列
		m.FieldSpecies = 0

		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			utils.LogError("添加用户字段失败。" + err.Error())
			c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}

	} else {
		if _, err = o.Update(&m, "FieldName", "FieldType", "FieldTypeValue", "ListShow", "AddShow", "QueryShow"); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			utils.LogError("更新用户字段失败。" + err.Error())
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}

}

// Delete 批量删除
func (c *UserClientFieldController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if num, err := models.UserClientFieldBatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}

// 更新字段信息
func (c *UserClientFieldController) UpdateField() {
	id, _ := c.GetInt("id", 0)
	field := c.GetString("field")
	value := c.GetString("value")

	if id == 0 {
		c.jsonResult(enums.JRCodeFailed, "未知记录", 0)
	}

	query := orm.NewOrm().QueryTable(models.UserClientFieldTBName())
	if _, err := query.Filter("id", id).Update(orm.Params{
		utils.SnakeString(field): value,
	}); err != nil {
		c.jsonResult(enums.JRCodeFailed, "更新"+field+"失败", 0)
	} else {
		c.jsonResult(enums.JRCodeSucc, "更新"+field+"成功", 0)
	}

}
