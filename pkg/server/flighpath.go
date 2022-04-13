package server

import (
	"context"
	"net"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "passenger-tracker/api/pb/v1alpha1/flightpath"
	"passenger-tracker/internal/errors"
	"passenger-tracker/internal/flightpath"
	"passenger-tracker/internal/models"
)

// flightPathServer is the API handler and server for the microservice.
type flightPathServer struct {
}

// NewServer is the constructor for the server.
// It will return object with more fields when the microservice will grow.
func NewServer() pb.FlightPathServiceServer {
	return &flightPathServer{}
}

// GetFlightPath will take the request of list of paths between airports and
// returns the response with the start and end of the travel.
func (server *flightPathServer) GetFlightPath(
	_ context.Context,
	request *pb.GetFlightPathRequest,
) (*pb.GetFlightPathResponse, error) {
	// TODO: user authentication and authorization using ctx

	list := request.GetFlights()
	if len(list) == 0 {
		return nil, errors.InvalidArgumentError
	}

	var paths []*models.Path

	for _, path := range list {
		if !models.AirPortCode(path.GetStart()).Validate() ||
			!models.AirPortCode(path.GetEnd()).Validate() {
			return nil, errors.InvalidArgumentError
		}

		temp := &models.Path{}
		temp.ConvertFromAPIObject(path)

		paths = append(paths, temp)
	}

	return &pb.GetFlightPathResponse{
		Path:     flightpath.Normalize(paths).ConvertToAPIObject(),
		Datetime: timestamppb.Now(),
	}, nil
}

func Run(ctx context.Context, network, grpcAddr, httpAddr string) error {
	l, err := net.Listen(network, grpcAddr)
	if err != nil {
		return err
	}
	defer func() {
		if err := l.Close(); err != nil {
			glog.Errorf("Failed to close %s %s: %v", network, grpcAddr, err)
		}
	}()

	svr := grpc.NewServer()
	pb.RegisterFlightPathServiceServer(svr, NewServer())

	go func() {
		defer svr.GracefulStop()
		<-ctx.Done()
	}()

	mux := runtime.NewServeMux()

	pb.RegisterFlightPathServiceHandlerServer(ctx, mux, NewServer())
	httpSvr := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		glog.Infof("Shutting down the http gateway server")
		if err := httpSvr.Shutdown(context.Background()); err != nil {
			glog.Errorf("Failed to shutdown http gateway server: %v", err)
		}
	}()

	if err := httpSvr.ListenAndServe(); err != http.ErrServerClosed {
		glog.Errorf("Failed to listen and serve: %v", err)
		return err
	}

	return svr.Serve(l)
}
