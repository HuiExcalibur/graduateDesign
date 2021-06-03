package room

import (
	"Shaw/goWeb/chatRoom/data"
	"fmt"
	"log"
	"sync"
)

type Hub struct {
	Rooms map[string]*Room

	Users map[string]*User

	Message chan Msg

	RegisterRoom chan *Room

	UnRegisterRoom chan *Room

	RegisterUser chan *User

	UnRegisterUser chan *User
}

var MyHub *Hub

func GetHub() *Hub {
	log.Println("try to get Hub")
	if MyHub != nil {
		return MyHub
	}

	var once sync.Once
	once.Do(func() {
		log.Println("try to new a Hub")
		newHub()
	})

	return MyHub
}

func newHub() {
	MyHub = &Hub{
		Rooms:          make(map[string]*Room),
		Users:          make(map[string]*User),
		Message:        make(chan Msg),
		RegisterRoom:   make(chan *Room),
		UnRegisterRoom: make(chan *Room),
		RegisterUser:   make(chan *User),
		UnRegisterUser: make(chan *User),
	}

	go MyHub.run()
	log.Println("get a new Hub")

	rooms, err := data.GetDB().GetAllRoom()
	if err != nil {
		log.Println("get a new hub failed ", err)
	}

	for _, v := range rooms {
		log.Println("register room ", v, len(rooms))
		room := NewRoom(v)
		MyHub.RegisterRoom <- room
	}
	log.Println("init a Hub success")
}

func (h Hub) run() {
	for {
		select {
		case room := <-h.RegisterRoom:
			h.Rooms[room.name] = room
			fmt.Println("new room", room.name)
			go room.Run()

		case room := <-h.UnRegisterRoom:
			delete(h.Rooms, room.name)

		case user := <-h.RegisterUser:
			fmt.Println("a new user", user)
			h.Users[user.Name] = user
			for _, v := range user.Rooms {
				fmt.Println("try to enter a room", v)
				h.Rooms[v].Register <- user
			}

		case user := <-h.UnRegisterUser:
			if _, ok := h.Users[user.Name]; ok {
				for _, v := range user.Rooms {
					h.Rooms[v].UnRegister <- user
				}
				delete(h.Users, user.Name)
			}

		case message := <-h.Message:
			log.Println("hub get message", message)
			h.Rooms[message.Room].broadcast <- message
			// err := data.GetDB().NewMessage(message.Username, message.Room, message.Data)
			// if err != nil {
			// 	log.Println("hub try to write a message failed ", err)
			// }
			// for room, name := range h.rooms {
			// 	if name == target {
			// 		room.broadcast <- message
			// 		break
			// 	}
			// }
		}
	}
}
