package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-module/carbon"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"wz-car-worker/internal/biz"
	"wz-car-worker/pkg/redis"
	"wz-car-worker/pkg/util/codec"
)

// deviceStateRepo 设备状态
type deviceStateRepo struct {
	data *Data
	log  *log.Helper
	*Mq
	onlineAtRepo  *OnlineAtRepo
	portStateRepo *PortStateRepo
}

func NewDeviceStateRepo(
	data *Data,
	logger log.Logger,
	mq *Mq,
	onlineAtRepo *OnlineAtRepo,
	portStateRepo *PortStateRepo,
) biz.DeviceStateRepo {
	return &deviceStateRepo{
		data:          data,
		log:           log.NewHelper(logger),
		Mq:            mq,
		onlineAtRepo:  onlineAtRepo,
		portStateRepo: portStateRepo,
	}
}

// Online 设备上线
func (repo *deviceStateRepo) Online(ctx context.Context, req *biz.OnlineCmdReq) error {
	res := repo.onlineAtRepo.HSet(ctx, "", req.DeviceNum, carbon.Now().ToDateTimeString())
	if res.Err() != nil {
		return errors.Wrap(res.Err(), "set online_at fail")
	}

	// todo 保存设备信息
	return nil
}

// Report 上报设备状态
func (repo *deviceStateRepo) Report(ctx context.Context, req *biz.StateReportCmdReq) error {
	// 保存设备状态
	state := &biz.DevicePortState{}
	if err := copier.Copy(state, req); err != nil {
		return err
	}
	res := repo.portStateRepo.HSet(ctx, req.DeviceNum+codec.Uint8ToBCD(req.PortNum), map[string]interface{}{
		"state":      state.State,
		"updated_at": carbon.Now().ToDateTimeString(),
	})
	if res.Err() != nil {
		return errors.Wrap(res.Err(), "set port state fail")
	}

	return nil
}

// Restart 设备重启
func (repo *deviceStateRepo) Restart(ctx context.Context, deviceNum string) error {
	cmd := &DeviceRestartCmdReq{
		DeviceNum: deviceNum,
		Type:      1,
	}

	return repo.Pub(&EncodeData{
		FrameType: "92",
		DeviceNum: deviceNum,
		Data:      cmd.Format(),
		MsgId:     2,
	})
}

// DeviceRestartCmdReq 设备重启命令
type DeviceRestartCmdReq struct {
	DeviceNum string // 设备编号
	Type      uint8  // 执行控制 0x01 立即执行 0x02 空闲执行
}

func (req *DeviceRestartCmdReq) Format() string {
	return req.DeviceNum + codec.Uint8ToHex(req.Type)
}

// SyncTime 同步时间
func (repo *deviceStateRepo) SyncTime(ctx context.Context, deviceNum string, ts int64) error {
	cmd := &SyncTimeCmdReq{
		DeviceNum: deviceNum,
		Time:      ts,
	}

	return repo.Pub(&EncodeData{
		FrameType: "56",
		DeviceNum: deviceNum,
		Data:      cmd.Format(),
		MsgId:     1,
	})
}

// SyncTimeCmdReq 同步时间命令
type SyncTimeCmdReq struct {
	DeviceNum string // 设备编号
	Time      int64
}

func (req *SyncTimeCmdReq) Format() string {
	return req.DeviceNum +
		codec.CP56Time2a(time.Unix(req.Time, 0))
}

// Refresh 刷新设备状态
func (repo *deviceStateRepo) Refresh(ctx context.Context, deviceNum string, portNum uint8) error {
	cmd := &GetStateCmdReq{
		DeviceNum: deviceNum,
		PortNum:   portNum,
	}

	return repo.Pub(&EncodeData{
		FrameType: "12",
		DeviceNum: deviceNum,
		Data:      cmd.Format(),
		MsgId:     1,
	})
}

// GetStateCmdReq 主动获取状态
type GetStateCmdReq struct {
	DeviceNum string
	PortNum   uint8
}

func (req *GetStateCmdReq) Format() string {
	return req.DeviceNum + codec.Uint8ToHex(req.PortNum)
}

// Get 获取设备端口状态
func (repo *deviceStateRepo) Get(ctx context.Context, deviceNum string, portNum uint8) (*biz.DevicePortState, error) {
	state := &biz.DevicePortState{}
	res := repo.portStateRepo.HGetAll(ctx, deviceNum+codec.Uint8ToBCD(portNum))
	err := res.Scan(state)

	return state, err
}

// OnlineAtRepo 设备上线时间 标示是否在线
type OnlineAtRepo struct {
	redis.Hash
}

func NewOnlineAtRepo(data *Data) *OnlineAtRepo {
	return &OnlineAtRepo{redis.Hash{
		Client:   data.rdb,
		RedisKey: "online_at",
	}}
}

// PortStateRepo 端口状态
type PortStateRepo struct {
	redis.Hash
}

func NewPortStateRepo(data *Data) *PortStateRepo {
	return &PortStateRepo{redis.Hash{
		Client:   data.rdb,
		RedisKey: "port_state",
	}}
}
