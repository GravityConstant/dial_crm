package controllers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "time"
    "strings"

    "zq/callout_crm/models"
    "zq/callout_crm/utils"

    "github.com/gorilla/websocket"
)

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
    BaseController
}

const (
    // Time allowed to write a message to the peer.
    writeWait = 10 * time.Second

    // Time allowed to read the next pong message from the peer.
    pongWait = 60 * time.Second

    // Send pings to peer with this period. Must be less than pongWait.
    pingPeriod = (pongWait * 9) / 10

    // Maximum message size allowed from peer.
    maxMessageSize = 512
)

var (
    newline = []byte{'\n'}
    space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
    hub *Hub

    // The websocket connection.
    conn *websocket.Conn

    // Buffered channel of outbound messages.
    send chan []byte

    // the same extNo send
    extNo string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
    for {
        msg := <- models.WsMessage
        message, _ := json.Marshal(msg)
        message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
        c.hub.broadcast <- message
    }
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()
    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                // The hub closed the channel.
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            // Add queued chat messages to the current websocket message.
            n := len(c.send)
            for i := 0; i < n; i++ {
                w.Write(newline)
                w.Write(<-c.send)
            }

            if err := w.Close(); err != nil {
                return
            }
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

var HubInstance *Hub

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
    // Registered clients.
    clients map[*Client]bool

    // Inbound messages from the clients.
    broadcast chan []byte

    // Register requests from the clients.
    register chan *Client

    // Unregister requests from clients.
    unregister chan *Client
}

func CallInHubInit() {
    HubInstance = newHub()
    go HubInstance.run()
}

func newHub() *Hub {
    return &Hub{
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        clients:    make(map[*Client]bool),
    }
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
        case message := <-h.broadcast:
            for client := range h.clients {
                var msg models.WsCallInfo
                json.Unmarshal(message, &msg)
                if msg.ExtNo == client.extNo {
                    select {
                    case client.send <- message:
                    default:
                        close(client.send)
                        delete(h.clients, client)
                    }
                }
            }
        }
    }
}

// Join method handles WebSocket requests for WebSocketController.
func (c *WebSocketController) Echo() {
    utils.Info("RemoteAddr: %s\n", c.Ctx.Request.RemoteAddr)
    // Upgrade from http request to WebSocket.
    conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
    if err != nil {
        utils.LogError("Cannot setup WebSocket connection:" + err.Error())
        return
    }
    client := &Client{hub: HubInstance, conn: conn, send: make(chan []byte, 256), extNo: c.curUser.ExtNo}
    client.hub.register <- client

    // Allow collection of memory referenced by the caller by doing all work in
    // new goroutines.
    go client.readPump()
    client.writePump()
    

    c.StopRun()
}

// 收
func readPump(ws *websocket.Conn, done chan struct{}) {
    for {
        mt, message, err := ws.ReadMessage()
        if err != nil {
            utils.LogError("read:" + err.Error())
            done <- struct{}{}
            break
        }
        utils.Info("recv: %s, messagetype: %v", message, mt)
    }

}

// 发
func writePump() {

}

// 20191228写的
func readWriteBak(ws *websocket.Conn, extNo string) {
    for {
        mt, message, err := ws.ReadMessage()
        if err != nil {
            utils.LogError("read:" + err.Error())
            break
        }
        utils.Info("recv: %s, messagetype: %v", message, mt)

        msg := <- models.WsMessage
        utils.Info("%v, %v\n", msg, extNo)
        if strings.Compare(extNo, msg.ExtNo) != 0 {
            msg.ExtNo = "0000"
            
        }
        data, err := json.Marshal(msg)
        if err != nil {
            utils.LogError("Fail to marshal event:" + err.Error())
            break
        }
        err = ws.WriteMessage(websocket.TextMessage, data)
        if err != nil {
            utils.LogError("write:" + err.Error())
            break
        }
    }
}

// 定时发送
func writeTicker(ws *websocket.Conn, mt int, message []byte, done chan struct{}) {
    if mt == -1 {
        mt = websocket.TextMessage
    }
    if len(message) == 0 {
        message = []byte("Hello World")
    }
    utils.Info("here? %s", message)
    pingTicker := time.NewTicker(time.Second)
    for {
        select {
        case <-pingTicker.C:
            err := ws.WriteMessage(mt, message)
            if err != nil {
                fmt.Println("write:", err)
                done <- struct{}{}
                break
            }
        }
        
    }
}
