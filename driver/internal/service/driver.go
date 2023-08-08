package service

import (
	"context"
	"driver/internal/biz"

	pb "driver/api/driver"
)

type DriverService struct {
	pb.UnimplementedDriverServer
	bz *biz.DriverBiz
}

func NewDriverService(bz *biz.DriverBiz) *DriverService {
	return &DriverService{
		bz: bz,
	}
}

func (s *DriverService) GetVerifiyCode(ctx context.Context, req *pb.GetVerifyCodeReq) (*pb.GetVerifyCodeResp, error) {
	return &pb.GetVerifyCodeResp{}, nil
}
func (s *DriverService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	return &pb.LoginResp{}, nil
}
func (s *DriverService) Loginout(ctx context.Context, req *pb.LoginoutReq) (*pb.LoginoutResp, error) {
	return &pb.LoginoutResp{}, nil
}
