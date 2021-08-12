package rabbitmq

import (
	"github.com/streadway/amqp"
	"wz-car-worker/internal/conf"
)

type Exchange struct {
	name       string
	kind       string
	durable    bool
	autoDelete bool
	internal   bool
	noWait     bool
	args       amqp.Table
}

func newExchange(c *conf.Data_RabbitMq_Exchange) *Exchange {
	return &Exchange{
		name:       c.Name,
		kind:       c.Kind,
		durable:    true,
		autoDelete: false,
		internal:   false,
		noWait:     false,
		args:       nil,
	}
}

type Queue struct {
	name       string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
	args       amqp.Table
}

func newQueue(c *conf.Data_RabbitMq_Queue) *Queue {
	return &Queue{
		name:       c.Name,
		durable:    true,
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
		args:       nil,
	}
}
