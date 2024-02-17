package gochat

type Group struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	AdditionalInfo map[string]string `json:"additionalInfo,omitempty"`
	users          map[string]*User
	registerUser   chan *User
	unregisterUser chan *User
}

func NewGroup(ID string, name string, additionalInfo map[string]string) *Group {
	return &Group{
		ID:             ID,
		Name:           name,
		AdditionalInfo: additionalInfo,
		users:          make(map[string]*User),
		registerUser:   make(chan *User),
		unregisterUser: make(chan *User),
	}
}

func (group *Group) Run() {
	for {
		select {
		case user := <-group.registerUser:
			group.handleRegisterUser(user)
		case user := <-group.unregisterUser:
			group.handleUnregisterUser(user)

		}
	}
}

func (group *Group) handleRegisterUser(user *User) {
	group.users[user.ID] = user
}

func (group *Group) handleUnregisterUser(user *User) {
	delete(group.users, user.ID)
}
