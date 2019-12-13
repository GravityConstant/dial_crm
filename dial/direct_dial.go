package dial

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"zq/callout_crm/utils"

	"github.com/vma/esl"
)

type Handler struct {
	ExtensionNumber string
	DialplanNumber  string
	Gateway         string
	VirtualPhone    string
	CalloutPhone    string
	Err             error
	ALegUId         string
	BLegUId         string
	ALegState       map[string]esl.EventName
	BLegState       map[string]esl.EventName
	Done            chan bool
	UserId          int
}

func Init(handler *Handler) {
	con, err := esl.NewConnection("127.0.0.1:8021", handler)
	if err != nil {
		log.Fatal("ERR connecting to freeswitch:", err)
	}
	con.HandleEvents()
}

func (h *Handler) OnConnect(con *esl.Connection) {
	con.SendRecv("event", "plain", "ALL")
	// aleg uuid
	if h.ALegUId, h.Err = con.Api("create_uuid"); h.Err != nil {
		utils.Error("create_uuid error: %s", h.Err)
		con.Close()
	}
	// bleg uuid
	if h.BLegUId, h.Err = con.Api("create_uuid"); h.Err != nil {
		utils.Error("create_uuid error: %s", h.Err)
		con.Close()
	}
	utils.Info("aleg: %s, bleg: %s", h.ALegUId, h.BLegUId)

	// 主叫
	outbound := `{origination_uuid=%s,dialplan_id=%d,hangup_after_bridge=true}%s`
	outbound = fmt.Sprintf(outbound, h.ALegUId, h.UserId, h.ExtensionNumber)
	// 被叫
	// {origination_caller_id_number=28324295,sip_h_Diversion=<sip:28324295@ip>}sofia/gateway/zqzj/$1
	calloutString := `&bridge({origination_uuid=%s,dialplan_id=%d,origination_caller_id_number=%s,sip_h_Diversion=<sip:%s@ip>}sofia/gateway/%s/%s)`
	calloutString = fmt.Sprintf(calloutString, h.BLegUId, h.UserId, h.VirtualPhone, h.DialplanNumber, h.Gateway, h.CalloutPhone)
	utils.Info("callout string: %s, %s", outbound, calloutString)

	var bgJobId string
	if bgJobId, h.Err = con.BgApi("originate", outbound, calloutString); h.Err != nil {
		log.Println(h.Err)
	} else {
		log.Printf("bg job id: %s", bgJobId)
	}

}

// func (h *Handler) OnConnect(con *esl.Connection) {
// 	con.SendRecv("event", "plain", "ALL")

// 	// 呼叫座席
// 	if h.ALegUId, h.Err = con.Api("create_uuid"); h.Err != nil {
// 		utils.Error("create_uuid error: %s", h.Err)
// 		con.Close()
// 	}
// 	alegStr := `{origination_uuid=%s}%s`
// 	alegStr = fmt.Sprintf(alegStr, h.ALegUId, h.ExtensionNumber)
// 	if _, h.Err = con.BgApi("originate", alegStr, "&park()"); h.Err != nil {
// 		utils.Error("originate aleg error: %s", h.Err)
// 		con.Close()
// 	}

// 	// 外呼号码
// 	if h.BLegUId, h.Err = con.Api("create_uuid"); h.Err != nil {
// 		utils.Error("create_uuid error: %s", h.Err)
// 		con.Close()
// 	}
// 	blegStr := `{origination_uuid=%s,origination_caller_id_number=%s,sip_h_Diversion=<sip:%s@ip>}sofia/gateway/%s/%s`
// 	blegStr = fmt.Sprintf(blegStr, h.BLegUId, h.VirtualPhone, h.DialplanNumber, h.Gateway, h.CalloutPhone)
// 	if _, h.Err = con.BgApi("originate", blegStr, "&park()"); h.Err != nil {
// 		utils.Error("originate bleg error: %s", h.Err)
// 		con.Close()
// 	}

// }

func (h *Handler) OnDisconnect(con *esl.Connection, ev *esl.Event) {
	log.Println("esl disconnected:", ev)
}

func (h *Handler) OnClose(con *esl.Connection) {
	h.Done <- false
	log.Println("esl connection closed")
}

func (h *Handler) OnEvent(con *esl.Connection, ev *esl.Event) {
	log.Printf("%s - event %s %s %s\n", ev.UId, ev.Name, ev.App, ev.AppData)
	// fmt.Println(ev) // send to stderr as it is very verbose
	// copy code
	resp := ev.GetTextBody()
	if strings.HasPrefix(resp, "-ERR") {
		h.Err = errors.New(resp[5:])
		utils.Error("call terminated with cause %s", resp)
		con.Close()
	}
	// 不是本次通话的uuid的事件不处理
	if !(ev.UId == h.ALegUId || ev.UId == h.BLegUId) {
		return
	}
	// end
	switch ev.Name {
	case esl.CHANNEL_BRIDGE:
		// fmt.Println(ev)
	case esl.CHANNEL_ANSWER:
		// fmt.Println(ev)
		if ev.UId == h.ALegUId {
			h.ALegState[ev.UId] = esl.CHANNEL_ANSWER
		} else if ev.UId == h.BLegUId {
			h.BLegState[ev.UId] = esl.CHANNEL_ANSWER
		}
		utils.Info("channel answer: %v", h)
		if h.ALegState[h.ALegUId] == esl.CHANNEL_ANSWER && h.BLegState[h.BLegUId] == esl.CHANNEL_ANSWER {
			utils.Info("two channel has been answered")
			h.Done <- true
			// 设置录音
			rec := &Record{}
			exts := strings.Split(h.ExtensionNumber, "/")
			rec.Init(exts[1], h.CalloutPhone)
			// uuid_setvar <uuid> <varname> [value]
			recparamsFormat := `%s %s %s`
			con.BgApi("uuid_setvar", fmt.Sprintf(recparamsFormat, h.BLegUId, rec.Name, rec.File))
			// uuid_record <uuid> [start|stop|mask|unmask] <path> [<limit>]
			con.BgApi("uuid_record", fmt.Sprintf(recparamsFormat, h.BLegUId, "start", filepath.Join(rec.PrefixPath, rec.File)))
		}
	case esl.CHANNEL_HANGUP:
		hupcause := ev.Get("Hangup-Cause")
		log.Printf("call terminated with cause %s", hupcause)
		utils.Info("hangup uuid: %s", ev.UId)
		if ev.UId == h.ALegUId {
			if h.BLegState[h.BLegUId] != esl.CHANNEL_HANGUP {
				con.BgApi("uuid_kill", h.BLegUId)
			}
			h.ALegState[ev.UId] = esl.CHANNEL_HANGUP
		} else if ev.UId == h.BLegUId {
			if h.ALegState[h.ALegUId] != esl.CHANNEL_HANGUP {
				con.BgApi("uuid_kill", h.ALegUId)
			}
			h.BLegState[ev.UId] = esl.CHANNEL_HANGUP
		}
		con.Close()
	}
}
