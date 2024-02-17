package gochat

type (
	Channel struct {
		ID              string            `json:"id"`
		Name            string            `json:"name"`
		AdditionalInfo  map[string]string `json:"additionalInfo,omitempty"`
		users           map[string]*User
		groups          map[string]*Group
		registerUser    chan *User
		unregisterUser  chan *User
		registerGroup   chan *Group
		unregisterGroup chan *Group
	}
)

func NewChannel(ID string, Name string, AdditionalInfo map[string]string) *Channel {
	return &Channel{
		ID:              ID,
		Name:            Name,
		AdditionalInfo:  AdditionalInfo,
		users:           make(map[string]*User),
		groups:          map[string]*Group{},
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
}

func (channel *Channel) handleUnregisterUser(user *User) {
	delete(channel.users, user.ID)
}

func (channel *Channel) handleRegisterGroup(group *Group) {
	channel.groups[group.ID] = group
}

func (channel *Channel) handleUnregisterGroup(group *Group) {
	delete(channel.groups, group.ID)
}

func (channel *Channel) findUserByID(id string) *User {
	return channel.users[id]
}

func (channel *Channel) findGroupByID(id string) *Group {
	return channel.groups[id]
}
