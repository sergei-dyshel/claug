package main

import (
	"net"
	"os"

	"github.com/sergei-dyshel/claug/internal/server"
	"github.com/sergei-dyshel/claug/internal/utils"
	"github.com/spf13/cobra"
)

const (
	defaultPort = 3489
	// TODO: remove
	// defaultLogname = ".claug.log"
)

func startServer(*cobra.Command, []string) error {
	err := initLogging(false /* client */)
	if err != nil {
		return err
	}
	sock, err := getSocket()
	if err != nil {
		return err
	}

	if err := os.RemoveAll(sock); err != nil {
		return utils.Wrapf(err, "Can not delete socket file")
	}
	listener, err := net.Listen("unix", sock)
	if err != nil {
		return utils.Wrapf(err, "failed to listen on socket %s", sock)
	}
	logger := getLogger("server")
	logger.Infof("Listening on %s", listener.Addr().String())
	return server.Run(listener, logger)
}

func startServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start server process",
		RunE:  startServer,
	}
	return cmd
}
