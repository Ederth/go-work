package data

import (
	"encoding/json"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-module/carbon"
	"github.com/segmentio/ksuid"
	"github.com/streadway/amqp"
	"wz-car-worker/pkg/rabbitmq"
	"wz-car-worker/pkg/util/codec"
)

const (
	ToWorker  string = "to_worker"
	ToGateway string = "to_gateway"
)

type EncodeData struct {
	FrameType string
	DeviceNum string
	Data      string
	MsgId     uint16
}

type encodeData struct {
	FrameType   string `json:"frame_type"`
	DeviceNum   string `json:"device_num"`
	Data        string `json:"data"`
	MsgId       string `json:"msg_id"`
	IsEncrypted bool   `json:"is_encrypted"`
}

func Encode(data *encodeData) ([]byte, error) {
	return json.Marshal(data)
}

type Mq struct {
	mq   *rabbitmq.RabbitMq
	log  *log.Helper
	name string
}

func NewMq(data *Data, logger log.Logger) *Mq {
	return &Mq{
		mq:   data.mq,
		log:  log.NewHelper(logger),
		name: ToGateway,
	}
}

func (t *Mq) Pub(data *EncodeData) error {
	d := &encodeData{
		FrameType:   data.FrameType,
		DeviceNum:   data.DeviceNum,
		Data:        data.Data,
		MsgId:       codec.Uint16ToHex(data.MsgId),
		IsEncrypted: false,
	}
	body, err := Encode(d)
	if err != nil {
		return err
	}

	con, err := t.mq.Connection(t.name)
	if err != nil {
		return err
	}
	con.Pub(amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		MessageId:    ksuid.New().String(),
		Timestamp:    carbon.Now().Time,
		Body:         body,
	})
	out, _ := json.Marshal(d)
	t.log.Info(t.name, string(out))

	return nil
}
