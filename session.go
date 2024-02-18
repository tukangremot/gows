package gochat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

const (
	SessionDriverRedis        = "redis"
	SessionUserChannelDataKey = "uc:%s:%s"
	SessionUserGroupDataKey   = "ug:%s:%s"
)

type (
	Session struct {
		driver string
		conn   interface{}
	}
)

func NewSession(driver string, conn interface{}) *Session {
	return &Session{
		driver: driver,
		conn:   conn,
	}
}

func (session *Session) registerUserChannel(ctx context.Context, channel *Channel, user *User) error {
	switch session.driver {
	case SessionDriverRedis:
		redisClient := session.conn.(*redis.Client)

		userData, err := json.Marshal(user)
		if err != nil {
			return err
		}

		err = redisClient.Set(ctx, fmt.Sprintf(SessionUserChannelDataKey, channel.ID, user.ID), userData, 0).Err()
		if err != nil {
			return err
		}

	}

	return nil
}

func (session *Session) unregisterUserChannel(ctx context.Context, channel *Channel, user *User) error {
	switch session.driver {
	case SessionDriverRedis:
		redisClient := session.conn.(*redis.Client)

		err := redisClient.Del(ctx, fmt.Sprintf(SessionUserChannelDataKey, channel.ID, user.ID)).Err()
		if err != nil {
			return err
		}

	}

	return nil
}

func (session *Session) findUserChannelByID(ctx context.Context, channel *Channel, userID string) (*User, error) {
	switch session.driver {
	case SessionDriverRedis:
		redisClient := session.conn.(*redis.Client)

		data, err := redisClient.Get(ctx, fmt.Sprintf(SessionUserChannelDataKey, channel.ID, userID)).Result()
		if err != nil {
			return nil, err
		}

		var user *User
		err = json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}

		user.onDifferentServer = true

		return user, nil
	}

	return nil, nil
}

func (session *Session) registerUserGroup(ctx context.Context, group *Group, user *User) error {
	switch session.driver {
	case SessionDriverRedis:
		redisClient := session.conn.(*redis.Client)

		userData, err := json.Marshal(user)
		if err != nil {
			return err
		}

		err = redisClient.Set(ctx, fmt.Sprintf(SessionUserGroupDataKey, group.ID, user.ID), userData, 0).Err()
		if err != nil {
			return err
		}

	}

	return nil
}

func (session *Session) unregisterUserGroup(ctx context.Context, group *Group, user *User) error {
	switch session.driver {
	case SessionDriverRedis:
		redisClient := session.conn.(*redis.Client)
		err := redisClient.Del(ctx, fmt.Sprintf(SessionUserGroupDataKey, group.ID, user.ID)).Err()
		if err != nil {
			return err
		}

	}

	return nil
}

func (session *Session) getUsersByGroup(ctx context.Context, group *Group) (map[string]*User, error) {
	users := make(map[string]*User)

	switch session.driver {
	case SessionDriverRedis:
		redisClient := session.conn.(*redis.Client)
		iter := redisClient.Scan(ctx, 0, fmt.Sprintf(SessionUserGroupDataKey, group.ID, "*"), 0).Iterator()
		for iter.Next(ctx) {
			data, err := redisClient.Get(ctx, iter.Val()).Result()
			if err != nil {
				log.Println(err)

				return nil, err
			}

			var user *User
			err = json.Unmarshal([]byte(data), &user)
			if err != nil {
				log.Println(err)

				return nil, err
			}

			users[user.ID] = user
		}
		if err := iter.Err(); err != nil {
			log.Println(err)
		}

	}

	return users, nil
}
