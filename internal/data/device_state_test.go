package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeviceRestartCmdReq_Format(t *testing.T) {
	cmd := &DeviceRestartCmdReq{
		DeviceNum: "18000000000001",
		Type:      1,
	}

	data := cmd.Format()
	assert.Equal(t, "1800000000000101", data)
}

func TestSyncTimeCmdReq_Format(t *testing.T) {
	cmd := &SyncTimeCmdReq{
		DeviceNum: "18000000000001",
		Time:      1625723643,
	}

	data := cmd.Format()
	assert.Equal(t, "18000000000001b80b360d880715", data)
}

func TestGetStateCmdReq_Format(t *testing.T) {
	cmd := &GetStateCmdReq{
		DeviceNum: "18000000000001",
		PortNum:   1,
	}

	data := cmd.Format()
	assert.Equal(t, "1800000000000101", data)
}
