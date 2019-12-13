package sms

import (
    "bytes"
    "encoding/json"
    // "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "time"

    "zq/callout_crm/utils"
)

const (
    send_url = `https://smssh1.253.com/msg/send/json`   // 五参数发送接口
)

type SendMsgParam struct {
    Account   string
    Password  string
    Signature string
    Msg       string
    Report    string
    Phone     string // phones是逗号分隔的
}

func (self *SendMsgParam) SendMsg() ([]byte, error) {
    //请登录zz.253.com获取API账号、密码以及短信发送的URL
    params := make(map[string]interface{})

    params["account"] = self.Account   //创蓝API账号
    params["password"] = self.Password //创蓝API密码
    //设置您要发送的内容：其中“【】”中括号为运营商签名符号，多签名内容前置添加提交
    params["msg"] = url.QueryEscape(fmt.Sprintf("【%s】%s", self.Signature, self.Msg))
    params["report"] = self.Report     // 一般为"true"
    params["phone"] = self.Phone

    // fmt.Printf("%#v\n", params)
    // // 模拟请求返回
    // return simulateMsgReturn(false)
    
    bytesData, err := json.Marshal(params)
    if err != nil {
        utils.LogError(err.Error())
        return bytesData, err
    }

    reader := bytes.NewReader(bytesData)
    request, err := http.NewRequest("POST", send_url, reader)

    if err != nil {
        utils.LogError(err.Error())
        return []byte("新建请求失败"), err
    }
    request.Header.Set("Content-Type", "application/json;charset=UTF-8")
    client := http.Client{}
    resp, err := client.Do(request)
    if err != nil {
        utils.LogError(err.Error())
        return []byte("返回数据失败"), err
    }
    respBytes, err := ioutil.ReadAll(resp.Body)

    return respBytes, err
}

func simulateMsgReturn(success bool) ([]byte, error) {
    m := make(map[string]string)
    if success {
        m["code"] = "0"
        m["msgId"] = "18052415065227118"
        m["time"] = time.Now().Format("20060102150405000")
        m["errorMsg"] = ""
    } else {
        m["code"] = "102"
        m["msgId"] = ""
        m["time"] = time.Now().Format("20060102150405000")
        m["errorMsg"] = "密码错误"
    }
    b, err := json.Marshal(m)
    return b, err
}