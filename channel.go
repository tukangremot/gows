package gochat

type (
	Channel struct {
		ID             string            `json:"id"`
		Name           string            `json:"name"`
		AdditionalInfo map[string]string `json:"additionalInfo,omitempty"`
		users          map[string]*User
		registerUser   chan *User
		unregisterUser chan *User
	}
)

func NewChannel(ID string, Name string, AdditionalInfo map[string]string) *Channel {
	return &Channel{
		ID:             ID,
		Name:           Name,
		AdditionalInfo: AdditionalInfo,
		users:          make(map[string]*User),
		registerUser:   make(chan *User),
		unregisterUser: make(chan *User),
	}

}

func (channel *Channel) Run() {
	for {
		select {
		case user := <-channel.registerUser:
			channel.handleRegisterUser(user)

		case user := <-channel.unregisterUser:
			channel.handleUnregisterUser(user)
		}
	}
}

func (channel *Channel) handleRegisterUser(user *User) {
	channel.users[user.ID] = user
}

func (channel *Channel) handleUnregisterUser(user *User) {
	delete(channel.users, user.ID)
}

func (channel *Channel) getUserByID(id string) *User {
	return channel.users[id]
}
