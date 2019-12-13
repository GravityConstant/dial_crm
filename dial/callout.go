package dial

import (
    "errors"
    "strings"
    "sync"

    "zq/callout_crm/models"
)

// ["80000"]["M"] = []*Callout
var AgentCalloutParam map[string]*sync.Map
// [TrunkNo+GwType]ExtNo
var DialResource = make(map[string]string)

type Callout struct {
    ExtNo    string // 坐席号码
    Gateway  string // 中继号对应的网关
    TrunkNo  string // 中继号
    State    int    // 状态：[-1:等待, 0:准备, 1:已拨，2:已到达上限]
    Count    int    // 已拨次数
    MaxCount int    // 规定的最大次数
}

func CalloutInit() {
    AgentCalloutParam = make(map[string]*sync.Map)

    // 从数据库中取值
    var params models.GatewayPhoneQueryParam
    list := models.GatewayPhoneDataList(&params)
    CalloutUpdate(list)
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
    if !picked {
        // 临时修改了最大数量
        for i, item := range cos {
            if item.MaxCount > item.Count && item.State == 2 {
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

func NewCallout() *Callout {
    return &Callout{
        State: -1,
    }
}

func (co *Callout) SetMaxCount(maxCount int) {
    co.MaxCount = maxCount
}

func (co *Callout) SetCount(count int) {
    co.Count = count
}

func (co *Callout) IncCount() {
    if co.State == 0 {
        co.State = 1
        co.Count++
    }
    if co.Count >= co.MaxCount {
        co.State = 2
    }
}

// 获取当前坐席要被轮询的一组中继号码，分电话线中继和插手机卡中继
// 第一次错误如果返回的是try的话，要再调用一次这个函数
func CalloutGetTrunkPhones(extNo string, gwType int, try bool) ([]*Callout, error) {
    empty := make([]*Callout, 0)

    // 获取网关类型
    gwTypeStr := CalloutGetGwType(gwType)
    // 从map中获取坐席信息
    if agentInfo, ok := AgentCalloutParam[extNo]; ok {
        if ccs, ok := agentInfo.Load(gwTypeStr); ok {
            if yccs, ok := ccs.([]*Callout); ok {
                return yccs, nil
            }
        } else {
            // 不是第二次进来重试的话，加载数据库的值
            if !try {
                // query db
                var params models.GatewayPhoneQueryParam
                params.ExtNoLike = extNo
                list := models.GatewayPhoneDataList(&params)
                CalloutUpdate(list)
                return empty, errors.New("try")
            } else {
                // if not found, switch gateway type
                gwTypeStr = CalloutSwitchGwType(gwType)
                if ccs, ok := agentInfo.Load(gwTypeStr); ok {
                    if yccs, ok := ccs.([]*Callout); ok {
                        return yccs, nil
                    }
                }
            }
        }
    } else {
        if !try {
            // query db
            var params models.GatewayPhoneQueryParam
            params.ExtNoLike = extNo
            list := models.GatewayPhoneDataList(&params)
            CalloutUpdate(list)
            return empty, errors.New("try")
        }
    }

    return empty, errors.New("not found")
}

// 获取网关类型
func CalloutGetGwType(gwType int) string {
    var gwTypeStr string

    switch gwType {
    case 0:
        gwTypeStr = "P"
    case 1:
        gwTypeStr = "M"
    default:
        gwTypeStr = "P"
    }

    return gwTypeStr
}

// 切换网关类型
func CalloutSwitchGwType(gwType int) string {
    var gwTypeStr string

    switch gwType {
    case 0:
        gwTypeStr = "M"
    case 1:
        gwTypeStr = "P"
    default:
        gwTypeStr = "P"
    }

    return gwTypeStr
}


// 更新数据
func CalloutUpdate(list []*models.GatewayPhoneWithRelated) {
    for _, item := range list {
        if len(item.ExtNo) > 0 && len(item.GatewayUrl) > 0 {
            // if exist and extno or gwtype change, delete
            CalloutDelModifiedExtnoOrGwtype(item)
            // create
            m := NewCallout()
            m.ExtNo = strings.TrimSpace(item.ExtNo)
            m.Gateway = strings.TrimSpace(item.GatewayUrl)
            m.TrunkNo = strings.TrimSpace(item.Phone)
            if _, ok := AgentCalloutParam[m.ExtNo]; ok {
                if ccs, ok := AgentCalloutParam[m.ExtNo].Load("P"); ok {
                    if item.GatewayType == 0 {
                        if yccs, ok := ccs.([]*Callout); ok {
                            if !CalloutIsExist(item, yccs) {
                                yccs = append(yccs, m)
                                AgentCalloutParam[m.ExtNo].Store("P", yccs)
                                continue
                            }
                        }
                    }
                }
                if ccs, ok := AgentCalloutParam[m.ExtNo].Load("M"); ok {
                    if item.GatewayType == 1 {
                        if yccs, ok := ccs.([]*Callout); ok {
                            if !CalloutIsExist(item, yccs) {
                                yccs = append(yccs, m)
                                AgentCalloutParam[m.ExtNo].Store("M", yccs)
                                continue
                            }
                        }
                    }

                }
                ccs := make([]*Callout, 0)
                ccs = append(ccs, m)
                if item.GatewayType == 0 {
                    AgentCalloutParam[m.ExtNo].Store("P", ccs)
                } else if item.GatewayType == 1 {
                    AgentCalloutParam[m.ExtNo].Store("M", ccs)
                }

            } else {
                var mSync sync.Map
                AgentCalloutParam[m.ExtNo] = &mSync
                ccs := make([]*Callout, 0)
                ccs = append(ccs, m)
                if item.GatewayType == 0 {
                    AgentCalloutParam[m.ExtNo].Store("P", ccs)
                } else if item.GatewayType == 1 {
                    AgentCalloutParam[m.ExtNo].Store("M", ccs)
                }
            }
        }
    }
}

// 坐席或网关改动的话，丢弃，让它update
func CalloutDelModifiedExtnoOrGwtype(item *models.GatewayPhoneWithRelated) {
    gwType := CalloutGetGwType(item.GatewayType)
    if extNo, ok := DialResource[item.Phone+gwType]; ok {
        // fmt.Println("original extNo: ", extNo)
        if ccs, ok := AgentCalloutParam[extNo].Load(gwType); ok {
            if yccs, ok := ccs.([]*Callout); ok {
                holdCcs := make([]*Callout, 0)
                for _, cc := range yccs {
                    // fmt.Println("callout: ", cc)
                    if cc.TrunkNo == item.Phone {
                        if cc.ExtNo != item.ExtNo || cc.Gateway != item.GatewayUrl {
                            // fmt.Println("depreciate")
                            // 删除相应中继号的资源
                            delete(DialResource, item.Phone+gwType)
                        } else {
                            holdCcs = append(holdCcs, cc)
                        }
                    } else {
                        holdCcs = append(holdCcs, cc)
                    }
                }
                AgentCalloutParam[extNo].Store(gwType, holdCcs)
            }
        }
    } else {
        // 新的网关改变的情况
        switchGwType := CalloutSwitchGwType(item.GatewayType)
        if extNo, ok := DialResource[item.Phone+switchGwType]; ok {
            if ccs, ok := AgentCalloutParam[extNo].Load(switchGwType); ok {
                if yccs, ok := ccs.([]*Callout); ok {
                    holdCcs := make([]*Callout, 0)
                    for _, cc := range yccs {
                        if cc.TrunkNo == item.Phone {
                            // 删除相应中继号的资源
                            delete(DialResource, item.Phone+switchGwType)
                        } else {
                            holdCcs = append(holdCcs, cc)
                        }
                    }
                    AgentCalloutParam[extNo].Store(switchGwType, holdCcs)
                }
            }
        }
    }
}

// 从数据库获取的网关号码是否已插入到AgentCalloutParam
func CalloutIsExist(item *models.GatewayPhoneWithRelated, data []*Callout) bool {
    for _, cc := range data {
        if item.Phone == cc.TrunkNo {
            return true
            break
        }
    }
    return false
}

// 删除整个AgentCalloutParam
func CalloutDestory() {
    for key, _ := range AgentCalloutParam {
        delete(AgentCalloutParam, key)
    }
}

// 刷新整个AgentCalloutParam
func CalloutRefresh() {
    CalloutDestory()
    // 从数据库中取值
    var params models.GatewayPhoneQueryParam
    list := models.GatewayPhoneDataList(&params)
    CalloutUpdate(list)
}