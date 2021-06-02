package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/junereycasuga/gokit-grpc-demo/endpoints"
	"github.com/junereycasuga/gokit-grpc-demo/pb"
	"github.com/junereycasuga/gokit-grpc-demo/repository/postgres"
	"github.com/junereycasuga/gokit-grpc-demo/service"
	transport "github.com/junereycasuga/gokit-grpc-demo/transports"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	var logger log.Logger
	logger = log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	dbConfig := postgres.ConnParam{"localhost", "5432", "demo_cpq2", "postgres", "123", "none", 5, 5, 6000000000}
	dbConn, _ := postgres.NewPostgresSql(&dbConfig)

	addservice := service.NewService(logger, dbConn)
	addendpoint := endpoints.MakeEndpoints(addservice)
	grpcServer := transport.NewGRPCServer(addendpoint, logger)

	//db config

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	grpcListener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log("during", "Listen", "err", err)
		os.Exit(1)
	}

	go func() {
		baseServer := grpc.NewServer()
		pb.RegisterMathServiceServer(baseServer, grpcServer)
		level.Info(logger).Log("msg", "Server started successfully 🚀")
		baseServer.Serve(grpcListener)
	}()

	level.Error(logger).Log("exit", <-errs)
}
