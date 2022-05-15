package main

import (
	"fmt"
	"os"
	"path"

	"github.com/sergei-dyshel/claug/cmd/claug/flags"

	"github.com/sergei-dyshel/claug/internal/utils"
	"github.com/spf13/cobra"
)

const (
	progNameLower = "claug"
	progNameUpper = "CLAUG"

	defaultRootName = "." + progNameLower
	defaultSocket   = "default"
	rootDirEnv      = progNameUpper + "_ROOT"
	socketEnv       = progNameUpper + "_SOCKET"
	serverLogEnv    = progNameUpper + "_SERVER_LOG"
	clientLogEnv    = progNameUpper + "_CLIENT_LOG"
	configEnv       = progNameUpper + "_CONFIG"

	defaultConfigName = "config.yml"

	embedSocket = "embed"
	logStderr   = "-"
)

var logFile *os.File

func getRootDir() (string, error) {
	dir := utils.FirstNonZero(
		flags.RootDir,
		os.Getenv(rootDirEnv),
		path.Join(utils.HomeDir(), defaultRootName),
	)
	err := os.MkdirAll(dir, os.ModePerm)
	return dir, utils.Wrapf(err, "could not create root directory %s", dir)
}

func addCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().
		CountVarP(&flags.Verbose, "verbose", "v", utils.Doc(`
			Log verbosity ''level''.
			Use -vv and -vvv for increased verbosity.
			Client will pass the level to server which will use it during processing the command`,
		))
	cmd.PersistentFlags().
		StringVarP(&flags.RootDir, "root-dir", "r", "", utils.Docf(`
			Use ''path'' as root directory to store logs, history etc.
			When not specified will take from $%s.
			Defaults to ~/%s.
			`, defaultRootName, rootDirEnv,
		))
	cmd.PersistentFlags().
		StringVarP(&flags.Socket, "socket", "s", "", utils.Docf(`
			Use socket ''path'' for client-server communication.
			If not specified, will be taken from $%s.
			Defaults to <root>/%s.sock.
			If just name without extension and directory is specified, will use <root>/<socket>.sock .
			Use e.g. ./name.sock for socket in current directory.
			When set to "%s", will create both server and client in the same process during invokation (usefull for debugging).
			`, socketEnv, defaultSocket, embedSocket,
		))
	cmd.PersistentFlags().
		StringVarP(&flags.LogFileName, "log", "l", "", utils.Docf(`
			Write log to ''file''. Specify '-' for stderr.
			If not specified, sever logs will be output to %s and client logs will be output to %s.
			Otherwise will output to stderr when running in TTY or to  <root>/server.<socket>.log, <root>/client.<socket>.log respectively.
			`, serverLogEnv, clientLogEnv,
		))
	cmd.PersistentFlags().
		StringVarP(&flags.ConfigFileName, "config", "c", "", utils.Docf(`
			Read config from ''file''.
			If not specified, will read from $% or from <root>/%s.
		`, configEnv, defaultConfigName))
}

func getSocket() (string, error) {
	sock := utils.FirstNonZero(flags.Socket, os.Getenv(socketEnv), "default")
	if sock == embedSocket {
		return embedSocket, nil
	}
	if path.Ext(sock) == "" && path.Base(sock) == sock {
		root, err := getRootDir()
		if err != nil {
			return "", err
		}
		sock = path.Join(root, sock+".sock")
	}
	if err := os.MkdirAll(path.Dir(sock), os.ModePerm); err != nil {
		return "", utils.Wrapf(err, "failed to create socket directory")
	}
	return sock, nil
}

func initLogging(client bool) error {
	sock, err := getSocket()
	if err != nil {
		return err
	}

	logFileName := logStderr
	if sock == "" {
		logFileName = utils.FirstNonZero(flags.LogFileName, logFileName)
	} else {
		envName := serverLogEnv
		typeStr := "server"
		if client {
			envName = clientLogEnv
			typeStr = "client"
		}

		root, err := getRootDir()
		if err != nil {
			return err
		}
		if !utils.IsTerminal(os.Stderr) {
			logFileName = path.Join(root, fmt.Sprintf("%s.%s.log", typeStr, path.Base(sock)))
		}
		envFileName := os.Getenv(envName)
		logFileName = utils.FirstNonZero(flags.LogFileName, envFileName, logFileName)
	}

	logFile = os.Stderr
	if logFileName != logStderr {
		var err error
		logFile, err = os.Create(logFileName)
		if err != nil {
			return utils.Wrapf(err, "could not open/create log file name %s", logFileName)
		}
	}
	return nil
}

func getLogger(prefix string) utils.Logger {
	showTimestamp := true
	if logFile == os.Stderr {
		showTimestamp = false
	}
	return utils.NewLogger(logFile, prefix, showTimestamp)
}

func getConfigFileName() (string, error) {
	fname := utils.FirstNonZero(flags.ConfigFileName, os.Getenv(configEnv))
	if fname != "" {
		return fname, nil
	}

	root, err := getRootDir()
	if err != nil {
		return "", err
	}

	return path.Join(root, defaultConfigName), nil
}
