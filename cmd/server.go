package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"wtpl/conf"
	"wtpl/pkg/logger"
	"wtpl/router"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	serverIP   string
	serverPort int
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"serve"},
	Short:   "run the http server",
	Long:    "run the http server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return startServer()
	},
}

func init() {
	serverCmd.Flags().StringVarP(&serverIP, "ip", "", "127.0.0.1", "http server bind ip")
	serverCmd.Flags().IntVarP(&serverPort, "port", "", 8080, "http server bind port")
	_ = viper.BindPFlag("server.ip", serverCmd.Flags().Lookup("ip"))
	_ = viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
}

func startServer() error {
	go startRpcServer()
	return startHttpServer()
}


func startHttpServer() error {
	addr := net.JoinHostPort(conf.GetConfig().Server.Ip, fmt.Sprintf("%d", conf.GetConfig().Server.Port))
	route := router.NewRoute()
	serve := &http.Server{
		Addr:    addr,
		Handler: route,
	}
	go func() {
		if conf.GetConfig().Server.TLSCert != "" && conf.GetConfig().Server.TLSKey != "" {
			if err := serve.ListenAndServeTLS(conf.GetConfig().Server.TLSCert, conf.GetConfig().Server.TLSKey); err != http.ErrServerClosed {

			}
		} else {
			if err := serve.ListenAndServe(); err != http.ErrServerClosed {

			}
		}
	}()

	signChan := make(chan os.Signal)
	signal.Notify(signChan, os.Kill, os.Interrupt, syscall.SIGHUP)
	for {
		select {
		case <-signChan:
			serve.Shutdown(context.Background())
			os.Exit(0)
		}
	}
}


func startRpcServer() {
	addr := net.JoinHostPort(conf.GetConfig().Rpc.Host, conf.GetConfig().Rpc.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.SErrorf("create rpc tcp listener failed:%v", err)
		return
	}
	s := grpc.NewServer()
	// s.RegisterService()
	if err := s.Serve(listener); err != nil {
		logger.SErrorf("rpc server start failed:%v", err)
	}
}