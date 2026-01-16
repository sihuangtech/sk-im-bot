package api

import (
	"encoding/json"
	"net/http"

	"sk-im-bot/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Client 代表一个已连接的 WebSocket 客户端
type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

// Hub 负责维护活跃客户端集合并向客户端广播消息。
type Hub struct {
	// 已注册的客户端。
	Clients map[*Client]bool

	// 来自客户端的入站消息。
	Broadcast chan []byte

	// 来自客户端的注册请求。
	Register chan *Client

	// 来自客户端的注销请求。
	Unregister chan *Client
}

// WSHub 是全局 WebSocket hub 实例
var WSHub = &Hub{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[*Client]bool),
}

// Run 处理 hub 的主循环
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			utils.Logger.Info("WebSocket 客户端已连接")
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			utils.Logger.Info("WebSocket 客户端已断开连接")
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 开发环境下允许所有跨域请求
	},
}

// WSHandler 处理来自对端的 WebSocket 请求。
func WSHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.Logger.Error("WebSocket 连接升级失败", zap.Error(err))
		return
	}
	client := &Client{Hub: WSHub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	go client.writePump()
	go client.readPump()
}

// readPump 将消息从 websocket 连接泵送到 hub。
func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				utils.Logger.Error("WebSocket 错误", zap.Error(err))
			}
			break
		}
	}
}

// writePump 将消息从 hub 泵送到 websocket 连接。
func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()
	for message := range c.Send {
		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		// 将队列中的聊天消息添加到当前的 websocket 消息中。
		n := len(c.Send)
		for i := 0; i < n; i++ {
			w.Write(<-c.Send)
		}

		if err := w.Close(); err != nil {
			return
		}
	}
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}

// BroadcastEvent 向所有连接的客户端发送事件
func BroadcastEvent(event interface{}) {
	jsonBytes, err := json.Marshal(event)
	if err != nil {
		utils.Logger.Error("序列化广播事件失败", zap.Error(err))
		return
	}
	WSHub.Broadcast <- jsonBytes
}
