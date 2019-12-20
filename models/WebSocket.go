package models

type WsCallInfo struct {
    Caller string
    Callee string
    ExtNo  string
}

var (
    WsMessage = make(chan WsCallInfo)
)