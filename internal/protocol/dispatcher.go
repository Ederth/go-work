package protocol

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-module/carbon"
	"github.com/streadway/amqp"
	"wz-car-worker/internal/biz"
	"wz-car-worker/internal/conf"
	"wz-car-worker/internal/data"
	"wz-car-worker/pkg/rabbitmq"
	"wz-car-worker/pkg/util/codec"
)

type Dispatcher struct {
	log    *log.Helper
	router map[string]Handler
}

func NewDispatcher(
	c *conf.Data,
	logger log.Logger,
	online *biz.OnlineCmdUseCase,
	stateReport *biz.StateReportCmdUseCase,
	syncTime *biz.SyncTimeCmdUseCase,
	deviceRestart *biz.DeviceRestartCmdUseCase) *Dispatcher {
	r := map[string]Handler{
		"01": online,
		"02": stateReport,
		"03": syncTime,
		"04": deviceRestart,
	}

	dispatcher := &Dispatcher{log: log.NewHelper(logger), router: r}

	name := data.ToWorker
	rabbitmq.Consume(logger, name, c.Mq[name], func(d amqp.Delivery) {
		if err := dispatcher.Dispatch(d.Body); err != nil {
			dispatcher.log.Error(err)
			if err = d.Reject(false); err != nil {
				dispatcher.log.Error(err)
			}
			return
		}

		if err := d.Ack(false); err != nil {
			dispatcher.log.Error(err)
			return
		}
	})

	return dispatcher
}

type Handler interface {
	Handle(ctx context.Context, d *biz.DecodeData) error
}

func (p *Dispatcher) Dispatch(b []byte) error {
	d, err := Decode(b)
	if err != nil {
		p.log.Info(string(b))
		return err
	}

	pd := &biz.DecodeData{
		FrameType:   d.FrameType,
		DeviceNum:   d.DeviceNum,
		DataDecrypt: d.DataDecrypt,
		Timestamp:   d.Timestamp,
		MsgId:       codec.Hex2Uint16(d.MsgId),
		Crc:         d.Crc,
	}
	p.log.Debugf("%s %+v", carbon.CreateFromTimestamp(d.Timestamp), pd)

	h, ok := p.router[pd.FrameType]
	if !ok {
		return nil
	}

	return h.Handle(context.Background(), pd)
}
