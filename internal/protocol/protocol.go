package protocol

import (
	"encoding/json"

	"github.com/go-playground/validator"
	"github.com/google/wire"
)

// ProviderSet is protocol providers.
var ProviderSet = wire.NewSet(NewDispatcher)

type DecodeData struct {
	FrameType   string `json:"frame_type" validate:"required,hexadecimal,len=2"`
	DeviceNum   string `json:"device_num" validate:"required,numeric,len=14"`
	DataDecrypt string `json:"data_decrypt" validate:"required"`
	Timestamp   int64  `json:"timestamp" validate:"required"`
	MsgId       string `json:"msg_id" validate:"required,hexadecimal,len=4"`
	Crc         uint16 `json:"crc"`
}

func (d *DecodeData) Validate() error {
	v := validator.New()
	return v.Struct(d)
}

func Decode(b []byte) (*DecodeData, error) {
	var pd = &DecodeData{}
	if err := json.Unmarshal(b, pd); err != nil {
		return nil, err
	}

	err := pd.Validate()

	return pd, err
}
