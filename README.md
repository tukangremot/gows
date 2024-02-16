# Gochat

## Usage

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

### Send Direct Message
Payload
```json
{
    "command": "message-send",
    "user": {
        "id": "<your-user-id>",
        "name": "<your-user-name>",
        "additionalInfo": {} // object of string
    },
    "target": {
        "type": "direct",
        "user": {
            "id": "<your-user-target-id>",
            "name": "<your-user-target-name>",
            "additionalInfo": {} // object of string
        }
    },
    "message": {
        "type": "<message-type>", // text
        "text": "<message-text>",
        "additionalInfo": {} // object of string
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