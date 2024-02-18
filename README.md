# Gochat

## Usage

### Example

#### Simple
```go
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tukangremot/gochat"
)

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(server *gochat.Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	user := gochat.NewUser(conn, server)

	go user.WritePump()
	go user.ReadPump()
	go user.ReadActivity()
}

func main() {
	wsServer := gochat.NewServer(nil)
	go wsServer.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(wsServer, w, r)
	})

	httpServer := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```

#### Capture User Activity
```go
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tukangremot/gochat"
)

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(server *gochat.Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	user := gochat.NewUser(conn, server)

	go user.WritePump()
	go user.ReadPump()

	for activity := range user.GetActivity() {
		switch activity.Type {
		case gochat.TypeUserActivityChannelConnect:
			// do somthing when user connect to channel
		case gochat.TypeUserActivityGroupJoin:
			// do somthing when user join to group
		case gochat.TypeUserActivityGroupLeave:
			// do somthing when user leave from group
		case gochat.TypeUserActivityMessageSend:
			// do somthing when user send message
		case gochat.TypeUserActivityDisconnect:
			// do somthing when user disconnet
		}
	}

}

func main() {
	wsServer := gochat.NewServer(nil)
	go wsServer.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(wsServer, w, r)
	})

	httpServer := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```

#### Using Session and PubSub (supports horizontal scaling)
```go
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tukangremot/gochat"
)

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(server *gochat.Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	user := gochat.NewUser(conn, server)

	go user.WritePump()
	go user.ReadPump()

	for activity := range user.GetActivity() {
		switch activity.Type {
		case gochat.TypeUserActivityChannelConnect:
			// do somthing when user connect to channel
		case gochat.TypeUserActivityGroupJoin:
			// do somthing when user join to group
		case gochat.TypeUserActivityGroupLeave:
			// do somthing when user leave from group
		case gochat.TypeUserActivityMessageSend:
			// do somthing when user send message
		case gochat.TypeUserActivityDisconnect:
			// do somthing when user disconnet
		}
	}

}

func main() {
	wsServer := gochat.NewServer(nil)
	go wsServer.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(wsServer, w, r)
	})

	httpServer := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```

### Events

#### User Join

Payload
```json
{
    "command": "user-connect",
    "channel": {
        "id": "<your-channel-id>",
        "name": "<your-channel-name",
        "additionalInfo": {} // object of string
    },
    "user": {
        "id": "<your-user-id>",
        "name": "<your-user-name",
        "additionalInfo": {} // object of string
    }
}
```
Example:
```json
{
    "command": "user-connect",
    "channel": {
        "id": "1",
        "name": "Channel 1",
        "additionalInfo": {
            "icon": "https://example.com/avatar.jpg"
        }
    },
    "user": {
        "id": "1",
        "name": "John",
        "additionalInfo": {
            "avatar": "https://example.com/avatar.jpg"
        }
    }
}
```

Response
```json
{
    "command": "user-connect",
    "channel": {
        "id": "1",
        "name": "Channel 1",
        "additionalInfo": {
            "icon": "https://example.com/avatar.jpg"
        }
    },
    "user": {
        "id": "1",
        "name": "John",
        "additionalInfo": {
            "avatar": "https://example.com/avatar.jpg"
        }
    },
    "message": {
        "type": "text",
        "text": "connected successfully"
    },
    "response": {
        "status": true,
        "message": "success"
    }
}
```


### Group Join
Payload
```json
{
    "command": "group-join",
    "group": {
        "id": "<your-group-id>",
        "name": "<your-group-name>",
        "additionalInfo": {} // object of string
    }
}
```

Example
```json
{
    "command": "group-join",
    "group": {
        "id": "1",
        "name": "Group 1",
        "additionalInfo": {
            "icon": "https://example.com/avatar.jpg"
        }
    }
}
```

Response
```json
{
    "command": "group-join",
    "group": {
        "id": "1",
        "name": "Group 1",
        "additionalInfo": {
            "icon": "https://example.com/avatar.jpg"
        }
    },
    "user": {
        "id": "1",
        "name": "John",
        "additionalInfo": {
            "avatar": "https://example.com/avatar.jpg"
        }
    },
    "message": {
        "type": "text",
        "text": "join group successful"
    },
    "response": {
        "status": true,
        "message": "success"
    }
}
```

### Group Leave
Payload
```json
{
    "command": "group-leave",
    "group": {
        "id": "<your-group-id>"
    }
}
```

Example
```json
{
    "command": "group-leave",
    "group": {
        "id": "1"
    }
}
```

Response
```json
{
    "command": "group-leave",
    "group": {
        "id": "1",
        "name": "Group 1",
        "additionalInfo": {
            "icon": "https://example.com/avatar.jpg"
        }
    },
    "user": {
        "id": "1",
        "name": "John",
        "additionalInfo": {
            "avatar": "https://example.com/avatar.jpg"
        }
    },
    "message": {
        "type": "text",
        "text": "leave group successful"
    },
    "response": {
        "status": true,
        "message": "success"
    }
}
```

### Send Direct Message
Payload
```json
{
    "command": "message-send",
    "target": {
        "type": "direct",
        "user": {
            "id": "<your-user-target-id>"
        }
    },
    "message": {
        "type": "<message-type>", // text
        "text": "<message-text>",
        "additionalInfo": {} // object of string
    }
}
```

Example
```json
{
    "command": "message-send",
    "target": {
        "type": "direct",
        "user": {
            "id": "2"
        }
    },
    "message": {
        "type": "text",
        "text": "Hallo Emma",
        "additionalInfo": {}
    }
}

```

Response
```json
{
    "command": "message-send",
    "user": {
        "id": "1",
        "name": "John",
        "additionalInfo": {
            "avatar": "https://example.com/avatar.jpg"
        }
    },
    "message": {
        "type": "text",
        "text": "halo",
        "additionalInfo": {
            "link": "https://example.com/avatar.jpg"
        }
    },
    "target": {
        "type": "direct",
        "user": {
            "id": "2",
            "name": "Emma",
            "additionalInfo": {
                "avatar": "https://example.com/avatar.jpg"
            }
        }
    },
    "response": {
        "status": true,
        "message": "success"
    }
}
```

The message received by the target
```json
{
    "command": "message-send",
    "user": {
        "id": "1",
        "name": "John",
        "additionalInfo": {
            "avatar": "https://example.com/avatar.jpg"
        }
    },
    "message": {
        "type": "text",
        "text": "halo",
        "additionalInfo": {
            "link": "https://example.com/avatar.jpg"
        }
    },
    "target": {
        "type": "direct",
        "user": {
            "id": "2",
            "name": "Emma",
            "additionalInfo": {
                "avatar": "https://example.com/avatar.jpg"
            }
        }
    }
}
```

### Send Group Message
Payload
```json
{
    "command": "message-send",
    "target": {
        "type": "group",
        "group": {
            "id": "<your-group-target-id>"
        }
    },
    "message": {
        "type": "<message-type>", // text
        "text": "<message-text>",
        "additionalInfo": {} // object of string
    }
}
```

Example
```json
{
    "command": "message-send",
    "target": {
        "type": "group",
        "group": {
            "id": "1"
        }
    },
    "message": {
        "type": "text",
        "text": "Hallo Guys",
        "additionalInfo": {}
    }
}
```

Response
```json
{
    "command": "message-send",
    "user": {
        "id": "1",
        "name": "John",
        "additionalInfo": {
            "avatar": "https://example.com/avatar.jpg"
        }
    },
    "message": {
        "type": "text",
        "text": "halo",
        "additionalInfo": {
            "link": "https://example.com/avatar.jpg"
        }
    },
    "target": {
        "type": "group",
        "user": {
            "id": "1",
            "name": "Group 1",
            "additionalInfo": {
                "avatar": "https://example.com/avatar.jpg"
            }
        }
    },
    "response": {
        "status": true,
        "message": "success"
    }
}
```

The message received by the target
```json
{
    "command": "message-send",
    "user": {
        "id": "1",
        "name": "John",
        "additionalInfo": {
            "avatar": "https://example.com/avatar.jpg"
        }
    },
    "message": {
        "type": "text",
        "text": "halo",
        "additionalInfo": {
            "link": "https://example.com/avatar.jpg"
        }
    },
    "target": {
        "type": "group",
        "user": {
            "id": "1",
            "name": "Group 1",
            "additionalInfo": {
                "avatar": "https://example.com/avatar.jpg"
            }
        }
    }
}
```
