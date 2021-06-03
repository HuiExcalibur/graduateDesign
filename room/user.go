package room

import (
	"Shaw/goWeb/chatRoom/data"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

type User struct {
	Name string

	Hub *Hub

	Rooms []string

	Conn *websocket.Conn

	Send chan Msg
}

type Msg struct {
	//发送消息的用户名称
	Username string `json:"username"`

	//发送的消息
	Data string `json:"message"`

	//消息所属的房间
	Room string `json:"roomname"`
}

// var (
// 	newline = []byte{'\n'}
// 	space   = []byte{' '}
// )

func (u *User) Read() {
	defer func() {
		u.Hub.UnRegisterUser <- u
		u.Conn.Close()
	}()

	u.Conn.SetReadLimit(maxMessageSize)
	u.Conn.SetReadDeadline(time.Now().Add(pongWait))
	u.Conn.SetPongHandler(func(string) error {
		u.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var message Msg
		err := u.Conn.ReadJSON(&message)
		if err != nil {
			fmt.Println("read json error :", err)
			break
		}
		// _, message, err := u.conn.ReadMessage()
		// if err != nil {
		// 	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		// 		log.Println("error: ", err)
		// 	}
		// 	break
		// }

		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		fmt.Println(message)
		data.GetDB().NewMessage(message.Username, message.Room, message.Data)
		u.Hub.Message <- message
	}
}

func (u *User) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		u.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-u.Send:
			u.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				u.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := u.Conn.WriteJSON(message)
			if err != nil {
				fmt.Println("send json error ", err)
				return
			}

			// w, err := u.conn.NextWriter(websocket.TextMessage)
			// if err != nil {
			// 	return
			// }
			// msg := string(message) + u.addr
			// w.Write([]byte(msg))

			// n := len(u.send)
			// for i := 0; i < n; i++ {
			// 	//w.Write(newline)
			// 	w.Write(<-u.send)
			// }

			// if err := w.Close(); err != nil {
			// 	return
			// }
		case <-ticker.C:
			u.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := u.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
