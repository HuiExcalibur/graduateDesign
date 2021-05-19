package room

import (
	"fmt"
	"sync"
)

type Hub struct {
	Rooms map[string]*Room

	Users map[string]*User

	message chan Msg

	RegisterRoom chan *Room

	UnRegisterRoom chan *Room

	RegisterUser chan *User

	UnRegisterUser chan *User
}

var MyHub *Hub

func GetHub() *Hub {
	if MyHub != nil {
		return MyHub
	}

	var once sync.Once
	once.Do(func() {
		NewHub()
	})

	return MyHub
}

func NewHub() {
	MyHub = &Hub{
		Rooms:          make(map[string]*Room),
		Users:          make(map[string]*User),
		message:        make(chan Msg),
		RegisterRoom:   make(chan *Room),
		UnRegisterRoom: make(chan *Room),
		RegisterUser:   make(chan *User),
		UnRegisterUser: make(chan *User),
	}
}

func (h Hub) Run() {
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
		case message := <-h.message:
			h.Rooms[message.Room].broadcast <- message
			// for room, name := range h.rooms {
			// 	if name == target {
			// 		room.broadcast <- message
			// 		break
			// 	}
			// }
		}
	}
}
