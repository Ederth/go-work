package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"wz-car-worker/internal/biz"
)

type greeterRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *greeterRepo) CreateGreeter(ctx context.Context, g *biz.Greeter) error {
	return nil
}

func (r *greeterRepo) UpdateGreeter(ctx context.Context, g *biz.Greeter) error {
	return nil
}

func (r *greeterRepo) GetName(ctx context.Context, g *biz.Greeter) error {
	r.data.rdb.Set(ctx, "name", "Ederth", 0)

	return nil
}
