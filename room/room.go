package room

type Room struct {
	name string

	users map[*User]bool

	broadcast chan Msg

	Register chan *User

	UnRegister chan *User
}

func NewRoom(roomname string) *Room {
	return &Room{
		name:       roomname,
		users:      make(map[*User]bool),
		broadcast:  make(chan Msg),
		Register:   make(chan *User),
		UnRegister: make(chan *User),
	}
}

func (r *Room) Run() {
	for {
		select {
		case user := <-r.Register:
			r.users[user] = true
		case user := <-r.UnRegister:
			if _, ok := r.users[user]; ok {
				close(user.Send)
				delete(r.users, user)
			}
		case message := <-r.broadcast:
			for user := range r.users {
				select {
				case user.Send <- message:
				default:
					close(user.Send)
					delete(r.users, user)
				}
			}
		}
	}
}
