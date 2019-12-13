package controllers

import (
	"encoding/json"
	"path/filepath"
	"fmt"
	"os"
	"strconv"
	"strings"

	"zq/callout_crm/enums"
	"zq/callout_crm/models"

	"github.com/astaxie/beego/orm"
)

type GatewayController struct {
	BaseController
}

func (c *GatewayController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid", "DataList")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

}

func (c *GatewayController) Index() {
	//是否显示更多查询条件的按钮
	c.Data["showMoreQuery"] = false
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	//页面模板设置
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "gateway/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "gateway/index_footerjs.html"
	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("GatewayController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("GatewayController", "Delete")
}

func (c *GatewayController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值（要求配置文件里 copyrequestbody=true）
	var params models.GatewayQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//获取数据列表和总数
	data, total := models.GatewayPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

// DataList 列表
func (c *GatewayController) DataList() {
	var params = models.GatewayQueryParam{}
	//获取数据列表和总数
	data := models.GatewayDataList(&params)
	//定义返回的数据结构
	c.jsonResult(enums.JRCodeSucc, "", data)
}

// 获取当前用户所在公司网关的列表
func (c *GatewayController) DataListByUserCompany() {
    uId, _ := c.GetInt("uid", 0)

    //获取数据列表和总数
    data := models.DataListByUserCompany(uId)
    //定义返回的数据结构
    c.jsonResult(enums.JRCodeSucc, "", data)
}

//Edit 添加、编辑角色界面
func (c *GatewayController) Edit() {
    if c.Ctx.Request.Method == "POST" {
        c.Save()
    }
    Id, _ := c.GetInt(":id", 0)
    m := models.Gateway{Id: Id}
    if Id > 0 {
    	o := orm.NewOrm()
        err := o.Read(&m)
        if err != nil {
            c.pageError("数据无效，请刷新后重试")
        }
		o.LoadRelated(&m, "GatewayGroupMap")
    }
    c.Data["m"] = m
    // 获取关联的gateway_group_id列表
	var gatewayGroupIds []string
	for _, item := range m.GatewayGroupMap {
		gatewayGroupIds = append(gatewayGroupIds, strconv.Itoa(item.GatewayGroup.Id))
	}
	c.Data["gatewayGroupIds"] = strings.Join(gatewayGroupIds, ",")
    c.setTpl("gateway/edit.html", "shared/layout_pullbox.html")
    c.LayoutSections = make(map[string]string)
    c.LayoutSections["footerjs"] = "gateway/edit_footerjs.html"
}

//Save 添加、编辑页面 保存
func (c *GatewayController) Save() {
    var err error
    o := orm.NewOrm()
    m := models.Gateway{}
    //获取form里的值
    if err = c.ParseForm(&m); err != nil {
        c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
    }
    //删除已关联的历史数据
	if _, err := o.QueryTable(models.GatewayGroupMapTBName()).Filter("gateway__id", m.Id).Delete(); err != nil {
		c.jsonResult(enums.JRCodeFailed, "删除历史关系失败", "")
	}
    if m.Id == 0 {
        if _, err = o.Insert(&m); err != nil {
            c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
        }
    } else {
        if _, err := o.Update(&m); err != nil {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
    }
    //添加关系
	var relations []models.GatewayGroupMap
	for _, groupId := range m.GatewayGroupIds {
		r := models.GatewayGroup{Id: groupId}
		relation := models.GatewayGroupMap{Gateway: &m, GatewayGroup: &r}
		relations = append(relations, relation)
	}
	if len(relations) > 0 {
		//批量添加
		if _, err := o.InsertMulti(len(relations), relations); err == nil {
			c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "保存失败", m.Id)
		}
	} else {
		c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
	}
}

//Delete 批量删除
func (c *GatewayController) Delete() {
    strs := c.GetString("ids")
    ids := make([]int, 0, len(strs))
    for _, str := range strings.Split(strs, ",") {
        if id, err := strconv.Atoi(str); err == nil {
            ids = append(ids, id)
        }
    }
    if num, err := models.GatewayBatchDelete(ids); err == nil {
        c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
    } else {
        c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
    }
}

// 网关名取自freeswitch的配置文件/usr/local/freeswitch/conf/sip_profiles/external下的文件名
func (c *GatewayController) FsRegisterGatewayList() {
	data := make([]*models.FsRegisterGateway, 0)
	filepath.Walk("/usr/local/freeswitch/conf/sip_profiles/external", func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			filename := strings.Split(info.Name(), ".")
			fsFile := models.FsRegisterGateway{
				Name: filename[0],
			}
			data = append(data, &fsFile)
		}
		return nil
	})
	c.jsonResult(enums.JRCodeSucc, "", data)
}