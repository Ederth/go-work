package data

import (
	kLog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"wz-car-worker/internal/conf"
	"wz-car-worker/pkg/rabbitmq"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewMq,
	NewGreeterRepo,
	NewDeviceStateRepo,
	NewOnlineAtRepo,
	NewPortStateRepo,
)

// Data .
type Data struct {
	rdb *redis.Client
	mq  *rabbitmq.RabbitMq
}

// NewData .
func NewData(c *conf.Data, logger kLog.Logger) (*Data, func(), error) {
	log := kLog.NewHelper(logger)

	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Db),
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
	})
	rdb.AddHook(redisotel.TracingHook{})

	mq := rabbitmq.NewClient(logger, c.Mq)

	d := &Data{
		rdb: rdb,
		mq:  mq,
	}

	return d, func() {
		log.Info("closing the data resources")
		if err := d.rdb.Close(); err != nil {
			log.Error(err)
		}
	}, nil
}
