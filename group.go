package gochat

import "log"

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

	if user.server.Session != nil {
		err := user.server.Session.registerUserGroup(user.server.ctx, group, user)
		if err != nil {
			log.Println(err)
		}
	}
}

func (group *Group) handleUnregisterUser(user *User) {
	delete(group.users, user.ID)

	if user.server.Session != nil {
		err := user.server.Session.unregisterUserGroup(user.server.ctx, group, user)
		if err != nil {
			log.Println(err)
		}
	}
}
