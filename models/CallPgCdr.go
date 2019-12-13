package models

import (
	"fmt"
	"strconv"
    "strings"

	"github.com/astaxie/beego/orm"
	"zq/callout_crm/utils"
)

// TableName 设置表名
func (a *CallPgCdr) TableName() string {
	return CallPgCdrTBName()
}

// CallPgCdrQueryParam 用于搜索的类
type CallPgCdrQueryParam struct {
	BaseQueryParam
	CallerIdNumber    string // 主叫
	DestinationNumber string // 被叫
	StartStamp        string // 开始时间
	EndStamp          string // 结束时间
	BillsecStatus     string // 接通状态
	DialplanIds       []int
	DialplanNumber    string
	Direction         string
	BindPhones        []interface{}
}

// CallPgCdr 实体类
type CallPgCdr struct {
	Id                     int64
	LocalIpV4              string
	CallerIdName           string
	CallerIdNumber         string
	OutboundCallerIdNumber string
	DestinationNumber      string
	Context                string // 对应freeswitch中的conf里的context
	StartStamp             string
	AnswerStamp            string
	EndStamp               string
	Duration               int // 开始打到挂的时间
	Billsec                int // 接听的时间
	HangupCause            string
	Uuid                   string
	BlegUuid               string
	Accountcode            string
	ReadCodec              string
	WriteCodec             string
	RecordFile             string
	Direction              string
	SipHangupDisposition   string
	OriginationUuid        string
	SipGatewayName         string
	DialplanId             int64
	FeeRate                float64
}

// CallPgCdrPageList 获取分页数据
func CallPgCdrPageList(params *CallPgCdrQueryParam) ([]*CallPgCdr, int64) {
	o := orm.NewOrm()
	query := new(utils.QueryString)
	data := make([]*CallPgCdr, 0)
	// 自定义参数
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	case "Seq":
		sortorder = "Seq"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	// query.Filter("", "direction", "outbound")
	if len(params.DialplanIds) > 0 {
		var interfaceSlice []interface{} = make([]interface{}, len(params.DialplanIds))
		for i, d := range params.DialplanIds {
			   interfaceSlice[i] = d
		}
		query.Filter("", "dialplan_id__in", interfaceSlice...)
	} else {
		query.Filter("", "dialplan_id__gt", 0)
	}
	if len(params.DialplanNumber) > 0 {
		query.Filter("", "destination_number", params.DialplanNumber)
	}
	if len(params.Direction) > 0 {
		query.Filter("", "direction", params.Direction)
	}
	if len(params.CallerIdNumber) > 0 {
		query.Filter("", "caller_id_name__icontains", params.CallerIdNumber)
	} else {
		// 过滤数据
		// caller_id_number不能是0000000000，不能是1000等分机号形式的
		// query.Filter("", "caller_id_number__regex", "^1[0-9]{10}|[2-9]{6,7}|0{11}$")
	}
	if len(params.DestinationNumber) > 0 {
		query.Filter("", "destination_number__icontains", params.DestinationNumber)
	} else {
		// 过滤数据
		// query.Filter("", "destination_number__regex", "^1[0-9]{10}|[2-9]{6,7}|0{11}$")
		// if len(params.BindPhones) > 0 {
		// 	query.Filter("", "destination_number__in", params.BindPhones...)
		// }
	}
	if len(params.StartStamp) > 0 {
		query.Filter("", "start_stamp__gte", params.StartStamp)
	}
	if len(params.EndStamp) > 0 {
		query.Filter("", "start_stamp__lte", params.EndStamp)
	}
	if len(params.BillsecStatus) > 0 {
		if bs, err := strconv.Atoi(params.BillsecStatus); err == nil && bs > 0 {
			query.Filter("", "billsec__gt", 0)
		} else if err == nil && bs == 0 {
			query.Filter("", "billsec", 0)
		}
	}
	query.OrderBy("", sortorder)
	query.Limit(params.Limit, params.Offset)

	var sql string
	var total int64
	sqlCol := `select id, local_ip_v4, (select real_name from crm_backend_user where id=dialplan_id) as caller_id_name, (select ext_no from crm_agent where backend_user_id=dialplan_id) as caller_id_number, outbound_caller_id_number, destination_number, context, start_stamp, answer_stamp, end_stamp, duration, billsec, hangup_cause, uuid, bleg_uuid, accountcode, read_codec, write_codec, record_file, direction, sip_hangup_disposition, origination_uuid, sip_gateway_name, dialplan_id, fee_rate `
	sqlFrom := `from ` + CallPgCdrTBName() + ` where char_length(sip_gateway_name) > 0 `
	if len(query.FilterStr) > 0 {
		sql = sqlCol + sqlFrom + "and " + query.String()
	} else {
		sql = sqlCol + sqlFrom + query.String()
	}
	if num, err := o.Raw(sql).QueryRows(&data); err == nil && num > 0 {
		sqlCol = `select count(*) total `
		sql = sqlCol + sqlFrom
		if len(query.FilterStr) > 0 {
			sql = sqlCol + sqlFrom + "and " + query.FilterStr
		}
		type Option struct {
			Total int64
		}
		var option Option
		if err := o.Raw(sql).QueryRow(&option); err == nil {
			total = option.Total
		}
	}

	return data, total
}

// CallPgCdrOne 获取单条
func CallPgCdrOne(id int) (*CallPgCdr, error) {
	o := orm.NewOrm()
	m := CallPgCdr{Id: int64(id)}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// CallPgCdrDataList 获取列表
func CallPgCdrDataList(params *CallPgCdrQueryParam) []*CallPgCdr {
	params.Limit = -1
	params.Sort = "Id"
	params.Order = "asc"
	data, _ := CallPgCdrPageList(params)
	return data
}

type CalledCount struct {
	Called string
	Times  int
}

// 获取某个号码某个时间段被呼叫的次数
func CalledCountByTime(params *CallPgCdrQueryParam) *CalledCount {
	o := orm.NewOrm()
	m := &CalledCount{}
    var tbName string

    if len(params.StartStamp) == 0 {
        return m
    } else {
        ss := strings.Split(params.StartStamp, "-")
        tbName = "call_pg_cdr_" + ss[0] + "_" + ss[1]
    }

    // 主叫主动取消呼叫的不算一次 and sip_hangup_disposition!='send_cancel'
	sql := "SELECT destination_number as called, count(destination_number) as times " + 
		   "from (" + 
	           "SELECT destination_number from %s where start_stamp > '%s' " + 
           ") foo GROUP BY destination_number HAVING destination_number in ('%s', '0%s')" 
    sql = fmt.Sprintf(sql, tbName, params.StartStamp, params.DestinationNumber, params.DestinationNumber)

    o.Raw(sql).QueryRow(m)

    return m
}

type CallerCount struct {
	Caller string
	Times  int
}

// 获取坐席某个时间段外呼的次数
func CallerCountByTime(params *CallPgCdrQueryParam) []*CallerCount {
	o := orm.NewOrm()
	data := make([]*CallerCount, 0)
    var tbName string

    if len(params.StartStamp) == 0 {
        return data
    } else {
        ss := strings.Split(params.StartStamp, "-")
        tbName = "call_pg_cdr_" + ss[0] + "_" + ss[1]
    }

    // 主叫主动取消呼叫的不算一次 and sip_hangup_disposition!='send_cancel'
	sql := "SELECT caller_id_number as caller, count(caller_id_number) as times " + 
		   "from (" + 
	           "SELECT caller_id_number from %s where start_stamp > '%s' " + 
	           "and caller_id_number in (" +
	           		"SELECT gp.phone from %s gp INNER JOIN %s a " +
	           		"on a.id=gp.agent_id where a.ext_no='%s'" +
	           ")" +
           ") foo GROUP BY caller_id_number" 
    sql = fmt.Sprintf(sql, tbName, params.StartStamp, GatewayPhoneTBName(), AgentTBName(),
    				 params.CallerIdNumber)

    o.Raw(sql).QueryRows(&data)

    return data
}