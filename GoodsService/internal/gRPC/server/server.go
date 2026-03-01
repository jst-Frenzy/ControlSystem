package server

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/GoodService"
	gen "github.com/jst-Frenzy/ControlSystem/protobuf/gen/goods"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"strconv"
)

type Deps struct {
	GoodsService GoodService.GoodService
	Logger       *logrus.Logger
}

type Server struct {
	gen.UnimplementedGoodsServiceServer
	srv          *grpc.Server
	goodsService GoodService.GoodService
	logger       *logrus.Entry
}

func NewGRPCServer(d Deps) *Server {
	logrusEntry := logrus.NewEntry(d.Logger)

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
		goodsService: d.GoodsService,
		logger:       logrusEntry,
	}
}

func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	gen.RegisterGoodsServiceServer(s.srv, s)

	return s.srv.Serve(lis)
}

func (s *Server) Stop() {
	if s.srv != nil {
		s.srv.GracefulStop()
	}
}

func (s *Server) GetItemQuantityAndPrice(ctx context.Context, req *gen.ItemQuantityAndPriceRequest) (*gen.ItemQuantityAndPriceResponse, error) {
	if req.GetItemId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "item id is required")
	}

	info, err := s.goodsService.GetItemInfoForCart(req.GetItemId())
	if err != nil {
		return &gen.ItemQuantityAndPriceResponse{Valid: false}, err
	}

	return &gen.ItemQuantityAndPriceResponse{
		Valid:    true,
		Quantity: strconv.Itoa(info.Quantity),
		Price:    strconv.FormatFloat(info.Price, 'f', 2, 64),
	}, nil
}
