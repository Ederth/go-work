package rabbitmq

import (
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/streadway/amqp"
	"wz-car-worker/internal/conf"
)

var (
	errConnectionNotFound = errors.New("rabbitmq connection not found")
)

type RabbitMq struct {
	list map[string]*connection
}

func NewClient(logger log.Logger, c map[string]*conf.Data_RabbitMq) *RabbitMq {
	mq := &RabbitMq{
		list: make(map[string]*connection),
	}
	for k, v := range c {
		exchange := newExchange(v.Exchange)
		queue := newQueue(v.Queue)
		con := &connection{
			name:          k,
			addr:          v.Addr,
			publishCh:     make(chan amqp.Publishing),
			publishNum:    int(v.PublishNum),
			consumeNum:    int(v.ConsumerNum),
			prefetchCount: int(v.PrefetchCount),
			exchangeConf:  exchange,
			queueConf:     queue,
			log:           log.NewHelper(logger),
		}

		con.newPubWorker()
		mq.list[k] = con
	}

	return mq
}

func Consume(logger log.Logger, name string, conf *conf.Data_RabbitMq, f func(d amqp.Delivery)) {
	exchange := newExchange(conf.Exchange)
	queue := newQueue(conf.Queue)
	con := &connection{
		name:          name,
		addr:          conf.Addr,
		publishCh:     make(chan amqp.Publishing),
		publishNum:    int(conf.PublishNum),
		consumeNum:    int(conf.ConsumerNum),
		prefetchCount: int(conf.PrefetchCount),
		exchangeConf:  exchange,
		queueConf:     queue,
		log:           log.NewHelper(logger),
	}

	con.newConsumeWorker(f)
}

func (mq RabbitMq) Connection(name string) (*connection, error) {
	con, ok := mq.list[name]
	if !ok {
		return nil, errConnectionNotFound
	}

	return con, nil
}
