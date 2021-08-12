package service

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	pb "wz-car-worker/api/worker/v1"
	"wz-car-worker/internal/biz"
)

type DeviceStateService struct {
	pb.DeviceStateHTTPClientImpl
	log *log.Helper
	uc  *biz.DeviceStateUseCase
}

func NewDeviceStateService(logger log.Logger, uc *biz.DeviceStateUseCase) *DeviceStateService {
	return &DeviceStateService{
		log: log.NewHelper(logger),
		uc:  uc,
	}
}

// Restart 设备重启
func (s *DeviceStateService) Restart(ctx context.Context, req *pb.DeviceRestartRequest) (*pb.DeviceRestartReply, error) {
	for _, deviceNum := range req.DeviceNums {
		_ = s.uc.Restart(ctx, deviceNum)
	}
	return &pb.DeviceRestartReply{Data: nil}, nil
}

// SyncTime 同步时间
func (s *DeviceStateService) SyncTime(ctx context.Context, req *pb.SyncTimeRequest) (*pb.SyncTimeReply, error) {
	ts := req.GetTime()
	if req.Time == nil {
		ts = time.Now().Unix()
	}
	for _, deviceNum := range req.DeviceNums {
		_ = s.uc.SyncTime(ctx, deviceNum, ts)
	}
	return &pb.SyncTimeReply{}, nil
}

// Refresh 刷新状态
func (s *DeviceStateService) Refresh(ctx context.Context, req *pb.DeviceStateRefreshRequest) (*pb.DeviceStateRefreshReply, error) {
	err := s.uc.Refresh(ctx, req.DeviceNum, uint8(req.PortNum))

	return &pb.DeviceStateRefreshReply{}, err
}

// Get 获取设备状态
func (s *DeviceStateService) Get(ctx context.Context, req *pb.DeviceStateGetRequest) (*pb.DeviceStateGetReply, error) {
	state, err := s.uc.Get(ctx, req.DeviceNum, uint8(req.PortNum))
	if err != nil {
		return nil, err
	}

	data := &pb.DeviceStateGetReply_DeviceState{}
	if err = copier.Copy(data, state); err != nil {
		return nil, err
	}
	data.PayUsed /= 100 // 转为2位小数
	reply := &pb.DeviceStateGetReply{Data: data}

	return reply, nil
}
