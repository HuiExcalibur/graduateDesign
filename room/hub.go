package room

import "fmt"

type Hub struct {
	rooms map[string]*Room

	users map[*User]bool

	message chan Msg

	RegisterRoom chan *Room

	UnRegisterRoom chan *Room

	RegisterUser chan *User

	UnRegisterUser chan *User
}

func NewHub() *Hub {
	return &Hub{
		rooms:          make(map[string]*Room),
		users:          make(map[*User]bool),
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
			h.rooms[room.name] = room
			fmt.Println("new room", room.name)
			go room.Run()
		case room := <-h.UnRegisterRoom:
			delete(h.rooms, room.name)
		case user := <-h.RegisterUser:
			fmt.Println("a new user", user)
			h.users[user] = true
			for _, v := range user.Rooms {
				fmt.Println("try to enter a room", v)
				h.rooms[v].Register <- user
			}
		case user := <-h.UnRegisterUser:
			if _, ok := h.users[user]; ok {
				for _, v := range user.Rooms {
					h.rooms[v].UnRegister <- user
				}
				delete(h.users, user)
			}
		case message := <-h.message:
			h.rooms[message.Room].broadcast <- message
			// for room, name := range h.rooms {
			// 	if name == target {
			// 		room.broadcast <- message
			// 		break
			// 	}
			// }
		}
	}
}
