package biz

import (
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewGreeterUsecase,
	NewOnlineCmdUseCase,
	NewStateReportCmdUseCase,
	NewDeviceStateUseCase,
	NewSyncTimeCmdUseCase,
	NewDeviceRestartCmdUseCase,
)

type DecodeData struct {
	FrameType   string
	DeviceNum   string
	DataDecrypt string
	Timestamp   int64
	MsgId       uint16
	Crc         uint16
}
