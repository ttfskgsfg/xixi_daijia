package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"map/internal/biz"

	pb "map/api/mapService"
)

type MapServiceService struct {
	pb.UnimplementedMapServiceServer

	msbiz *biz.MapServiceBiz //将MapServiceBiz与mapService建立联系
}

func NewMapServiceService(msbiz *biz.MapServiceBiz) *MapServiceService {
	return &MapServiceService{
		msbiz: msbiz,
	}
}

func (s *MapServiceService) GetDrivingInfo(ctx context.Context, req *pb.GetDrivingInfoReq) (*pb.GetDrivingReply, error) {
	distance, duration, err := s.msbiz.GetDriverInfo(req.Origin, req.Destionation)
	if err != nil {
		return nil, errors.New(200, "LBS_ERROR", "lbs api error")
	}
	return &pb.GetDrivingReply{
		Origin:       req.Origin,
		Destionation: req.Destionation,
		Distance:     distance,
		Duration:     duration,
	}, nil
}
