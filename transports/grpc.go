package transport

import (
	"context"
	"reflect"

	"github.com/go-kit/kit/log"
	gt "github.com/go-kit/kit/transport/grpc"
	"github.com/junereycasuga/gokit-grpc-demo/endpoints"
	"github.com/junereycasuga/gokit-grpc-demo/pb"
	"github.com/junereycasuga/gokit-grpc-demo/repository"
)

type gRPCServer struct {
	add      gt.Handler
	subtract gt.Handler
	multiply gt.Handler
	divide   gt.Handler
	cda      gt.Handler
}

// NewGRPCServer initializes a new gRPC server
func NewGRPCServer(endpoints endpoints.Endpoints, logger log.Logger) pb.MathServiceServer {
	return &gRPCServer{
		add: gt.NewServer(
			endpoints.Add,
			decodeMathRequest,
			encodeMathResponse,
		),
		subtract: gt.NewServer(
			endpoints.Subtract,
			decodeMathRequest,
			encodeMathResponse,
		),
		multiply: gt.NewServer(
			endpoints.Multiply,
			decodeMathRequest,
			encodeMathResponse,
		),
		divide: gt.NewServer(
			endpoints.Divide,
			decodeMathRequest,
			encodeMathResponse,
		),
		cda: gt.NewServer(
			endpoints.Cda,
			decodeMathRequest,
			encodeOtherResponse,
		),
	}
}

func (s *gRPCServer) Add(ctx context.Context, req *pb.MathRequest) (*pb.MathResponse, error) {
	_, resp, err := s.add.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.MathResponse), nil
}

func (s *gRPCServer) Subtract(ctx context.Context, req *pb.MathRequest) (*pb.MathResponse, error) {
	_, resp, err := s.subtract.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.MathResponse), nil
}

func (s *gRPCServer) Multiply(ctx context.Context, req *pb.MathRequest) (*pb.MathResponse, error) {
	_, resp, err := s.multiply.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.MathResponse), nil
}

func (s *gRPCServer) Cda(ctx context.Context, req *pb.MathRequest) (*pb.OtherResponse, error) {
	_, resp, err := s.cda.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.OtherResponse), nil
}

func (s *gRPCServer) Divide(ctx context.Context, req *pb.MathRequest) (*pb.MathResponse, error) {
	_, resp, err := s.divide.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.MathResponse), nil
}

func decodeMathRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.MathRequest)
	return endpoints.MathReq{NumA: req.NumA, NumB: req.NumB}, nil
}

func encodeMathResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.MathResp)
	return &pb.MathResponse{Result: resp.Result}, nil
}

func decodeOtherRequest(_ context.Context, request interface{}) (interface{}, error) {
	return endpoints.OtherReq{NumA: "dummy"}, nil
}

func encodeOtherResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.OtherResp)
	// var jsonData, _ = json.Marshal(resp.Result)
	// var buf bytes.Buffer
	// enc := gob.NewEncoder(&buf)
	// err := enc.Encode(resp.Result)
	// if err != nil {
	// 	return nil, err
	// }
	var transforms []*pb.DataDb
	s := reflect.ValueOf(resp.Result)

	for i := 0; i < s.Len(); i++ {
		var row pb.DataDb
		curentData := s.Index(i).Interface().(repository.Transaction)
		row.Id = curentData.Id
		row.OfferId = curentData.Offer_id
		row.QuoteId = curentData.Quote_id
		transforms = append(transforms, &row)
	}
	// json.Unmarshal(buf.Bytes(), &transforms)
	return &pb.OtherResponse{Data: transforms}, nil
}
