package main

import (
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
	addr string

	room *Room

	conn *websocket.Conn

	send chan Msg
}

type Msg struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	Room     string `json:"roomname"`
}

// var (
// 	newline = []byte{'\n'}
// 	space   = []byte{' '}
// )

func (u *User) read() {
	defer func() {
		u.room.unRegister <- u
		u.conn.Close()
	}()

	u.conn.SetReadLimit(maxMessageSize)
	u.conn.SetReadDeadline(time.Now().Add(pongWait))
	u.conn.SetPongHandler(func(string) error {
		u.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var message Msg
		err := u.conn.ReadJSON(&message)
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
		u.room.broadcast <- message
	}
}

func (u *User) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		u.conn.Close()
	}()

	for {
		select {
		case message, ok := <-u.send:
			u.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				u.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := u.conn.WriteJSON(message)
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
			u.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := u.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
