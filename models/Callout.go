package models

import (
    "errors"
)

var AgentCalloutParam map[string][]*Callout

type Callout struct {
    ExtNo    string // 坐席号码
    Gateway  string // 中继号对应的网关
    TrunkNo  string // 中继号
    State    int    // 状态：[-1:等待, 0:准备, 1:已拨，2:已到达上限]
    Count    int    // 已拨次数
    MaxCount int    // 规定的最大次数
}

func CalloutInit() {
    AgentCalloutParam = make(map[string][]*Callout)
    co1 := &Callout{
        ExtNo:    "80000",
        Gateway:  "zqzj",
        TrunkNo:  "28324291",
        State:    -1,
        MaxCount: 3,
    }
    co2 := &Callout{
        ExtNo:    "80000",
        Gateway:  "zqzj",
        TrunkNo:  "28324292",
        State:    -1,
        MaxCount: 3,
    }
    co3 := &Callout{
        ExtNo:    "80000",
        Gateway:  "zqzj",
        TrunkNo:  "28324293",
        State:    -1,
        MaxCount: 3,
    }
    co4 := &Callout{
        ExtNo:    "80000",
        Gateway:  "zqzj",
        TrunkNo:  "28324294",
        State:    -1,
        MaxCount: 3,
    }
    AgentCalloutParam["80000"] = []*Callout{co1, co2, co3, co4}
}

func CalloutPickTrunkNo(cos []*Callout) (*Callout, error) {
    var picked bool // 是否已选到号码
    var index int
    for i, item := range cos {
        // 第一个state为-1的，改为0
        if item.State == -1 && !picked {
            item.State = 0
            picked = true
            index = i
        }
    }
    if !picked {
        for i, item := range cos {
            if item.State == 1 {
                if !picked {
                    item.State = 0
                    picked = true
                    index = i
                } else {
                    item.State = -1
                }
            }
        }
    }
    if picked {
        return cos[index], nil
    }
    return &Callout{}, errors.New("no valid trunk no.")
}

func CalloutSuccess(co *Callout) {
    if co.State == 0 {
        co.State = 1
        co.Count++
    }
    if co.Count >= co.MaxCount {
        co.State = 2
    }
}