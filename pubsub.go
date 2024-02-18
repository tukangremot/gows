package gochat

const (
	PubSubDriverRedis = "redis"
)

type (
	PubSub struct {
		driver string
		conn   interface{}
	}
)

func NewPubSub(driver string, conn interface{}) *PubSub {
	return &PubSub{
		driver: driver,
		conn:   conn,
	}
}
