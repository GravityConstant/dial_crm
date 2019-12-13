package controllers

import (
	"encoding/json"

	"zq/callout_crm/enums"
	"zq/callout_crm/models"
	"zq/callout_crm/utils"

	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/tealeg/xlsx"
)

//UserClientController 管理
type UserClientController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *UserClientController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid", "DataList", "UpdateFieldByIds", "DownloadTmpl", "UploadExcel", "Download", "IsExistedCalledPhone")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

//Index 角色管理首页
func (c *UserClientController) Index() {
	//是否显示更多查询条件的按钮
	c.Data["showMoreQuery"] = true
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "userclient/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "userclient/index_footerjs.html"
	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("UserClientController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("UserClientController", "Delete")
}

// DataGrid 角色管理首页 表格获取数据
func (c *UserClientController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.UserClientQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	// 根据uid获取记录
	if !c.curUser.IsSuper {
		params.CompanyId = c.curUser.UserCompanyId
		if len(params.BackendUserIds) > 0 {
			if params.BackendUserIds[0] == -1 {
				params.BackendUserIds = c.curUser.UIds
			}
		} else {
			params.BackendUserIds = []int{c.curUser.Id}
		}
	}

	//获取数据列表和总数
	data, total := models.UserClientPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

//DataList 角色列表
func (c *UserClientController) DataList() {
	var params = models.UserClientQueryParam{}
	//获取数据列表和总数
	data := models.UserClientDataList(&params)
	//定义返回的数据结构
	c.jsonResult(enums.JRCodeSucc, "", data)
}

//Edit 添加、编辑角色界面
func (c *UserClientController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	phone := c.GetString(":phone")
	dialSuccess, _ := c.GetBool(":dialSuccess", false)
	m := models.UserClient{Id: Id}
	if utils.IsPhoneNumber(phone) {
		m.MobilePhone = phone
	} else {
		m.MobilePhone = ""
	}
	
	// 获取自定义的字段信息
	var params models.UserClientFieldQueryParam
	params.UserCompanyId = c.curUser.UserCompanyId
	attrs := models.UserClientFieldDataList(&params)

	// 整理数据
	ma := make(map[string]*models.UserClientField)
	for _, item := range attrs {
		switch item.ColumnName {
		case "name":
			ma["Name"] = item
		case "mobile_phone":
			ma["MobilePhone"] = item
		case "contact_phone":
			ma["ContactPhone"] = item
		case "comment":
			ma["Comment"] = item
		case "address":
			ma["Address"] = item
		case "dial_state":
			ma["DialState"] = item
		case "state":
			ma["State"] = item
		case "feature":
			ma["Feature"] = item
		case "complaint":
			ma["Complaint"] = item
		case "clue_from":
			ma["ClueFrom"] = item
		case "email":
			ma["Email"] = item
		case "column1":
			ma["Column1"] = item
		case "column2":
			ma["Column2"] = item
		case "column3":
			ma["Column3"] = item
		case "column4":
			ma["Column4"] = item
		case "column5":
			ma["Column5"] = item
		case "column6":
			ma["Column6"] = item
		case "column7":
			ma["Column7"] = item
		case "column8":
			ma["Column8"] = item
		case "column9":
			ma["Column9"] = item
		case "column10":
			ma["Column10"] = item
		case "column11":
			ma["Column11"] = item
		case "column12":
			ma["Column12"] = item
		case "column13":
			ma["Column13"] = item
		case "column14":
			ma["Column14"] = item
		case "column15":
			ma["Column15"] = item
		case "column16":
			ma["Column16"] = item
		}
	}

	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["m"] = m
	c.Data["ma"] = ma
	c.Data["dialSuccess"] = dialSuccess
	c.setTpl("userclient/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "userclient/edit_footerjs.html"
}

//Save 添加、编辑页面 保存.199:3122731;200:3122807
func (c *UserClientController) Save() {
	var err error
	m := models.UserClient{}
	um := models.UserServerHistory{}
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		utils.Info("%v\n", err)
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	// 呼叫成功更新latest_communicated字段
	dialSuccess, _ := c.GetBool("dialSuccess", false)
	// 服务记录
	um.Context = c.GetString("Context")
	um.UserClientId = m.Id
	um.BackendUserId = c.curUser.Id
	um.BackendUserName = c.curUser.RealName
	// 获取客户字段属性
	params := models.UserClientFieldQueryParam{UserCompanyId: c.curUser.UserCompanyId}
	attrs := models.UserClientFieldDataList(&params)
	for _, attr := range attrs {
		switch attr.ColumnName {
		case "column1":
			if attr.FieldType == 3 {
				m.Column1 = strings.Join(c.GetStrings("Column1"), ",")
			}
		case "column2":
			if attr.FieldType == 3 {
				m.Column2 = strings.Join(c.GetStrings("Column2"), ",")
			}
		case "column3":
			if attr.FieldType == 3 {
				m.Column3 = strings.Join(c.GetStrings("Column3"), ",")
			}
		case "column4":
			if attr.FieldType == 3 {
				m.Column4 = strings.Join(c.GetStrings("Column4"), ",")
			}
		case "column5":
			if attr.FieldType == 3 {
				m.Column5 = strings.Join(c.GetStrings("Column5"), ",")
			}
		case "column6":
			if attr.FieldType == 3 {
				m.Column6 = strings.Join(c.GetStrings("Column6"), ",")
			}
		case "column7":
			if attr.FieldType == 3 {
				m.Column7 = strings.Join(c.GetStrings("Column7"), ",")
			}
		case "column8":
			if attr.FieldType == 3 {
				m.Column8 = strings.Join(c.GetStrings("Column8"), ",")
			}
		case "column9":
			if attr.FieldType == 3 {
				m.Column9 = strings.Join(c.GetStrings("Column9"), ",")
			}
		case "column10":
			if attr.FieldType == 3 {
				m.Column10 = strings.Join(c.GetStrings("Column10"), ",")
			}
		case "column11":
			if attr.FieldType == 3 {
				m.Column11 = strings.Join(c.GetStrings("Column11"), ",")
			}
		case "column12":
			if attr.FieldType == 3 {
				m.Column12 = strings.Join(c.GetStrings("Column12"), ",")
			}
		case "column13":
			if attr.FieldType == 3 {
				m.Column13 = strings.Join(c.GetStrings("Column13"), ",")
			}
		case "column14":
			if attr.FieldType == 3 {
				m.Column14 = strings.Join(c.GetStrings("Column14"), ",")
			}
		case "column15":
			if attr.FieldType == 3 {
				m.Column15 = strings.Join(c.GetStrings("Column15"), ",")
			}
		case "column16":
			if attr.FieldType == 3 {
				m.Column16 = strings.Join(c.GetStrings("Column16"), ",")
			}
		}
	}

	// utils.Info("UserClientController, Save")
	// fmt.Println(m)
	
	// c.jsonResult(enums.JRCodeFailed, "testing...", m.Id)

	o := orm.NewOrm()
	o.Begin()
	if m.Id == 0 {
		// 默认字段设置
		m.BackendUserId = c.curUser.Id
		m.BelongBackendUserId = c.curUser.Id
		m.UserCompanyId = c.curUser.UserCompanyId
		m.Created = time.Now().Format(enums.BaseTimeFormat)
		m.Updated = m.Created
		if dialSuccess {
			m.LatestCommunicated = m.Created
		} else {
			m.LatestCommunicated = `2000-01-01 00:00:00`
		}
		if _, err = o.Insert(&m); err == nil {
			um.UserClientId = m.Id
			if dialSuccess {
				if _, err = o.Insert(&um); err != nil  {
					o.Rollback()
					c.jsonResult(enums.JRCodeFailed, "添加客户资料成功，添加服务记录失败", m.Id)
				}
			}
			o.Commit()
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			o.Rollback()
			utils.LogError("添加客户信息失败。" + err.Error())
			c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}

	} else {
		m.Updated = time.Now().Format(enums.BaseTimeFormat)
		var updatedColumns []string
		if dialSuccess {
			m.LatestCommunicated = m.Updated
			updatedColumns = []string{"Name", "MobilePhone", "ContactPhone", "Comment", "Address", "DialState", "State", "Feature", "Complaint", "ClueFrom", "Email", "Updated", "LatestCommunicated", "Column1", "Column2", "Column3", "Column4", "Column5", "Column6", "Column7", "Column8", "Column9", "Column10", "Column11", "Column12", "Column13", "Column14", "Column15", "Column16"}
		} else {
			updatedColumns = []string{"Name", "MobilePhone", "ContactPhone", "Comment", "Address", "DialState", "State", "Feature", "Complaint", "ClueFrom", "Email", "Updated", "Column1", "Column2", "Column3", "Column4", "Column5", "Column6", "Column7", "Column8", "Column9", "Column10", "Column11", "Column12", "Column13", "Column14", "Column15", "Column16"}
		}

		if _, err = o.Update(&m, updatedColumns...); err == nil {
			if dialSuccess {
				if _, err = o.Insert(&um); err != nil  {
					o.Rollback()
					c.jsonResult(enums.JRCodeFailed, "更新客户资料成功，添加服务记录失败", m.Id)
				}
			}
			o.Commit()
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			o.Rollback()
			utils.LogError("编辑客户信息失败。" + err.Error())
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}

}

// Delete 批量删除
func (c *UserClientController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	if num, err := models.UserClientBatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}

// 根据ids修改一个字段信息
func (c *UserClientController) UpdateFieldByIds() {
	strs := c.GetString("ids")
	fieldName := c.GetString("fieldName")
	fieldValue := c.GetString("fieldValue")

	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	query := orm.NewOrm().QueryTable(models.UserClientTBName())
	if num, err := query.Filter("id__in", ids).Update(orm.Params{
		fieldName: fieldValue,
	}); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功更新 %d 项", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "更新字段失败", 0)
	}
}

// 下载导入模板
func (c *UserClientController) DownloadTmpl() {
	var err error
	m := models.UserClient{}
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		utils.LogError("下载导入模板，获取参数失败。" + err.Error())
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	utils.Info("%v\n", m)
	// c.StopRun()

	// 读取字段信息
	var params models.UserClientFieldQueryParam
	params.UserCompanyId = c.curUser.UserCompanyId
	fieldsInfo := models.UserClientFieldDataList(&params)

	f := xlsx.NewFile()
	sheet, _ := f.AddSheet("sheet1")
	th := sheet.AddRow()

	var val string
	if strings.Compare(m.Name, "on") == 0 {
		if val, err = getFieldName("name", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.MobilePhone, "on") == 0 {
		if val, err = getFieldName("mobile_phone", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.ContactPhone, "on") == 0 {
		if val, err = getFieldName("contract_phone", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Comment, "on") == 0 {
		if val, err = getFieldName("comment", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Address, "on") == 0 {
		if val, err = getFieldName("address", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.DialState, "on") == 0 {
		if val, err = getFieldName("dial_state", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if m.State == 1 {
		if val, err = getFieldName("state", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Complaint, "on") == 0 {
		if val, err = getFieldName("complaint", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Email, "on") == 0 {
		if val, err = getFieldName("email", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column1, "on") == 0 {
		if val, err = getFieldName("column1", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column2, "on") == 0 {
		if val, err = getFieldName("column2", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column3, "on") == 0 {
		if val, err = getFieldName("column3", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column4, "on") == 0 {
		if val, err = getFieldName("column4", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column5, "on") == 0 {
		if val, err = getFieldName("column5", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column6, "on") == 0 {
		if val, err = getFieldName("column6", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column7, "on") == 0 {
		if val, err = getFieldName("column7", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column8, "on") == 0 {
		if val, err = getFieldName("column8", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column9, "on") == 0 {
		if val, err = getFieldName("column9", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column10, "on") == 0 {
		if val, err = getFieldName("column10", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column11, "on") == 0 {
		if val, err = getFieldName("column11", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column12, "on") == 0 {
		if val, err = getFieldName("column12", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column13, "on") == 0 {
		if val, err = getFieldName("column13", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column14, "on") == 0 {
		if val, err = getFieldName("column14", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column15, "on") == 0 {
		if val, err = getFieldName("column15", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}
	if strings.Compare(m.Column16, "on") == 0 {
		if val, err = getFieldName("column16", fieldsInfo); err == nil {
			cell := th.AddCell()
			cell.SetString(val)
		}
	}

	filePath := "./static/download/客户资料模板.xlsx"
	if err := os.MkdirAll("static/download", 0666); err != nil {
		utils.LogError("创建下载文件夹失败。" + err.Error())
		// c.jsonResult(enums.JRCodeFailed, "创建下载文件夹失败。", 0)
	}
	if err := f.Save(filePath); err == nil {
		c.Ctx.Output.Download(filePath, "客户资料模板.xlsx")
	} else {
		utils.LogError("下载客户资料模板失败。" + err.Error())
		// c.jsonResult(enums.JRCodeFailed, "下载客户资料模板失败。", 0)
	}
	c.StopRun()
}

// 导入模板设置单元格值
func getFieldName(columnName string, fieldsInfo []*models.UserClientField) (string, error) {
	for _, item := range fieldsInfo {
		if strings.Compare(item.ColumnName, columnName) == 0 {
			switch item.FieldType {
			case 2, 3:
				format := "%s (%s)"
				return fmt.Sprintf(format, item.FieldName, item.FieldTypeValue), nil
			default:
				return item.FieldName, nil
			}
		}
	}
	return "", errors.New("invalid column")
}

// 上传excel文件导入表中
func (c *UserClientController) UploadExcel() {
	f, _, err := c.GetFile("custtemplet")
	if err != nil {
		c.jsonResult(enums.JRCodeFailed, "上传失败", "")
	}
	defer f.Close()

	//文件指针指向文件末尾 获取文件大小保存于buf_len
	buf_len, _ := f.Seek(0, os.SEEK_END)
	fmt.Println("buf_len", buf_len)
	//获取buf_len后把文件指针重新定位于文件开始
	f.Seek(0, os.SEEK_SET)

	buf := make([]byte, buf_len)
	f.Read(buf)
	xfile, err := xlsx.OpenBinary(buf)
	if err != nil {
		utils.LogError("上传打开客户文件失败。" + err.Error())
		c.jsonResult(enums.JRCodeFailed, "上传打开客户文件失败。", "")
	}

	// 读取字段信息
	var params models.UserClientFieldQueryParam
	params.UserCompanyId = c.curUser.UserCompanyId
	fieldsInfo := models.UserClientFieldDataList(&params)

	excelCol2structCol := make(map[int]*models.UserClientField)
	utils.Info("info...")
	sheet := xfile.Sheets[0]
	for k, th := range sheet.Rows[0].Cells {
		if tmp, err := excelColumnReflect2UserClientColumn(strings.TrimSpace(th.String()), fieldsInfo); err == nil {
			excelCol2structCol[k] = tmp
		}
	}
	// 获取相同公司下的手机号码列表
	mps := models.GetMobilePhoneList(c.curUser.UserCompanyId)
	// fixed value
	uidStr := strconv.Itoa(c.curUser.Id)
	uComIdStr := strconv.Itoa(c.curUser.UserCompanyId)
	nowStr := fmt.Sprintf(`'%s'`, time.Now().Format(enums.BaseTimeFormat))
	// sql prepare
	// insert
	sqlTmpl := `insert into ` + models.UserClientTBName() + ` (%s) values (%s)`
	SegStrs := make([]string, 0)
	ValStrs := make([]string, 0)
	// update
	updateTmpl := `update ` + models.UserClientTBName() + ` set %s where %s`
	setTmpl := `%s=%s`
	ColStrs := make([]string, 0)
	WhereStrs := make([]string, 0)
	// 消息提示
	msg := make(map[string]string)
	// 电话号码重复
	dupRows := make([]string, 0)
	// error rows
	errRows := make([]string, 0)
	for m, row := range sheet.Rows {
		if m == 0 {
			continue
		}
		SegStr := []string{}
		ValStr := []string{}
		errOccur := false
		for k, cell := range row.Cells {
			text := strings.TrimSpace(cell.String())
			if len(text) > 0 {
				if strings.Compare(excelCol2structCol[k].ColumnName, "mobile_phone") == 0 {
					// mobile_phone在同一个backend_user_id下相同的，做更新。
					for _, mp := range mps {
						if strings.Compare(mp.MobilePhone, text) == 0 {
							// 电话号码重复的，直接跳过
							errOccur = true
							dupRows = append(dupRows, strconv.Itoa(m+1))
							break
							// update
							idsCol := []string{"id"}
							idsVal := []string{strconv.Itoa(mp.Id)}
							// id放在第一个位置
							SegStr = append(idsCol, SegStr...)
							ValStr = append(idsVal, ValStr...)
						}
					}
				} else if strings.Compare(excelCol2structCol[k].ColumnName, "state") == 0 || strings.Compare(excelCol2structCol[k].ColumnName, "clue_from") == 0 {
					states := strings.Split(excelCol2structCol[k].FieldTypeValue, ",")
					var stateInit int
					for k, v := range states {
						if v == text {
							stateInit = k + 1
							break
						}
					}
					text = strconv.Itoa(stateInit)
				}
				SegStr = append(SegStr, excelCol2structCol[k].ColumnName)
				ValStr = append(ValStr, fmt.Sprintf(`'%s'`, text))

			} else {
				if excelCol2structCol[k].Required == 1 {
					errOccur = true
					errRows = append(errRows, strconv.Itoa(m+1))
					break
				}
			}
		}
		if errOccur {

		} else if SegStr[0] == "id" {
			// 更新
			// 补充updated字段
			SegStr = append(SegStr, "updated")
			ValStr = append(ValStr, nowStr)

			tmp := []string{}
			colLen := len(SegStr)
			for i := 1; i < colLen; i++ {
				tmp = append(tmp, fmt.Sprintf(setTmpl, SegStr[i], ValStr[i]))
			}
			ColStrs = append(ColStrs, strings.Join(tmp, ","))
			WhereStrs = append(WhereStrs, fmt.Sprintf(setTmpl, SegStr[0], ValStr[0]))
		} else {
			// 插入
			// 补充应添加的字段
			SegStr = append(SegStr, "backend_user_id", "belong_backend_user_id", "user_company_id", "created", "updated")
			ValStr = append(ValStr, uidStr, uidStr, uComIdStr, nowStr, nowStr)
			SegStrs = append(SegStrs, strings.Join(SegStr, ","))
			ValStrs = append(ValStrs, strings.Join(ValStr, ","))
		}

	}
	rowLen := len(SegStrs)
	sqls := make([]string, 0)
	for i := 0; i < rowLen; i++ {
		sql := fmt.Sprintf(sqlTmpl, SegStrs[i], ValStrs[i])
		sqls = append(sqls, sql)
	}
	rowLen = len(ColStrs)
	for i := 0; i < rowLen; i++ {
		sql := fmt.Sprintf(updateTmpl, ColStrs[i], WhereStrs[i])
		sqls = append(sqls, sql)
	}
	utils.Info("sql: %v\n", strings.Join(sqls, ";"))

	// 合并返回的提示消息
	msg["err"] = strings.Join(errRows, ",")
	msg["dup"] = strings.Join(dupRows, ",")
	// 导入数据库
	o := orm.NewOrm()
	if _, err := o.Raw(strings.Join(sqls, ";")).Exec(); err == nil {
		c.jsonResult(enums.JRCodeSucc, "导入成功", msg)
	} else {
		utils.LogError("导入客户资料，插入数据库失败。" + err.Error())
		c.jsonResult(enums.JRCodeFailed, "导入失败", "")
	}

	// c.jsonResult(enums.JRCodeFailed, "testing...", "")

	// filePath := "static/upload/" + h.Filename
	// // 保存位置在 static/upload, 没有文件夹要先创建
	// if err := c.SaveToFile("custtemplet", filePath); err == nil {
	// 	c.jsonResult(enums.JRCodeSucc, "上传成功", "/"+filePath)
	// } else {
	// 	utils.LogError("导入客户资料失败。" + err.Error())
	// 	c.jsonResult(enums.JRCodeFailed, "上传失败", "")
	// }

}

// 需要被导入的单元格
func excelColumnReflect2UserClientColumn(cellVal string, fieldsInfo []*models.UserClientField) (*models.UserClientField, error) {
	// state, clue_from是int类型的，会标记字段的值，去除它
	// 客户状态 (无意向，中等意向，有意向)
	if strings.Index(cellVal, "(") != -1 {
		arr := strings.Split(cellVal, "(")
		cellVal = strings.TrimSpace(arr[0])
	}
	for _, item := range fieldsInfo {
		if strings.Compare(item.FieldName, cellVal) == 0 {
			return item, nil
		}
	}
	return &models.UserClientField{}, errors.New("invalid column")
}

// 下载客户信息
func (c *UserClientController) Download() {
	//直接反序化获取json格式的requestbody里的值
	var params models.UserClientQueryParam
	var err error
	//获取form里的值
	if err = c.ParseForm(&params); err != nil {
		c.pageError("数据无效，请刷新后重试")
	}
	// 根据uid获取记录
	if !c.curUser.IsSuper {
		params.CompanyId = c.curUser.UserCompanyId
		if len(params.BackendUserIds) > 0 {
			if params.BackendUserIds[0] == -1 {
				params.BackendUserIds = c.curUser.UIds
			}
		} else {
			params.BackendUserIds = []int{c.curUser.Id}
		}
	}
	data := models.UserClientDataList(&params)

	// 获取自定义的字段信息
	var paramsAttr models.UserClientFieldQueryParam
	paramsAttr.UserCompanyId = c.curUser.UserCompanyId
	attrs := models.UserClientFieldDataList(&paramsAttr)
	// 获取归属人名字
	var paramsUser = models.BackendUserQueryParam{}
	// 当前登录用户的同个公司的list
	if !c.curUser.IsSuper {
		paramsUser.UIds = c.curUser.UIds
	}
	//获取数据列表和总数
	dataUsers := models.BackendUserDataList(&paramsUser)

	// file
	var file *xlsx.File
	var sheet *xlsx.Sheet
	// new
	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		utils.LogError("新建excel文件失败。" + err.Error())
	}

	segNeedExport := make(map[int]*models.UserClientField)

	for _, attr := range attrs {
		switch attr.ColumnName {
		case "name":
			segNeedExport[1] = attr
		case "mobile_phone":
			segNeedExport[2] = attr
		case "contact_phone":
			segNeedExport[3] = attr
		case "backend_user_id":
			segNeedExport[4] = attr
		case "created":
			segNeedExport[5] = attr
		case "comment":
			segNeedExport[6] = attr
		case "address":
			segNeedExport[7] = attr
		case "state":
			segNeedExport[8] = attr
		case "feature":
			segNeedExport[9] = attr
		case "complaint":
			segNeedExport[10] = attr
		case "latest_communicated":
			segNeedExport[11] = attr
		case "clue_from":
			segNeedExport[12] = attr
		case "email":
			segNeedExport[13] = attr
		case "belong_backend_user_id":
			segNeedExport[14] = attr
		case "user_company_id":
			segNeedExport[15] = attr
		case "updated":
			segNeedExport[16] = attr
		case "column1":
			segNeedExport[17] = attr
		case "column2":
			segNeedExport[18] = attr
		case "column3":
			segNeedExport[19] = attr
		case "column4":
			segNeedExport[20] = attr
		case "column5":
			segNeedExport[21] = attr
		case "column6":
			segNeedExport[22] = attr
		case "column7":
			segNeedExport[23] = attr
		case "column8":
			segNeedExport[24] = attr
		case "column9":
			segNeedExport[25] = attr
		case "column10":
			segNeedExport[26] = attr
		case "column11":
			segNeedExport[27] = attr
		case "column12":
			segNeedExport[28] = attr
		case "column13":
			segNeedExport[29] = attr
		case "column14":
			segNeedExport[30] = attr
		case "column15":
			segNeedExport[31] = attr
		case "column16":
			segNeedExport[32] = attr
		case "dial_state":
			segNeedExport[33] = attr
		}
	}
	// 表头
	th := sheet.AddRow()
	if e, ok := segNeedExport[1]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[2]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[3]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[4]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName

	}
	if e, ok := segNeedExport[5]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[6]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[7]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[8]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	// if segNeedExport[9].ListShow {
	// 	td := th.AddCell()
	// 	td.Value = item.Feature
	// }
	if e, ok := segNeedExport[10]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[11]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[12]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[13]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[14]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	// if segNeedExport[15].ListShow {
	// 	td := th.AddCell()
	// 	td.Value = item.UserCompanyId
	// }
	if e, ok := segNeedExport[16]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[17]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[18]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[19]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[20]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[21]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[22]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[23]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[24]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[25]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[26]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[27]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[28]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[29]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[30]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[31]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[32]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	if e, ok := segNeedExport[33]; ok && e.ListShow {
		td := th.AddCell()
		td.Value = e.FieldName
	}
	// 添加数据
	for _, item := range data {
		tr := sheet.AddRow()
		if e, ok := segNeedExport[1]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Name
		}
		if e, ok := segNeedExport[2]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.MobilePhone
		}
		if e, ok := segNeedExport[3]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.ContactPhone
		}
		if e, ok := segNeedExport[4]; ok && e.ListShow {
			td := tr.AddCell()
			for _, user := range dataUsers {
				if user.Id == item.BackendUserId {
					td.Value = user.RealName
				}
			}

		}
		if e, ok := segNeedExport[5]; ok && e.ListShow {
			td := tr.AddCell()
			arr := strings.Split(item.Created, " ")
			if len(arr) > 1 {
				td.Value = strings.Join([]string{arr[0], arr[1]}, " ")
			} else {
				td.Value = ""
			}

		}
		if e, ok := segNeedExport[6]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Comment
		}
		if e, ok := segNeedExport[7]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Address
		}
		if e, ok := segNeedExport[8]; ok && e.ListShow {
			td := tr.AddCell()
			tmp := strings.Split(segNeedExport[8].FieldTypeValue, ",")
			for k, v := range tmp {
				if (k + 1) == item.State {
					td.Value = v
					break
				}
			}

		}
		// if segNeedExport[9].ListShow {
		// 	td := tr.AddCell()
		// 	td.Value = item.Feature
		// }
		if e, ok := segNeedExport[10]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Complaint
		}
		if e, ok := segNeedExport[11]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.LatestCommunicated
			arr := strings.Split(item.LatestCommunicated, " ")
			if strings.Compare(arr[0], "2000-01-01") == 0 {
				td.Value = ""
			} else if len(arr) > 1 {
				td.Value = strings.Join([]string{arr[0], arr[1]}, " ")
			} else {
				td.Value = ""
			}
		}
		if e, ok := segNeedExport[12]; ok && e.ListShow {
			td := tr.AddCell()
			tmp := strings.Split(segNeedExport[12].FieldTypeValue, ",")
			for k, v := range tmp {
				if (k + 1) == item.ClueFrom {
					td.Value = v
					break
				}
			}
		}
		if e, ok := segNeedExport[13]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Email
		}
		if e, ok := segNeedExport[14]; ok && e.ListShow {
			td := tr.AddCell()
			for _, user := range dataUsers {
				if user.Id == item.BelongBackendUserId {
					td.Value = user.RealName
				}
			}
		}
		// if segNeedExport[15].ListShow {
		// 	td := tr.AddCell()
		// 	td.Value = item.UserCompanyId
		// }
		if e, ok := segNeedExport[16]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Updated
		}
		if e, ok := segNeedExport[17]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column1
		}
		if e, ok := segNeedExport[18]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column2
		}
		if e, ok := segNeedExport[19]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column3
		}
		if e, ok := segNeedExport[20]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column4
		}
		if e, ok := segNeedExport[21]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column5
		}
		if e, ok := segNeedExport[22]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column6
		}
		if e, ok := segNeedExport[23]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column7
		}
		if e, ok := segNeedExport[24]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column8
		}
		if e, ok := segNeedExport[25]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column9
		}
		if e, ok := segNeedExport[26]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column10
		}
		if e, ok := segNeedExport[27]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column11
		}
		if e, ok := segNeedExport[28]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column12
		}
		if e, ok := segNeedExport[29]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column13
		}
		if e, ok := segNeedExport[30]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column14
		}
		if e, ok := segNeedExport[31]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column15
		}
		if e, ok := segNeedExport[32]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.Column16
		}
		if e, ok := segNeedExport[33]; ok && e.ListShow {
			td := tr.AddCell()
			td.Value = item.DialState
		}
	}

	filePath := "./static/download/客户资料.xlsx"
	if err := os.MkdirAll("static/download", 0666); err != nil {
		utils.LogError("创建下载文件夹失败。" + err.Error())
		// c.jsonResult(enums.JRCodeFailed, "创建下载文件夹失败。", 0)
	}
	if err := file.Save(filePath); err == nil {
		c.Ctx.Output.Download(filePath, "客户资料.xlsx")
	} else {
		utils.LogError("下载客户资料模板失败。" + err.Error())
		// c.jsonResult(enums.JRCodeFailed, "下载客户资料模板失败。", 0)
	}
	c.StopRun()
}

// 通过号码获得一条客户信息
func (c *UserClientController) IsExistedCalledPhone() {
	phone := c.GetString("tocallphone")

	if len(phone) == 0 {
		c.jsonResult(enums.JRCodeFailed, "获取号码失败，请刷新重试。", "")
	}

	var err error
	var m models.UserClient
	query := orm.NewOrm().QueryTable(models.UserClientTBName())
	err = query.Filter("mobile_phone", phone).Filter("user_company_id", c.curUser.UserCompanyId).One(&m)

	if err == orm.ErrMultiRows {
		utils.Info("%v\n", m)
		c.jsonResult(enums.JRCodeSucc, "有多行", m)
	}
	if err == orm.ErrNoRows {
		c.jsonResult(enums.JRCodeSucc, "不存在此号码", m)
	} else {
		c.jsonResult(enums.JRCodeSucc, "有一行", m)
	}
}
