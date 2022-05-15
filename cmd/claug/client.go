package main

import (
	"context"
	"net"
	"time"

	"github.com/sergei-dyshel/claug/internal/server"

	"github.com/sergei-dyshel/claug/api"
	"github.com/sergei-dyshel/claug/internal/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func startEmbeddedServer(listener *bufconn.Listener) chan error {
	logger := getLogger("server")
	ch := make(chan error, 1)
	go func() {
		ch <- server.Run(listener, logger)
	}()
	return ch
}

func createConn(logger utils.Logger) (conn *grpc.ClientConn, err error) {
	sock, err := getSocket()
	if err != nil {
		return nil, err
	}
	err_ch := make(chan error, 1)
	err_ch <- nil
	if sock == embedSocket {
		logger.Infof("Using in-process server")
		// TODO: tweak buffer size
		listener := bufconn.Listen(4096)
		err_ch = startEmbeddedServer(listener) // TODO: check error
		conn, err = grpc.Dial(
			"bufnet",
			grpc.WithDialer(func(s string, duration time.Duration) (net.Conn, error) {
				return listener.Dial()
			}),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			err = utils.Wrapf(err, "failed to create embedded client-server connection")
		}
	} else {
		conn, err = grpc.Dial(
			"unix://"+sock,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			err = utils.Wrapf(err, "could not create client connection")
		}
	}
	return
}

func wrapClient(
	f func(client api.ServiceClient) error,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		err := initLogging(true /* client */)
		if err != nil {
			return err
		}
		logger := getLogger("client")
		conn, err := createConn(logger)
		if err != nil {
			return err
		}

		client := api.NewServiceClient(conn)
		err = f(client)
		if err != nil {
			return utils.Wrapf(err, "command failed")
		}
		return nil
	}
}

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping server",
	Long:  "Useful for checking that server is running and accepting commands",
	RunE: wrapClient(func(client api.ServiceClient) error {
		_, err := client.Ping(context.Background(), utils.Empty)
		return err
	}),
}

var exitServerCmd = &cobra.Command{
	Use:   "exit-server",
	Short: "Exit server",
	Long:  "Exit server immediately",
	RunE: wrapClient(func(client api.ServiceClient) error {
		_, err := client.Exit(context.Background(), utils.Empty)
		return err
	}),
}

var pressEnterCmd = &cobra.Command{
	Use:   "press-enter",
	Short: "Save current command in history and start it",
	Long:  "Intended for use as tmux binding, to save to history on each pressed Enter",
	RunE: wrapClient(func(client api.ServiceClient) error {
		_, err := client.PressEnter(context.Background(), utils.Empty)
		return err
	}),
}
