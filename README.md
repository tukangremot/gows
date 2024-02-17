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
