package controller

import (
	"Shaw/goWeb/chatRoom/data"
	"Shaw/goWeb/chatRoom/room"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WSHandler(w http.ResponseWriter, r *http.Request, username string) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("failed to upgrade websocket", err)
		return
	}

	rooms, _ := data.GetDB().GetRoom(username)

	user := &room.User{
		Name:  username,
		Hub:   room.GetHub(),
		Rooms: rooms,
		Conn:  conn,
		Send:  make(chan room.Msg, 64),
	}
	log.Println(user.Name)
	user.Hub.RegisterUser <- user
	log.Println("has register", user.Name)
	go user.Read()
	go user.Write()
}
