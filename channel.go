package gochat

import (
	"log"
)

type (
	Channel struct {
		ID              string            `json:"id"`
		Name            string            `json:"name"`
		AdditionalInfo  map[string]string `json:"additionalInfo,omitempty"`
		users           map[string]*User
		groups          map[string]*Group
		server          *Server
		registerUser    chan *User
		unregisterUser  chan *User
		registerGroup   chan *Group
		unregisterGroup chan *Group
	}
)

func NewChannel(ID string, Name string, AdditionalInfo map[string]string, server *Server) *Channel {
	return &Channel{
		ID:              ID,
		Name:            Name,
		AdditionalInfo:  AdditionalInfo,
		users:           make(map[string]*User),
		groups:          map[string]*Group{},
		server:          server,
		registerUser:    make(chan *User),
		unregisterUser:  make(chan *User),
		registerGroup:   make(chan *Group),
		unregisterGroup: make(chan *Group),
	}

}

func (channel *Channel) Run() {
	for {
		select {
		case user := <-channel.registerUser:
			channel.handleRegisterUser(user)
		case user := <-channel.unregisterUser:
			channel.handleUnregisterUser(user)
		case group := <-channel.registerGroup:
			channel.handleRegisterGroup(group)
		case group := <-channel.unregisterGroup:
			channel.handleUnregisterGroup(group)
		}
	}
}

func (channel *Channel) handleRegisterUser(user *User) {
	channel.users[user.ID] = user

	if channel.server.Session != nil {
		err := channel.server.Session.registerUserChannel(channel.server.ctx, channel, user)
		if err != nil {
			log.Println(err)
		}
	}
}

func (channel *Channel) handleUnregisterUser(user *User) {
	delete(channel.users, user.ID)

	if channel.server.Session != nil {
		err := channel.server.Session.unregisterUserChannel(channel.server.ctx, channel, user)
		if err != nil {
			log.Println(err)
		}
	}
}

func (channel *Channel) handleRegisterGroup(group *Group) {
	channel.groups[group.ID] = group
}

func (channel *Channel) handleUnregisterGroup(group *Group) {
	delete(channel.groups, group.ID)
}

func (channel *Channel) findUserByID(id string) *User {
	if channel.users[id] != nil {
		return channel.users[id]
	}

	if channel.server.Session != nil {
		user, err := channel.server.Session.findUserChannelByID(channel.server.ctx, channel, id)
		if err != nil {
			log.Println(err)

			return nil
		}

		return user
	}

	return nil
}

func (channel *Channel) findGroupByID(id string) *Group {
	return channel.groups[id]
}

func (channel *Channel) getUsersByGroup(group *Group) map[string]*User {
	usersGroupTarget := channel.groups[group.ID].users

	if channel.server.Session != nil {
		usersGroupInSession, err := channel.server.Session.getUsersByGroup(channel.server.ctx, group)
		if err != nil {
			log.Println(err)

			return usersGroupTarget
		}

		for userGroupInSessionID, userGroupTarget := range usersGroupInSession {
			if usersGroupTarget[userGroupInSessionID] == nil {
				userGroupTarget.onDifferentServer = true
				usersGroupTarget[userGroupInSessionID] = userGroupTarget
			}
		}
	}

	return usersGroupTarget
}
