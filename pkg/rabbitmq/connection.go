package rabbitmq

import (
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/streadway/amqp"
)

type connection struct {
	name          string
	addr          string
	publishCh     chan amqp.Publishing
	consumeCh     <-chan amqp.Delivery
	publishNum    int
	consumeNum    int
	prefetchCount int
	exchangeConf  *Exchange
	queueConf     *Queue
	con           []*amqp.Connection
	channel       []*amqp.Channel
	queue         []amqp.Queue
	log           *log.Helper
	lock          sync.RWMutex
}

func (c *connection) newPubWorker() {
	c.con = make([]*amqp.Connection, c.publishNum)
	c.channel = make([]*amqp.Channel, c.publishNum)
	c.queue = make([]amqp.Queue, c.publishNum)
	for i := 0; i < c.publishNum; i++ {
		go c.loop(i, func(id int) error {
			ch, err := c.dail(id, "send")
			if err != nil {
				return err
			}

			for v := range c.publishCh {
				if err = ch.Publish(c.exchangeConf.name, c.queueConf.name, false, false, v); err != nil {
					return err
				}
			}

			return nil
		})
	}
}

func (c *connection) newConsumeWorker(f func(d amqp.Delivery)) {
	c.con = make([]*amqp.Connection, c.consumeNum)
	c.channel = make([]*amqp.Channel, c.consumeNum)
	c.queue = make([]amqp.Queue, c.consumeNum)
	for i := 0; i < c.consumeNum; i++ {
		go c.loop(i, func(id int) error {
			ch, err := c.dail(id, "recv")
			if err != nil {
				return err
			}

			if err = ch.Qos(c.prefetchCount, 0, false); err != nil {
				return err
			}

			c.consumeCh, err = ch.Consume(c.queueConf.name, "", false, false, false, false, nil)
			if err != nil {
				return err
			}

			for d := range c.consumeCh {
				go f(d)
			}

			return nil
		})
	}
}

func (c *connection) loop(id int, f func(id int) error) {
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					c.log.Error("panic: ", r)
				}
				c.Close(id)
			}()

			if err := f(id); err != nil {
				c.log.Errorf("id:%d %s", id, err)
			}
		}()

		time.Sleep(time.Second)
		c.log.Errorf("id:%d try to reconnect amqp", id)
	}
}

func (c *connection) dail(id int, t string) (*amqp.Channel, error) {
	con, err := amqp.DialConfig("amqp://"+c.addr, amqp.Config{
		Properties: amqp.Table{
			"connection_name": t + "." + c.name,
		}})
	if err != nil {
		return nil, err
	}
	c.lock.Lock()
	c.con[id] = con
	c.lock.Unlock()

	ch, err := con.Channel()
	if err != nil {
		return nil, err
	}
	c.lock.Lock()
	c.channel[id] = ch
	c.lock.Unlock()

	if err = c.exchangeDeclare(ch); err != nil {
		return nil, err
	}

	q, err := c.queueDeclare(ch)
	if err != nil {
		return nil, err
	}
	c.lock.Lock()
	c.queue[id] = q
	c.lock.Unlock()

	err = ch.QueueBind(c.queueConf.name, c.queueConf.name, c.exchangeConf.name, false, nil)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (c *connection) exchangeDeclare(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		c.exchangeConf.name,
		c.exchangeConf.kind,
		c.exchangeConf.durable,
		c.exchangeConf.autoDelete,
		c.exchangeConf.internal,
		c.exchangeConf.noWait,
		c.exchangeConf.args,
	)
}

func (c *connection) queueDeclare(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		c.queueConf.name,
		c.queueConf.durable,
		c.queueConf.autoDelete,
		c.queueConf.exclusive,
		c.queueConf.noWait,
		c.queueConf.args,
	)
}

func (c *connection) Pub(v amqp.Publishing) {
	c.publishCh <- v
}

func (c *connection) Close(id int) {
	c.lock.RLock()
	channel := c.channel[id]
	if channel != nil {
		_ = channel.Close()
	}

	con := c.con[id]
	if con != nil {
		_ = con.Close()
	}
	c.lock.RUnlock()
}
