package dial

import (
    "fmt"
    "log"
    "time"
    "strings"
    "strconv"
    "path/filepath"

    "zq/callout_crm/models"
    "zq/callout_crm/utils"

    "github.com/vma/esl"
)

// 这个是一条腿一条腿的桥接的方式
// type CallinHandler struct {
//     Caller         map[string]string    // [aleg_uuid]bleg_uuid
//     CalloutNumbers map[string]int       // [aleg_uuid]count
//     Callee         map[string]string    // [belg_uuid]aleg_uuid
// }

type CallinHandler struct {
    Caller         map[string]struct{}  // 就想标记一下呼入
}

func CallinInit() {
    callInHandler := new(CallinHandler)
    callInHandler.Caller = make(map[string]struct{})
    // callInHandler.Caller = make(map[string]string)
    // callInHandler.CalloutNumbers = make(map[string]int)
    // callInHandler.Callee = make(map[string]string)

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
    con.SendRecv("event", "plain", "CHANNEL_CREATE CHANNEL_ANSWER CHANNEL_HANGUP")
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
    if !needHandleOnEvent(direction, destinationNumber) {
        utils.Info("no on event -- direction: %s, caller: %s, callee: %s", direction, callerNumber, destinationNumber)
        return
    }

    switch ev.Name {
    case esl.CHANNEL_CREATE:
        utils.Info("channel_create -- direction: %s, caller: %s, callee: %s", direction, callerNumber, destinationNumber)
        // fmt.Println(ev)
        if strings.Compare(direction, "inbound") == 0 {
            // extNo, bindPhone := models.GetExtNoByTrunkNo(destinationNumber)
            tau := models.GetExtNoByTrunkNo(destinationNumber)
            if len(tau.ExtNo) > 0 {
                // con.Execute("set", ev.UId, fmt.Sprintf("dialplan_id=%d", tau.BackendUserId))
                // con.Execute("loop_playback", ev.UId, fmt.Sprintf("+6 %s", "/home/voices/rings/common/ring_short.wav"))
                // 先桥接一端，成功了再桥接
                // if blegUId, err := con.BgApi("create_uuid"); err != nil {
                //     con.Execute("hangup", ev.UId)
                // } else {
                //     // 添加aleg=bleg信息
                //     h.Caller[ev.UId] = blegUId
                //     h.CalloutNumbers[ev.UId] = 2
                //     if _, err := con.BgApi("originate", fmt.Sprintf("{origination_uuid=%s,origination_caller_id_number=%s,effective_caller_id_number=%s}user/%s", blegUId, callerNumber, callerNumber, tau.ExtNo), "&park"); err != nil {
                //         utils.Error("%v\n", err)
                //         // con.Execute("playback", ev.UId, "/home/voices/rings/common/busy.wav")
                //         // con.Execute("hangup", ev.UId)
                //         if len(tau.BindPhone) > 0 {
                //             if blegUId, err := con.BgApi("create_uuid"); err != nil {
                //                 con.Execute("hangup", ev.UId)
                //             } else {
                //                 // aleg还存在才外呼
                //                 if _, ok := h.Caller[ev.UId]; ok {
                //                     // 添加aleg=bleg信息
                //                     h.Caller[ev.UId] = blegUId
                //                     if _, err := con.BgApi("originate", fmt.Sprintf("{origination_uuid=%s,origination_caller_id_number=%s,sip_h_Diversion=<sip:28324284@ip}sofia/gateway/%s/%s", blegUId, callerNumber, tau.Gateway, tau.BindPhone), "&park"); err != nil {
                //                         utils.Error("originate out number error: %v\n", err)
                //                         con.Execute("playback", ev.UId, "/home/voices/rings/common/busy.wav")
                //                         con.Execute("hangup", ev.UId)
                //                     } else {
                //                         con.BgApi("uuid_bridge", ev.UId, blegUId)
                //                     }
                //                 }
                //             }
                //         } else {
                //             _, ok1 := h.Caller[ev.UId]
                //             _, ok2 := h.Callee[blegUId]
                //             if ok1 && ok2 {
                //                 con.Execute("hangup", ev.UId)
                //             }
                //         }
                //     } else {
                //         _, ok1 := h.Caller[ev.UId]
                //         _, ok2 := h.Callee[blegUId]
                //         if ok1 && ok2 {
                //             con.BgApi("uuid_bridge", ev.UId, blegUId)
                //         }
                        
                //     }
                // }
                
                // 直接桥接
                con.Execute("set", ev.UId, "hangup_after_bridge=true")
                con.Execute("set", ev.UId, fmt.Sprintf("dialplan_id=%d", tau.BackendUserId))
                // 呼叫字符串
                calloutStrs := make([]string, 1, 2)
                calloutStrs[0] = fmt.Sprintf("{ringback=/home/voices/rings/common/ring_short.wav,origination_caller_id_number=%s,leg_timeout=?}user/%s", callerNumber, tau.ExtNo)
                if len(tau.BindPhone) > 0 {
                    calloutStrs = append(calloutStrs, fmt.Sprintf("[sip_h_Diversion=<sip:%s@ip>]sofia/gateway/%s/%s", tau.TrunkNo, tau.Gateway, tau.BindPhone))
                }
                strings.Replace(calloutStrs[0], "?", strconv.Itoa(60/len(calloutStrs)), 1)
                con.Execute("bridge", ev.UId, strings.Join(calloutStrs, "|"))
                if _, ok := h.Caller[ev.UId]; !ok {
                    h.Caller[ev.UId] = struct{}{}
                }
                // 推送到WebSocketController echo
                models.WsMessage <- models.WsCallInfo{
                    Caller: callerNumber,
                    Callee: destinationNumber,
                    ExtNo:  tau.ExtNo,
                }
            } else {
                con.Execute("hangup", ev.UId)
            }
        } 
        // else if strings.Compare(direction, "outbound") == 0 {
        //     callerUId := ev.Get("Channel-Call-UUID")
        //     if _, ok := h.Callee[ev.UId]; !ok {
        //         h.Callee[ev.UId] = callerUId
        //     }
        // } else {
        //     con.Execute("hangup", ev.UId)
        // }
        
    case esl.CHANNEL_ANSWER:
        if _, ok := h.Caller[ev.UId]; ok {
            // 设置录音
            rec := &Record{}
            rec.InitCallIn(callerNumber, destinationNumber)
            // uuid_setvar <uuid> <varname> [value]
            recparamsFormat := `%s %s %s`
            con.BgApi("uuid_setvar", fmt.Sprintf(recparamsFormat, ev.UId, rec.Name, rec.File))
            // uuid_record <uuid> [start|stop|mask|unmask] <path> [<limit>]
            con.BgApi("uuid_record", fmt.Sprintf(recparamsFormat, ev.UId, "start", filepath.Join(rec.PrefixPath, rec.File)))
        }

    case esl.CHANNEL_HANGUP:
        utils.Info("channel_hangup -- direction: %s, caller: %s, callee: %s", direction, callerNumber, destinationNumber)
        if _, ok := h.Caller[ev.UId]; ok {
            delete(h.Caller, ev.UId)
        }
        // if bUId, ok := h.Caller[ev.UId]; ok {
        //     if _, ok := h.Callee[bUId]; ok {
        //         con.Execute("hangup", bUId)
        //         delete(h.Callee, bUId)
        //     }
        //     delete(h.Caller, ev.UId)
        //     delete(h.CalloutNumbers, ev.UId)
        // }
        // if aUId, ok := h.Callee[ev.UId]; ok {
        //     if _, ok := h.Caller[aUId]; ok {
        //         if tmp, ok := h.CalloutNumbers[aUId]; ok {
        //             if tmp == 0 {
        //                 con.Execute("hangup", aUId)
        //                 delete(h.Caller, aUId)
        //             } else {
        //                 h.CalloutNumbers[aUId] -= 1
        //             }
                    
        //         }
        //     }
        //     delete(h.Callee, ev.UId)
        // }
    }
}

func needHandleOnEvent(direction, destinationNumber string) bool {
    // li and qin
    if strings.Compare(destinationNumber, "28324286") == 0 ||
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
            "80000", "80001", "80002", "80003", "80004", "13675017141",
        }
        for _, extNo := range extNos {
            if strings.Compare(destinationNumber, extNo) == 0 {
                return true
            }
        }
    }
    

    return false
}