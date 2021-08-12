package biz

import (
	"context"
	"encoding/hex"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"wz-car-worker/pkg/util/codec"
)

// DevicePortState 设备端口状态
type DevicePortState struct {
	State     uint8  `redis:"state,omitempty"`      // 状态 0离线 1故障 2空闲
	UpdatedAt string `redis:"updated_at,omitempty"` // 更新时间
}

type DeviceStateRepo interface {
	Online(ctx context.Context, req *OnlineCmdReq) error
	Report(ctx context.Context, req *StateReportCmdReq) error
	Restart(ctx context.Context, deviceNum string) error
	SyncTime(ctx context.Context, deviceNum string, ts int64) error
	Refresh(ctx context.Context, deviceNum string, portNum uint8) error
	Get(ctx context.Context, deviceNum string, portNum uint8) (*DevicePortState, error)
}

// DeviceStateUseCase 设备状态
type DeviceStateUseCase struct {
	log  *log.Helper
	repo DeviceStateRepo
}

func NewDeviceStateUseCase(logger log.Logger, repo DeviceStateRepo) *DeviceStateUseCase {
	return &DeviceStateUseCase{
		log:  log.NewHelper(logger),
		repo: repo,
	}
}

// Restart 设备重启
func (uc *DeviceStateUseCase) Restart(ctx context.Context, deviceNum string) error {
	return uc.repo.Restart(ctx, deviceNum)
}

// SyncTime 同步时间
func (uc *DeviceStateUseCase) SyncTime(ctx context.Context, deviceNum string, ts int64) error {
	return uc.repo.SyncTime(ctx, deviceNum, ts)
}

// Refresh 刷新设备状态
func (uc *DeviceStateUseCase) Refresh(ctx context.Context, deviceNum string, portNum uint8) error {
	return uc.repo.Refresh(ctx, deviceNum, portNum)
}

// Get 获取设备状态
func (uc *DeviceStateUseCase) Get(ctx context.Context, deviceNum string, portNum uint8) (*DevicePortState, error) {
	return uc.repo.Get(ctx, deviceNum, portNum)
}

type OnlineCmdReq struct {
	DeviceNum  string `validate:"required,numeric,len=14"` // 设备编号
	DeviceType uint8  `validate:"required,max=1"`          // 设备类型
}

func (r *OnlineCmdReq) Validate() error {
	v := validator.New()
	return v.Struct(r)
}

// OnlineCmdUseCase 设备上线命令
type OnlineCmdUseCase struct {
	log  *log.Helper
	repo DeviceStateRepo
}

func NewOnlineCmdUseCase(logger log.Logger, repo DeviceStateRepo) *OnlineCmdUseCase {
	return &OnlineCmdUseCase{log: log.NewHelper(logger), repo: repo}
}

func (uc *OnlineCmdUseCase) Parse(d *DecodeData) (*OnlineCmdReq, error) {
	v, err := hex.DecodeString(d.DataDecrypt)
	if err != nil {
		return nil, errors.Wrapf(err, "DecodeData: %+v", d)
	}

	return &OnlineCmdReq{
		DeviceNum:  d.DataDecrypt[:14],
		DeviceType: v[7],
	}, nil
}

func (uc *OnlineCmdUseCase) Handle(ctx context.Context, d *DecodeData) error {
	req, err := uc.Parse(d)
	uc.log.Infof("%+v", req)
	// todo 测试数据解析 网络类型 运营商
	if err != nil {
		return err
	}

	if err = req.Validate(); err != nil {
		return err
	}
	// todo 上线失败处理
	return uc.repo.Online(ctx, req)
}

type StateReportCmdReq struct {
	DeviceNum string `validate:"required,numeric,len=14"` // 设备编号
	PortNum   uint8  // 端口号
	State     uint8  `validate:"max=5"` // 状态 0x00 离线 0x01 故障 0x02 空闲
}

func (r *StateReportCmdReq) Validate() error {
	// todo 设备编号不一致
	// todo 设备是否上线
	v := validator.New()
	return v.Struct(r)
}

// StateReportCmdUseCase 实时监测数据上报命令
type StateReportCmdUseCase struct {
	log  *log.Helper
	repo DeviceStateRepo
}

func NewStateReportCmdUseCase(logger log.Logger, repo DeviceStateRepo) *StateReportCmdUseCase {
	return &StateReportCmdUseCase{log: log.NewHelper(logger), repo: repo}
}

func (uc *StateReportCmdUseCase) Parse(d *DecodeData) (*StateReportCmdReq, error) {
	_, err := hex.DecodeString(d.DataDecrypt)
	if err != nil {
		return nil, errors.Wrapf(err, "DecodeData: %+v", d)
	}

	req := &StateReportCmdReq{
		DeviceNum: "",
		PortNum:   1,
		State:     0,
	}

	return req, nil
}

func (uc *StateReportCmdUseCase) Handle(ctx context.Context, d *DecodeData) error {
	req, err := uc.Parse(d)
	uc.log.Infof("%+v", req)
	if err != nil {
		return err
	}

	if err = req.Validate(); err != nil {
		return err
	}

	return uc.repo.Report(ctx, req)
}

type SyncTimeCmdResp struct {
	DeviceNum string `validate:"required,numeric,len=14"` // 设备编号
	Now       int64  // 当前时间 CP65Time2a 格式
}

func (r *SyncTimeCmdResp) Validate() error {
	// todo 设备编号不一致
	// todo 设备是否上线
	v := validator.New()
	return v.Struct(r)
}

// SyncTimeCmdUseCase 同步时间命令
type SyncTimeCmdUseCase struct {
	log *log.Helper
}

func NewSyncTimeCmdUseCase(logger log.Logger) *SyncTimeCmdUseCase {
	return &SyncTimeCmdUseCase{log: log.NewHelper(logger)}
}

func (uc *SyncTimeCmdUseCase) Parse(d *DecodeData) *SyncTimeCmdResp {
	return &SyncTimeCmdResp{
		DeviceNum: d.DataDecrypt[:14],
		Now:       codec.ParseCP56Time2a(d.DataDecrypt[14:28]).Unix(),
	}
}

func (uc *SyncTimeCmdUseCase) Handle(ctx context.Context, d *DecodeData) error {
	resp := uc.Parse(d)
	uc.log.Infof("%+v", resp)

	return resp.Validate()
}

type DeviceRestartResp struct {
	DeviceNum string `validate:"required,numeric,len=14"` // 设备编号
	State     uint8  `validate:"max=1"`                   // 启动结果 0x00 失败 0x01 成功
}

func (r *DeviceRestartResp) Validate() error {
	// todo 设备编号不一致
	// todo 设备是否上线
	v := validator.New()
	return v.Struct(r)
}

// DeviceRestartCmdUseCase 设备重启命令
type DeviceRestartCmdUseCase struct {
	log *log.Helper
}

func NewDeviceRestartCmdUseCase(logger log.Logger) *DeviceRestartCmdUseCase {
	return &DeviceRestartCmdUseCase{log: log.NewHelper(logger)}
}

func (uc *DeviceRestartCmdUseCase) Parse(d *DecodeData) *DeviceRestartResp {
	return &DeviceRestartResp{
		DeviceNum: d.DataDecrypt[:14],
		State:     codec.Hex2Uint8(d.DataDecrypt[14:16]),
	}
}

func (uc *DeviceRestartCmdUseCase) Handle(ctx context.Context, d *DecodeData) error {
	resp := uc.Parse(d)
	uc.log.Infof("%+v", resp)

	return resp.Validate()
}
