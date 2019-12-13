package controllers

import (
    "encoding/json"

    "zq/callout_crm/enums"
    "zq/callout_crm/models"
    "zq/callout_crm/sms"

    "fmt"
    "strconv"
    "strings"
    "time"

    "github.com/astaxie/beego/orm"
)

//SmsRecordController 角色管理
type SmsRecordController struct {
    BaseController
}

//Prepare 参考beego官方文档说明
func (c *SmsRecordController) Prepare() {
    //先执行
    c.BaseController.Prepare()
    //如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
    c.checkAuthor("DataGrid", "DataList", "SendMsg")
    //如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
    //权限控制里会进行登录验证，因此这里不用再作登录验证
    //c.checkLogin()
}

// DataGrid 角色管理首页 表格获取数据
func (c *SmsRecordController) DataGrid() {
    //直接反序化获取json格式的requestbody里的值
    var params models.SmsRecordQueryParam
    json.Unmarshal(c.Ctx.Input.RequestBody, &params)
    //获取数据列表和总数
    if !c.curUser.IsSuper {
        params.UIds = c.curUser.UIds
    }
    data, total := models.SmsRecordPageList(&params)
    //定义返回的数据结构
    result := make(map[string]interface{})
    result["total"] = total
    result["rows"] = data
    c.Data["json"] = result
    c.ServeJSON()
}

// DataList 角色列表
func (c *SmsRecordController) DataList() {
    var params = models.SmsRecordQueryParam{}
    //获取数据列表和总数
    data := models.SmsRecordDataList(&params)
    //定义返回的数据结构
    c.jsonResult(enums.JRCodeSucc, "", data)
}

// Edit 添加、编辑角色界面
func (c *SmsRecordController) Edit() {
    if c.Ctx.Request.Method == "POST" {
        c.Save()
    }
    Id, _ := c.GetInt(":id", 0)
    m := models.SmsRecord{Id: Id}
    if Id > 0 {
        o := orm.NewOrm()
        err := o.Read(&m)
        if err != nil {
            c.pageError("数据无效，请刷新后重试")
        }
    }
    c.Data["m"] = m
    c.setTpl("smsrecord/edit.html", "shared/layout_pullbox.html")
    c.LayoutSections = make(map[string]string)
    c.LayoutSections["footerjs"] = "smsrecord/edit_footerjs.html"
}

//Save 添加、编辑页面 保存
func (c *SmsRecordController) Save() {
    var err error
    m := models.SmsRecord{}
    //获取form里的值
    if err = c.ParseForm(&m); err != nil {
        c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
    }
    o := orm.NewOrm()
    if m.Id == 0 {
        // insert添加created时间
        m.SendTime = time.Now().Format(enums.BaseTimeFormat)
        if _, err = o.Insert(&m); err == nil {
            c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
        } else {
            c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
        }

    } else {
        if _, err = o.Update(&m, "Mobile", "Content", "Result"); err == nil {
            c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
        } else {
            c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
        }
    }

}

//Delete 批量删除
func (c *SmsRecordController) Delete() {
    strs := c.GetString("ids")
    ids := make([]int, 0, len(strs))
    for _, str := range strings.Split(strs, ",") {
        if id, err := strconv.Atoi(str); err == nil {
            ids = append(ids, id)
        }
    }
    if num, err := models.SmsRecordBatchDelete(ids); err == nil {
        c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
    } else {
        c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
    }
}

// 发送短信
func (c *SmsRecordController) SendMsg() {
    ucid, _ := c.GetInt("ucid", 0)
    classify, _ := c.GetInt("classify", 0)
    content := c.GetString("content")
    phone := c.GetString("phone")
    // check received params
    if len(content) == 0 || len(phone) == 0 {
        c.jsonResult(enums.JRCodeFailed, "发送失败，发送内容不能为空！", 0)
    }
    // query db, get needed params
    params := map[string]interface{}{
        "user_company_id": ucid,
        "classify": classify,
    }
    m := models.SmsSetOneByParams(params)
    if len(m.Account) == 0 {
        c.jsonResult(enums.JRCodeFailed, "获取发送参数失败，请联系管理员！", 0)
    }

    sms := &sms.SendMsgParam{
        Account  : m.Account,
        Password : m.Password,
        Signature: m.Signature,
        Msg      : content,
        Report   : "true",
        Phone    : phone,
    }
    if respBytes, err := sms.SendMsg(); err != nil {
        // 与创蓝平台无关的错误处理
        c.jsonResult(enums.JRCodeFailed, string(respBytes), 0)
    } else {
        // 与创蓝平台相关
        res := make(map[string]string)
        if err := json.Unmarshal(respBytes, &res); err == nil {
            if res["code"] == "0" {
                // 成功
                o := orm.NewOrm()
                rcd := models.SmsRecord{
                    Mobile: phone,
                    Content: content,
                    Result: 1,
                    SendTime: time.Now().Format(enums.BaseTimeFormat),
                    Classify: 0,
                    BackendUserId: c.curUser.Id,
                    UserCompanyId: c.curUser.UserCompanyId,
                }
                if _, err = o.Insert(&rcd); err == nil {
                    c.jsonResult(enums.JRCodeSucc, "发送成功，添加发送记录成功", m.Id)
                } else {
                    c.jsonResult(enums.JRCodeFailed, "发送成功，添加发送记录失败", m.Id)
                }
            } else {
                // 失败
                c.jsonResult(enums.JRCodeFailed, "发送失败，" + res["errorMsg"], 0)
            }
        } else {
            c.jsonResult(enums.JRCodeFailed, "解析返回数据失败，请联系管理员！", 0)
        }
    }
}