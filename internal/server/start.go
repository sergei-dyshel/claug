package server

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/sergei-dyshel/claug/pkg/tmux"

	"github.com/sergei-dyshel/claug/api"
	"github.com/sergei-dyshel/claug/internal/utils"
	"google.golang.org/grpc"
)

func Run(listener net.Listener, logger utils.Logger) error {
	server := grpc.NewServer()
	api.RegisterServiceServer(server, &serviceServer{server: server, logger: logger})
	err := server.Serve(listener)
	if err != nil {
		return utils.Wrapf(err, "failed to serve requests")
	}
	return nil
}

type serviceServer struct {
	api.UnimplementedServiceServer

	logger utils.Logger
	server *grpc.Server
}

func (server *serviceServer) Ping(
	ctx context.Context,
	req *api.Empty,
) (*api.Empty, error) {
	server.logger.Infof("Requested ping")
	return utils.Empty, nil
}

func (server *serviceServer) Exit(context.Context, *api.Empty) (*api.Empty, error) {
	server.logger.Infof("Requsted exit")
	go func() {
		server.server.GracefulStop()
		log.Printf("Server stopped, exiting")
		os.Exit(0)
	}()
	return utils.Empty, nil
}

func (server *serviceServer) PressEnter(context.Context, *api.Empty) (*api.Empty, error) {
	server.logger.Infof("Requested press-enter")
	_, err := tmux.Run(tmux.CmdLine, &tmux.SendKeys{}, "Enter")
	return utils.Empty, err
}
