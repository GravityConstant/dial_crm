package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"os"
	"os/exec"
	"path/filepath"

	"zq/callout_crm/models"
	"zq/callout_crm/utils"

	"github.com/tealeg/xlsx"
)

//CallPgCdrController 话单管理
type CallPgCdrController struct {
	BaseController
}

//Prepare 参考beego官方文档说明
func (c *CallPgCdrController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid", "DownloadCdr", "DownloadRecord")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

//Index 角色管理首页
func (c *CallPgCdrController) Index() {
	//是否显示更多查询条件的按钮
	c.Data["showMoreQuery"] = true
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "callpgcdr/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "callpgcdr/index_footerjs.html"
}

// DataGrid 角色管理首页 表格获取数据
func (c *CallPgCdrController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.CallPgCdrQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	fmt.Println("CallPgCdrController, DataGrid================")
	fmt.Printf("%#v\n", params)

	if !c.curUser.IsSuper {
		if len(params.DialplanIds) == 0 {
			params.DialplanIds = c.curUser.UIds
		}
	}

	//获取数据列表和总数
	data, total := models.CallPgCdrPageList(&params)

	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

// 现在是按时间下载的
func (c *CallPgCdrController) DownloadCdr() {
	//直接反序化获取json格式的requestbody里的值
	var params models.CallPgCdrQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	if !c.curUser.IsSuper {
		if len(params.DialplanIds) == 0 {
			params.DialplanIds = c.curUser.UIds
		}
	}
	data := models.CallPgCdrDataList(&params)

	f := xlsx.NewFile()
	sheet, _ := f.AddSheet("cdr1")
	// 头行
	th := sheet.AddRow()
	thCon := []string{"通话时间", "主叫号码", "被叫号码", "通话时间(秒)", "通话状态", "呼叫类型"}
	for _, v := range thCon {
		td := th.AddCell()
		td.Value = v
	}

	callerFmt := `%s(%s)`
	
	for _, item := range data {
		tr := sheet.AddRow()

		td1 := tr.AddCell()
		arr := strings.Split(item.StartStamp, " ")
		if len(arr) > 1 {
			td1.Value = strings.Join([]string{arr[0], arr[1]}, " ")
		}
		

		td2 := tr.AddCell()
		td2.Value = fmt.Sprintf(callerFmt, item.CallerIdName, item.CallerIdNumber)

		td3 := tr.AddCell()
		td3.Value = item.DestinationNumber

		td4 := tr.AddCell()
		tf := time.Unix(int64(item.Billsec), 0)
		tf = tf.Add(-8 * time.Hour)
		td4.Value = tf.Format("15:04:05")

		td5 := tr.AddCell()
		if item.Billsec > 0 {
			td5.Value = "已接听"
		} else {
			td5.Value = "未接听"
		}
		
		td6 := tr.AddCell()
		if strings.Compare(item.Direction, "outbound") == 0 {
			td6.Value = "呼出"
		} else {
			td6.Value = "呼入"
		}
		
	}

	filePath := "./static/download/cdr.xlsx"
	if err := os.MkdirAll("static/download", 0666); err != nil {
		utils.LogError("创建下载文件夹失败。" + err.Error())
		// c.jsonResult(enums.JRCodeFailed, "创建下载文件夹失败。", 0)
	}
	if err := f.Save(filePath); err == nil {
		c.Ctx.Output.Download(filePath, "cdr.xlsx")
	} else {
		utils.LogError("下载客户资料模板失败。" + err.Error())
		// c.jsonResult(enums.JRCodeFailed, "下载客户资料模板失败。", 0)
	}
	c.StopRun()
}

// 下载录音
func (c *CallPgCdrController) DownloadRecord() {
	// 直接反序化获取json格式的requestbody里的值
	var params models.CallPgCdrQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	st := c.GetString("StartStamp")
	et := c.GetString("EndStamp")

	params.StartStamp = st
	params.EndStamp = et

	// 获取数据列表和总数
	if !c.curUser.IsSuper {
		if len(params.DialplanIds) == 0 {
			params.DialplanIds = c.curUser.UIds
		}
	}
	data := models.CallPgCdrDataList(&params)

	// 2019/9/25/80000/80000-63368908-20190925155344000.wav
	datetime := make([]string, 0)

	for _, item := range data {
		rs := strings.Split(item.RecordFile, `/`)
		if (len(rs) > 2) {
			recStr := strings.Join([]string{rs[0], rs[1], rs[2]}, `/`)
			datetime = addUniqueArray(datetime, recStr)
		}

	}

	filenames := make([]string, 0)
	if len(datetime) > 0 {
		if agent, err := models.AgentOneByUid(c.curUser.Id); err == nil {
			for _, val := range datetime {
				tmp := filepath.Join(val, agent.ExtNo)
				filenames = append(filenames, tmp)
			}
		}
	}

	
	files := make([]string, 0)
	if len(filenames) > 0 {
		// 验证是否文件是否存在
		for _, val := range filenames {
			tmp := `/home/voices/records/` + val
			if _, err := os.Stat(tmp); err == nil {
				files = append(files, tmp)
			}
		}
		utils.Info("%#v", files)
		// 打包
		if len(files) > 0 {
			nowFileName := time.Now().Format("20060102150405") + `.tar.gz`
			output := `./static/download/records.tar.gz`
			if err := os.MkdirAll("static/download", 0666); err != nil {
				utils.LogError("创建下载文件夹失败。" + err.Error())
			} else {
				args := []string{"-zcf", output}
				args = append(args, files...)
				cmd := exec.Command("tar", args...)
				if err := cmd.Run(); err != nil {
					utils.LogError("执行tar命令失败。" + err.Error())
				}
				utils.Info("%v", nowFileName)
				c.Ctx.Output.Download(output, nowFileName)
			}
		}
		
	}

	c.StopRun()
}


func addUniqueArray(datetime []string, s string) []string {
	for _, val := range datetime {
		if val == s {
			return datetime
		}
	}
	return append(datetime, s)
}