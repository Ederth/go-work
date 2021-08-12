package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeData_Validate(t *testing.T) {
	data := DecodeData{
		FrameType:   "0a",
		DeviceNum:   "12000000000001",
		DataDecrypt: "120000000000010102",
		Timestamp:   1624256661,
		MsgId:       "0001",
		Crc:         1234,
	}

	t.Run("FrameType", func(t *testing.T) {
		input := data
		s := []string{"00", "01", "0a", "10", "ff"}
		for _, v := range s {
			input.FrameType = v
			err := input.Validate()
			assert.Nil(t, err, v)
		}

		s = []string{"", "123", "xx"}
		for _, v := range s {
			input.FrameType = v
			err := input.Validate()
			assert.Error(t, err, v)
		}
	})

	t.Run("DeviceNum", func(t *testing.T) {
		input := data
		s := []string{"", "1200000000000a", "123", "1511111111111111111"}
		for _, v := range s {
			input.DeviceNum = v
			err := input.Validate()
			assert.Error(t, err, v)
		}
	})

	t.Run("MsgId", func(t *testing.T) {
		input := data
		s := []string{"0000", "000a", "ffff", "1001"}
		for _, v := range s {
			input.MsgId = v
			err := input.Validate()
			assert.Nil(t, err, v)
		}

		s = []string{"", "123", "12345", "000x"}
		for _, v := range s {
			input.MsgId = v
			err := input.Validate()
			assert.Error(t, err, v)
		}
	})
}
