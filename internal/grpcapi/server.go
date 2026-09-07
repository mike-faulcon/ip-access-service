package grpcapi

import (
	"context"
	"net"

	accessv1 "github.com/mike-faulcon/ip-access-service/gen/access/v1"
	"github.com/mike-faulcon/ip-access-service/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	accessv1.UnimplementedAccessServiceServer

	service *service.AccessService
}

func NewServer(svc *service.AccessService) *Server {
	return &Server{
		service: svc,
	}
}

func (s *Server) CheckAccess(
	ctx context.Context,
	req *accessv1.CheckAccessRequest,
) (*accessv1.CheckAccessResponse, error) {
	ip := net.ParseIP(req.GetIp())
	if ip == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid IP address")
	}

	accessRequst := service.AccessRequest{
		IP:               ip,
		AllowedCountries: req.GetAllowedCountries(),
	}

	result, err := s.service.Check(
		ctx,
		accessRequst,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check access")
	}

	return &accessv1.CheckAccessResponse{
		Allowed: result.Allowed,
		Country: result.CountryISO,
	}, nil
}
