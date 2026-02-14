package gRPC

import (
	"context"
	"errors"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jst-Frenzy/ControlSystem/AuthService/internal/AuthService"
	gen "github.com/jst-Frenzy/ControlSystem/protobuf/gen"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"strconv"
)

type Deps struct {
	Logger      *logrus.Logger
	AuthService AuthService.AuthService
}

type Server struct {
	gen.UnimplementedAuthServiceServer
	srv         *grpc.Server
	authService AuthService.AuthService
	logger      *logrus.Entry
}

func NewGRPCServer(deps Deps) *Server {
	logrusEntry := logrus.NewEntry(deps.Logger)

	return &Server{
		srv: grpc.NewServer(
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
				grpc_logrus.StreamServerInterceptor(logrusEntry),
				grpc_recovery.StreamServerInterceptor(),
			)),
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_logrus.UnaryServerInterceptor(logrusEntry),
				grpc_recovery.UnaryServerInterceptor(),
			)),
		),
		authService: deps.AuthService,
		logger:      logrusEntry,
	}
}

func (s *Server) StartGRPC(port int) error {
	addr := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	gen.RegisterAuthServiceServer(s.srv, s)

	return s.srv.Serve(lis)
}

func (s *Server) Stop() {
	if s.srv != nil {
		s.srv.GracefulStop()
	}
}

func (s *Server) ValidateToken(ctx context.Context, req *gen.ValidateTokenRequest) (*gen.ValidateTokenResponse, error) {
	if req.GetAccessToken() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "access_token is required")
	}

	info, err := s.authService.ParseToken(req.GetAccessToken())
	if err != nil {
		return &gen.ValidateTokenResponse{Valid: false}, errors.New("can't parse token")
	}

	return &gen.ValidateTokenResponse{
		Valid:    true,
		UserId:   strconv.Itoa(info.ID),
		Role:     info.Role,
		UserName: info.UserName,
	}, nil
}
