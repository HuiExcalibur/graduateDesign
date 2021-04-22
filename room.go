package main

type Room struct {
	users map[*User]bool

	broadcast chan Msg

	register chan *User

	unRegister chan *User
}

func newRoom() *Room {
	return &Room{
		users:      make(map[*User]bool),
		broadcast:  make(chan Msg),
		register:   make(chan *User),
		unRegister: make(chan *User),
	}
}

func (r *Room) run() {
	for {
		select {
		case user := <-r.register:
			r.users[user] = true
		case user := <-r.unRegister:
			if _, ok := r.users[user]; ok {
				close(user.send)
				delete(r.users, user)
			}
		case message := <-r.broadcast:
			for user := range r.users {
				select {
				case user.send <- message:
				default:
					close(user.send)
					delete(r.users, user)
				}
			}
		}
	}
}
