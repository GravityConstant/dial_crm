package dial

import (
    "fmt"
    "log"
    "time"
    "strings"

    "zq/callout_crm/models"
    "zq/callout_crm/utils"

    "github.com/vma/esl"
)

type CallinHandler struct {
    Caller map[string]int    // [aleg_uuid]1
    Callee map[string]string // [belg_uuid]aleg_uuid
}

func CallinInit() {
    callInHandler := new(CallinHandler)
    callInHandler.Caller = make(map[string]int)
    callInHandler.Callee = make(map[string]string)

    go func() {
    RESTART:
        con, err := esl.NewConnection("127.0.0.1:8021", callInHandler)
        if err != nil {
            log.Fatal("ERR connecting to freeswitch:", err)
        }
        
        if err := con.HandleEvents(); err != nil {
            utils.LogError("call_in.go handle events error. " + err.Error())
            time.Sleep(time.Second)
            goto RESTART
        }
    }()
}

func (h *CallinHandler) OnConnect(con *esl.Connection) {
    con.SendRecv("event", "plain", "CHANNEL_CREATE CHANNEL_HANGUP")
    // con.SendRecv("event", "plain", "ALL")
}

func (h *CallinHandler) OnDisconnect(con *esl.Connection, ev *esl.Event) {
    log.Println("esl disconnected:", ev)
}

func (h *CallinHandler) OnClose(con *esl.Connection) {
    log.Println("esl connection closed")
}

func (h *CallinHandler) OnEvent(con *esl.Connection, ev *esl.Event) {
    log.Printf("%s - event %s %s %s\n", ev.UId, ev.Name, ev.App, ev.AppData)
    // fmt.Println(ev) // send to stderr as it is very verbose 182.61.34.55
    direction := ev.Get("Call-Direction")
    callerNumber := ev.Get("Caller-Caller-ID-Number")
    destinationNumber := ev.Get("Caller-Destination-Number")
    if ok := needHandleOnEvent(direction, destinationNumber); !ok {
        return
    }

    switch ev.Name {
    case esl.CHANNEL_CREATE:
        utils.Info("direction: %s, caller: %s, callee: %s", direction, callerNumber, destinationNumber)
        // fmt.Println(ev)
        if strings.Compare(direction, "inbound") == 0 {
            extNo := models.GetExtNoByTrunkNo(destinationNumber)
            if len(extNo) > 0 {
                // if blegUId, err := con.Api("create_uuid"); err != nil {
                //     con.Execute("hangup", ev.UId)
                // } else {
                //     if _, err := con.Api("originate", fmt.Sprintf("{origination_uuid=%s,origination_caller_id_number=%s,effective_caller_id_number=%s}user/%s", blegUId, callerNumber, callerNumber, extNo), "&park"); err != nil {
                //         utils.Error("%v\n", err)
                //         con.Execute("playback", ev.UId, "/home/voices/rings/common/busy.wav")
                //     } else {
                //         con.Api("uuid_bridge", ev.UId, blegUId)
                //     }
                // }
                con.Execute("set", ev.UId, "hangup_after_bridge=true")
                con.Execute("bridge", ev.UId, fmt.Sprintf("{ringback=/home/voices/rings/common/ring_short.wav,origination_caller_id_number=%s}user/%s", callerNumber, extNo))
                if _, ok := h.Caller[ev.UId]; !ok {
                    h.Caller[ev.UId] = 1
                }
                // 推送到WebSocketController echo
                models.WsMessage <- models.WsCallInfo{
                    Caller: callerNumber,
                    Callee: destinationNumber,
                    ExtNo:  extNo,
                }
            } else {
                con.Execute("hangup", ev.UId)
            }
        } else if strings.Compare(direction, "outbound") == 0 {
            callerUId := ev.Get("Channel-Call-UUID")
            if _, ok := h.Callee[ev.UId]; !ok {
                h.Callee[ev.UId] = callerUId
            }
        } else {
            con.Execute("hangup", ev.UId)
        }
        

    case esl.CHANNEL_HANGUP:
        if _, ok := h.Caller[ev.UId]; ok {
            delete(h.Caller, ev.UId)
        }
        if aUId, ok := h.Callee[ev.UId]; ok {
            if _, ok := h.Caller[ev.UId]; ok {
                con.Execute("hangup", aUId)
            }
            delete(h.Callee, ev.UId)
        }
    }
}

func needHandleOnEvent(direction, destinationNumber string) bool {
    // li and qin
    if strings.Compare(destinationNumber, "28324286") == 0 &&
        strings.Compare(destinationNumber, "28324300") == 0 {
            return false
    }

    if strings.Compare(direction, "inbound") == 0 {
        trunkNos := []string{
            "28324284","28324285","28324286","28324287","28324288","28324289",
            "28324290","28324291","28324292","28324293","28324294","28324295","28324296","28324297","28324298","28324299",
            "28324300","28324301","28324302","28324303","28324304","28324305","28324306","28324307","28324308","28324309",
            "28324310","28324311","28324312","28324313",
        }
        for _, trunkNo := range trunkNos {
            if strings.Compare(destinationNumber, trunkNo) == 0 {
                return true
            }
        }
    } else if strings.Compare(direction, "outbound") == 0 {
        extNos := []string{
            "80000", "80001", "80002", "80003", "80004",
        }
        for _, extNo := range extNos {
            if strings.Compare(destinationNumber, extNo) == 0 {
                return true
            }
        }
    }
    

    return false
}